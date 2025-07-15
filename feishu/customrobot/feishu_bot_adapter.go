package customrobot

import (
	"github.com/magic-lib/go-notice-service/msg"
	"github.com/magic-lib/go-notice-service/msgbuild"
)
import "context"

// feiShuBotAdapter 飞书机器人适配器
type feiShuBotAdapter struct {
	feiShuBot *feiShuBot
}

// NewFeiShuCustomRoBotAdapter 创建飞书适配器
func NewFeiShuCustomRoBotAdapter(token string) msgbuild.ChannelAdapter {
	return &feiShuBotAdapter{
		feiShuBot: NewFeiShuBot(token),
	}
}

func (f *feiShuBotAdapter) SupportedChannels() []msg.ChannelType {
	return []msg.ChannelType{msg.ChannelFeiShuCustomRoBot}
}
func (f *feiShuBotAdapter) Send(ctx context.Context, msg msg.Message) (string, error) {
	return "", nil
}
