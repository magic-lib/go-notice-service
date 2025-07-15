package customrobot

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/magic-lib/go-notice-service/msg"
	"github.com/magic-lib/go-notice-service/msgbuild"
	"github.com/magic-lib/go-plat-curl/curl"
	"github.com/magic-lib/go-plat-utils/cond"
	"github.com/magic-lib/go-plat-utils/conv"
	"github.com/magic-lib/go-plat-utils/templates"
	"github.com/tidwall/sjson"
	"golang.org/x/net/context"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultCustomRobotFeiShuUrlPath = "https://open.feishu.cn/open-apis/bot/v2/hook/"
	tmplPath                        = "tmpl"
	templateSuffix                  = ".json"
)

var (
	allTemplateFiles = make(map[string]string)
)

type sendResp struct {
	Code int `json:"code"`
	Data struct {
	} `json:"data"`
	Msg string `json:"msg"`
}

type feiShuBot struct {
	botUrl string
}

func NewFeiShuBot(botUrl string) *feiShuBot {
	if len(allTemplateFiles) == 0 {
		templateTemp, err := readAllTemplateJsonFiles(tmplPath)
		if err != nil {
			panic(err)
		}
		allTemplateFiles = templateTemp
	}
	return &feiShuBot{
		botUrl: botUrl,
	}
}

func (f *feiShuBot) getFeiShuUrl() (string, error) {
	if f.botUrl == "" {
		return "", fmt.Errorf("invalid bot url: empty, default: %s", defaultCustomRobotFeiShuUrlPath)
	}
	if strings.HasPrefix(f.botUrl, defaultCustomRobotFeiShuUrlPath) {
		return f.botUrl, nil
	}
	if cond.IsUUID(f.botUrl) {
		f.botUrl = defaultCustomRobotFeiShuUrlPath + f.botUrl
		return f.botUrl, nil
	}
	return f.botUrl, fmt.Errorf("invalid bot url: %s", f.botUrl)
}
func (f *feiShuBot) getFeiShuBotToken() (string, error) {
	rawUrl, _ := f.getFeiShuUrl()
	if rawUrl == "" {
		return "", fmt.Errorf("invalid bot url: empty")
	}
	parsedURL, err := url.Parse(rawUrl)
	if err != nil {
		return "", fmt.Errorf("解析URL失败: %v", err)
	}

	path := parsedURL.Path

	// 按 "/" 分割路径，取最后一段即为Token
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		return "", fmt.Errorf("URL路径格式异常")
	}

	token := parts[len(parts)-1]
	if token == "" {
		return "", fmt.Errorf("未找到有效的Token")
	}
	if cond.IsUUID(token) {
		return strings.TrimSpace(token), nil
	}

	return token, fmt.Errorf("URL Token 异常: %s", token)
}
func (f *feiShuBot) getFeiShuSignature(timeNow time.Time) (string, error) {
	token, err := f.getFeiShuBotToken()
	if err != nil {
		return "", err
	}
	return genSign(token, timeNow.Unix())
}

func genSign(secret string, timestamp int64) (string, error) {
	//timestamp + key 做sha256, 再进行base64 encode
	stringToSign := fmt.Sprintf("%v\n%s", timestamp, secret)
	var data []byte
	h := hmac.New(sha256.New, []byte(stringToSign))
	_, err := h.Write(data)
	if err != nil {
		return "", err
	}
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature, nil
}

func (f *feiShuBot) botPost(ctx context.Context, text string) (string, error) {
	botUrl, err := f.getFeiShuUrl()
	if err != nil {
		return "", err
	}
	resp := curl.NewClient().NewRequest(&curl.Request{
		Url:    botUrl,
		Data:   text,
		Method: http.MethodPost,
	}).Submit(ctx)
	if resp.Error != nil {
		return "", resp.Error
	}
	sendBotResp := new(sendResp)
	err = conv.Unmarshal(resp.Response, sendBotResp)
	if err != nil {
		return "", err
	}
	if sendBotResp.Code != 0 {
		return "", fmt.Errorf("发送消息失败: %s", sendBotResp.Msg)
	}

	return resp.Response, nil
}

