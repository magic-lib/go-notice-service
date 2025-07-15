package msgbuild

import (
	"github.com/magic-lib/go-notice-service/msg"
)

// MessageBuilder 消息构建器
type MessageBuilder struct {
	msg messageImpl
}

// NewMessageBuilder 创建消息构建器
func NewMessageBuilder() *MessageBuilder {
	return &MessageBuilder{
		msg: messageImpl{
			receivers:    make([]*msg.Receiver, 0),
			options:      make(map[string]any),
			templateData: make(map[string]any),
		},
	}
}

// WithTemplateId 设置模版ID
func (b *MessageBuilder) WithTemplateId(templateId string) *MessageBuilder {
	b.msg.templateId = templateId
	return b
}

// WithTemplateData 设置模版数据
func (b *MessageBuilder) WithTemplateData(data map[string]any) *MessageBuilder {
	if len(data) > 0 {
		for k, v := range data {
			b.msg.templateData[k] = v
		}
	}
	return b
}

// WithChannel 设置消息渠道
func (b *MessageBuilder) WithChannel(channel msg.ChannelType, channelKey msg.ChannelKey) *MessageBuilder {
	b.msg.channel = channel
	b.msg.channelKey = channelKey
	return b
}
func (b *MessageBuilder) WithChannelAdapter(adapter ChannelAdapter) *MessageBuilder {
	channels := adapter.SupportedChannels()
	if len(channels) == 0 {
		return b
	}
	b.msg.channel = channels[0]
	b.msg.channelKey = adapter.ChannelKey()
	return b
}

// WithType 设置消息类型
func (b *MessageBuilder) WithType(msgType msg.MessageType) *MessageBuilder {
	b.msg.msgType = msgType
	return b
}

// WithReceiver 添加接收者
func (b *MessageBuilder) WithReceiver(typ msg.ReceiverType, id string) *MessageBuilder {
	b.msg.receivers = append(b.msg.receivers, &msg.Receiver{
		Type: typ,
		Id:   id,
	})
	return b
}
func (b *MessageBuilder) WithOneReceiver(receiver *msg.Receiver) *MessageBuilder {
	if receiver == nil {
		return b
	}
	b.msg.receivers = append(b.msg.receivers, receiver)
	return b
}

// WithTitle 设置标题
func (b *MessageBuilder) WithTitle(title string) *MessageBuilder {
	b.msg.title = title
	return b
}

// WithContent 设置内容
func (b *MessageBuilder) WithContent(content any) *MessageBuilder {
	b.msg.content = content
	return b
}

// WithOption 添加选项
func (b *MessageBuilder) WithOption(key string, value any) *MessageBuilder {
	b.msg.options[key] = value
	return b
}
func (b *MessageBuilder) WithOptions(data map[string]any) *MessageBuilder {
	if len(data) > 0 {
		for k, v := range data {
			b.msg.options[k] = v
		}
	}
	return b
}

// Build 构建消息
func (b *MessageBuilder) Build() msg.MessageTemplate {
	return &b.msg
}
