package approbot

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-multierror"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/magic-lib/go-notice-service/msg"
	"github.com/magic-lib/go-notice-service/msgbuild"
	"github.com/magic-lib/go-plat-utils/conv"
	"github.com/magic-lib/go-plat-utils/templates"
	"github.com/samber/lo"
	"log"
	"strings"
)

type feiShuBot struct {
	appId     string
	appSecret string
	client    *lark.Client
}

var (
	allTemplateFiles = make(map[string]string)
	tmplPath         = "tmpl"
	templateSuffix   = ".json"
)

func NewFeiShuAppBot(appId, appSecret string) *feiShuBot {
	if len(allTemplateFiles) == 0 {
		templateTemp, err := readAllTemplateJsonFiles(tmplPath)
		if err != nil {
			panic(err)
		}
		allTemplateFiles = templateTemp
	}

	client := lark.NewClient(appId, appSecret)
	if client == nil {
		log.Println("new feiShuBot failed")
		return new(feiShuBot)
	}
	return &feiShuBot{
		appId:     appId,
		appSecret: appSecret,
		client:    client,
	}
}

func (m *feiShuBot) getMessageMap(messages msg.Message, mapData map[string]any) (map[string]any, error) {
	if messages == nil {
		messages = msgbuild.NewMessageBuilder().Build()
	}
	if mapData == nil {
		mapData = make(map[string]any)
	}

	optMap := messages.Options()
	for k, v := range optMap {
		if _, ok := mapData[k]; ok {
			continue
		}
		mapData[k] = v
	}
	if messages.Title() != "" {
		mapData["title"] = messages.Title()
	}

	mapData["content"] = conv.String(messages.Content())
	return mapData, nil
}
func (m *feiShuBot) Send(ctx context.Context, msgInfo msg.MessageTemplate) (string, error) {
	if m.client == nil {
		return "", fmt.Errorf("client is nil")
	}

	receives := msgInfo.Receivers()
	if len(receives) == 0 {
		return "", fmt.Errorf("receives is empty")
	}

	var content, err = m.getContent(msgInfo)
	if err != nil {
		return "", err
	}

	var retString = make([]string, 0)
	var retErr error
	lo.ForEach(receives, func(oneReceiver *msg.Receiver, index int) {
		retStr, err := m.sendToOne(ctx, oneReceiver, content, msgInfo)
		if err != nil {
			retErr = multierror.Append(retErr, err)
		} else {
			retString = append(retString, retStr)
		}
	})

	if retErr != nil {
		return "", retErr
	}
	if len(retString) == 0 {
		return "", nil
	}
	return retString[0], nil
}

func (m *feiShuBot) getContent(msgInfo msg.MessageTemplate) (string, error) {
	optMap, err := m.getMessageMap(msgInfo, msgInfo.TemplateData())
	if err != nil {
		return "", err
	}
	content := ""
	templateId := msgInfo.TemplateId()
	if templateId != "" {
		msgBodyTmpl := m.getContentByTemplateId(templateId)
		if msgBodyTmpl == "" {
			return "", fmt.Errorf("invalid templateId: %s", templateId)
		}
		content, _ = templates.Template(msgBodyTmpl, optMap)
		content = strings.ReplaceAll(content, "<no value>", "")
	} else {
		content, err = m.getContentJsonByType(msgInfo.MsgType(), msgInfo)
		if err != nil {
			return "", err
		}
	}
	return content, nil
}

func (m *feiShuBot) getContentByTemplateId(templateId string) string {
	if templateId == "" {
		return ""
	}
	fileName := templateId + templateSuffix
	if retStr, ok := allTemplateFiles[fileName]; ok {
		return retStr
	}
	return ""
}

func (m *feiShuBot) sendToOne(ctx context.Context, oneReceiver *msg.Receiver, content string, msgInfo msg.Message) (string, error) {
	if m.client == nil {
		return "", fmt.Errorf("client is nil")
	}

	if content == "" {
		return "", fmt.Errorf("content is empty")
	}

	// 创建请求对象
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(oneReceiver.Type.String()).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(oneReceiver.Id).
			MsgType(msgInfo.MsgType().String()).
			Content(content).
			Build()).Build()

	log.Printf("app bot sendToOne req, type: %s, receiver: %s, msgType: %s", oneReceiver.Type.String(), oneReceiver.Id, msgInfo.MsgType().String())

	// 发起请求
	resp, err := m.client.Im.V1.Message.Create(ctx, req)
	if err != nil {
		return "", err
	}
	if !resp.Success() {
		return "", fmt.Errorf("logId: %s, error response: \n%s", resp.RequestId(), larkcore.Prettify(resp.CodeError))
	}
	return larkcore.Prettify(resp), nil
}

// https://open.feishu.cn/document/server-docs/im-v1/message-content-description/create_json
// https://open.feishu.cn/cardkit?from=open_docs_tool_overview
func (m *feiShuBot) getContentJsonByType(msgType msg.MessageType, msgInfo msg.Message) (string, error) {
	if msgType == msg.MsgTypeText {
		content := map[string]any{
			"text": conv.String(msgInfo.Content()),
		}
		return conv.String(content), nil
	} else if msgType == msg.MsgTypePost {
		content := map[string]map[string]any{
			"zh_cn": {},
		}
		if msgInfo.Title() != "" {
			content["zh_cn"]["title"] = msgInfo.Title()
		}
		if msgInfo.Content() != "" {
			content["zh_cn"]["content"] = msgInfo.Content()
		}
		return conv.String(content), nil
	} else if msgType == "share_chat" {
		content := map[string]string{
			"chat_id": conv.String(msgInfo.Content()),
		}
		return conv.String(content), nil
	} else if msgType == msg.MsgTypeImage {
		content := map[string]string{
			"image_key": conv.String(msgInfo.Content()),
		}
		return conv.String(content), nil
	} else if msgType == msg.MsgTypeInteractive {
		return conv.String(msgInfo.Content()), nil
	}

	return "", fmt.Errorf("msg type not support")
}
