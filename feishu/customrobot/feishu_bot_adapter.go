package customrobot

import (
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
func NewFeiShuCustomRoBotAdapter(channelKey msg.ChannelKey, token string) msgbuild.ChannelAdapter {
	fb := NewFeiShuBot(token)
	if channelKey == "" {
		channelKeyTemp, _ := fb.getFeiShuBotToken()
		channelKey = msg.ChannelKey(channelKeyTemp)
	}
	return &feiShuBotAdapter{
		feiShuBot:  fb,
		channelKey: channelKey,
	}
}

func (f *feiShuBotAdapter) SupportedChannels() []msg.ChannelType {
	return []msg.ChannelType{
		msg.ChannelFeiShuCustomRoBot,
	}
}
func (f *feiShuBotAdapter) Send(ctx context.Context, message msg.MessageTemplate) (string, error) {
	return f.feiShuBot.Send(ctx, message)
}
