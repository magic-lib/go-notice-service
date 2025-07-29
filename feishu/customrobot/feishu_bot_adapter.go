package customrobot

import (
	"fmt"
	"github.com/magic-lib/go-notice-service/msg"
	"github.com/magic-lib/go-notice-service/msgbuild"
)
import "context"

// feiShuBotAdapter 飞书机器人适配器
type feiShuBotAdapter struct {
	feiShuBot  *feiShuBot
	channelKey msg.ChannelKey
}

func (f *feiShuBotAdapter) ChannelKey() msg.ChannelKey {
	return f.channelKey
}

// NewFeiShuCustomRoBotAdapter 创建飞书适配器
func NewFeiShuCustomRoBotAdapter(channelKey msg.ChannelKey, token string) (msgbuild.ChannelAdapter, error) {
	if token == "" {
		return nil, fmt.Errorf("token is null")
	}
	fb := NewFeiShuBot(token)
	if channelKey == "" {
		channelKeyTemp, err := fb.getFeiShuBotToken()
		if err != nil || channelKeyTemp == "" {
			return nil, fmt.Errorf("获取飞书机器人Token失败: %v", err)
		}
		channelKey = msg.ChannelKey(channelKeyTemp)
	}
	return &feiShuBotAdapter{
		feiShuBot:  fb,
		channelKey: channelKey,
	}, nil
}

func (f *feiShuBotAdapter) SupportedChannels() []msg.ChannelType {
	return []msg.ChannelType{
		msg.ChannelFeiShuCustomRoBot,
	}
}
func (f *feiShuBotAdapter) Send(ctx context.Context, message msg.MessageTemplate) (string, error) {
	return f.feiShuBot.Send(ctx, message)
}
