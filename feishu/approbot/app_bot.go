package approbot

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-multierror"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/magic-lib/go-notice-service/msg"
	"github.com/magic-lib/go-plat-utils/conv"
	"github.com/samber/lo"
)

type feiShuBot struct {
	appId     string
	appSecret string
	client    *lark.Client
}

func newFeiShuAppBot(appId, appSecret string) (*feiShuBot, error) {
	client := lark.NewClient(appId, appSecret)
	if client == nil {
		return nil, fmt.Errorf("new client failed")
	}
	return &feiShuBot{
		appId:     appId,
		appSecret: appSecret,
		client:    client,
	}, nil
}

func (m *feiShuBot) SendToOne(ctx context.Context, msgInfo msg.Message) ([]string, error) {
	if m.client == nil {
		return nil, fmt.Errorf("client is nil")
	}
	receives := msgInfo.Receivers()
	if len(receives) == 0 {
		return nil, fmt.Errorf("receives is empty")
	}
	var retString = make([]string, 0)
	var retErr error
	lo.ForEach(receives, func(oneReceiver *msg.Receiver, index int) {
		retStr, err := m.sendToOne(ctx, oneReceiver, msgInfo)
		if err != nil {
			retErr = multierror.Append(retErr, err)
		} else {
			retString = append(retString, retStr)
		}
	})
	if retErr != nil {
		return retString, retErr
	}
	return retString, nil
}

func (m *feiShuBot) sendToOne(ctx context.Context, oneReceiver *msg.Receiver, msgInfo msg.Message) (string, error) {
	if m.client == nil {
		return "", fmt.Errorf("client is nil")
	}
	content, err := m.getContentJsonByType(msgInfo.MsgType(), msgInfo.Content())
	if err != nil {
		return "", err
	}

	// 创建请求对象
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(oneReceiver.Type.String()).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(oneReceiver.Id).
			MsgType(msgInfo.MsgType().String()).
			Content(content).
			Build()).Build()

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

func (m *feiShuBot) getContentJsonByType(msgType msg.MessageType, content any) (string, error) {
	if msgType == msg.MsgTypeText {
		content := map[string]any{
			"text": conv.String(content),
		}
		return conv.String(content), nil
	}

	return "", fmt.Errorf("msg type not support")
}
