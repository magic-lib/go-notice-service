package msgbuild_test

import (
	"context"
	"github.com/magic-lib/go-notice-service/msg"
	"github.com/magic-lib/go-notice-service/msgbuild"
	"log"
	"testing"
)

func TestNewChannelAdapterManager(t *testing.T) {
	// 注册飞书适配器
	registry := msgbuild.NewChannelAdapterManager()
	registry.Register(NewFeishuAdapter()) // 假设已实现飞书适配器

	// 创建消息发送器
	sender := msgbuild.NewMessageSender(registry)

	// 构建一条飞书文本消息
	msgInfo := msgbuild.NewMessageBuilder().
		WithChannel(msg.ChannelFeiShu).
		WithType(msg.TypeText).
		WithReceiver(msg.ReceiverUser, "user_id_123", "", nil).
		WithTitle("系统通知").
		WithContent("您的订单已发货！").
		WithOption("priority", "high").
		Build()

	// 发送消息
	ctx := context.Background()
	msgID, err := sender.Send(ctx, msgInfo)
	if err != nil {
		log.Fatalf("发送失败: %v", err)
	}
	log.Printf("消息发送成功，ID: %s", msgID)
}
