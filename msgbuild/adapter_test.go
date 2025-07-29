package msgbuild_test

import (
	"context"
	"fmt"
	"github.com/magic-lib/go-notice-service/feishu/approbot"
	"github.com/magic-lib/go-notice-service/msg"
	"github.com/magic-lib/go-notice-service/msgbuild"
	"github.com/magic-lib/go-plat-utils/conv"
	"log"
	"testing"
)

func TestFeiShuCustomRobotAdapterManager(t *testing.T) {
	appId := ""
	appSecret := ""
	openId := ""

	// 注册飞书适配器
	//customRoBot := customrobot.NewFeiShuCustomRoBotAdapter("", "")
	appRoBot, _ := approbot.NewFeiShuAppRoBotAdapter("", appId, appSecret)

	registry := msgbuild.NewChannelAdapterManager()
	//registry.RegisterAll(customRoBot, appRoBot)
	registry.RegisterAll(appRoBot)

	mm := registry.GetChannels()
	fmt.Println(conv.String(mm))

	// 创建消息发送器
	sender := msgbuild.NewMessageSender(registry)

	msgInfo := msgbuild.NewMessageBuilder().WithTitle("这是一个无用的标题").
		//WithChannelAdapter(customRoBot).
		WithTemplateId("text_default").
		WithContent("<at user_id=\"ou_xxx\">Tom</at> 新更新提醒").Build()
	msgID, err := sender.Send(context.Background(), msgInfo)
	if err != nil {
		log.Printf("发送失败: %v", err)
	}
	log.Printf("消息发送成功，ID: %s", msgID)

	msgInfo = msgbuild.NewMessageBuilder().WithTitle("这是一个无用的标题").
		WithChannelAdapter(appRoBot).
		WithType(msg.MsgTypeText).
		WithReceiver(msg.ReceiverOpenId, openId).
		WithContent("<at user_id=\"ou_xxx\">Tom</at> 新更新提醒").Build()
	msgID, err = sender.Send(context.Background(), msgInfo)
	if err != nil {
		log.Printf("发送失败: %v", err)
	}
	log.Printf("消息发送成功，ID: %s", msgID)
}
