package approbot

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

// NewFeiShuAppRoBotAdapter 创建飞书适配器
func NewFeiShuAppRoBotAdapter(channelKey msg.ChannelKey, appId, appSecret string) msgbuild.ChannelAdapter {
	if channelKey == "" {
		channelKey = msg.ChannelKey(fmt.Sprintf("%s/%s", appId, appSecret))
	}

	return &feiShuBotAdapter{
		feiShuBot:  NewFeiShuAppBot(appId, appSecret),
		channelKey: channelKey,
	}
}

func (f *feiShuBotAdapter) SupportedChannels() []msg.ChannelType {
	return []msg.ChannelType{
		msg.ChannelFeiShuAppRoBot,
	}
}
func (f *feiShuBotAdapter) Send(ctx context.Context, message msg.MessageTemplate) (string, error) {
	return f.feiShuBot.Send(ctx, message)
}