func (f *feiShuBot) sendWithJson(ctx context.Context, msgBody string, optMap map[string]any) (string, error) {
	if _, ok := optMap["sign"]; ok {
		if _, ok = optMap["timestamp"]; ok {
			msgBody, _ = sjson.Set(msgBody, "sign", optMap["sign"])
			msgBody, _ = sjson.Set(msgBody, "timestamp", optMap["timestamp"])
		}
	}

	var data interface{}
	err := json.Unmarshal([]byte(msgBody), &data)
	if err != nil {
		return "", err
	}
	formattedJson := conv.String(data)

	resp, err := f.botPost(ctx, formattedJson)
	if err != nil {
		return "", err
	}
	return resp, nil
}

func (f *feiShuBot) getMessageMap(messages msg.Message, mapData map[string]any) (map[string]any, error) {
	timeNow := time.Now()
	sign, err := f.getFeiShuSignature(timeNow)
	if err != nil {
		return nil, err
	}
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
	mapData["sign"] = sign
	mapData["timestamp"] = timeNow.Unix()

	if _, ok := mapData["title_color"]; !ok {
		mapData["title_color"] = "orange"
	}
	if _, ok := mapData["subtitle"]; !ok {
		mapData["subtitle"] = ""
	}

	mapData["content"] = conv.String(messages.Content())
	return mapData, nil
}

func (f *feiShuBot) getContentByTemplateId(templateId string) string {
	if templateId == "" {
		return ""
	}
	fileName := templateId + templateSuffix
	if retStr, ok := allTemplateFiles[fileName]; ok {
		return retStr
	}
	return ""
}

func (f *feiShuBot) getContentJsonByType(msgType msg.MessageType, msgInfo msg.Message) (string, error) {
	if msgType == msg.MsgTypeText {
		content := map[string]any{
			"msg_type": "text",
			"content": map[string]any{
				"text": conv.String(msgInfo.Content()),
			},
		}
		return conv.String(content), nil
	} else if msgType == msg.MsgTypePost {
		content := map[string]any{
			"msg_type": "post",
			"content": map[string]any{
				"post": map[string]any{
					"zh_cn": map[string]any{
						"title":   msgInfo.Title(),
						"content": msgInfo.Content(),
					},
				},
			},
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
		content := map[string]any{
			"msg_type": "interactive",
			"card":     conv.String(msgInfo.Content()),
		}
		return conv.String(content), nil
	}

	return "", fmt.Errorf("msg type not support")
}

// Send https://open.feishu.cn/document/client-docs/bot-v3/add-custom-bot?lang=zh-CN#%E6%94%AF%E6%8C%81%E5%8F%91%E9%80%81%E7%9A%84%E6%B6%88%E6%81%AF%E7%B1%BB%E5%9E%8B%E8%AF%B4%E6%98%8E
func (f *feiShuBot) Send(ctx context.Context, messages msg.MessageTemplate) (string, error) {
	if messages == nil {
		return "", fmt.Errorf("invalid message: empty")
	}

	optMap, err := f.getMessageMap(messages, messages.TemplateData())
	if err != nil {
		return "", err
	}

	templateId := messages.TemplateId()
	if templateId != "" {
		msgBodyTmpl := f.getContentByTemplateId(templateId)
		if msgBodyTmpl == "" {
			return "", fmt.Errorf("invalid templateId: %s", templateId)
		}
		msgBody, _ := templates.Template(msgBodyTmpl, optMap)

		return f.sendWithJson(ctx, msgBody, optMap)
	}
	content, err := f.getContentJsonByType(messages.MsgType(), messages)
	if err != nil {
		return "", err
	}
	return f.sendWithJson(ctx, content, optMap)
}
