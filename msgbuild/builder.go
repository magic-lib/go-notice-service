package msgbuild

import "github.com/magic-lib/go-notice-service/msg"

// MessageBuilder 消息构建器
type MessageBuilder struct {
	msg messageImpl
}

// NewMessageBuilder 创建消息构建器
func NewMessageBuilder() *MessageBuilder {
	return &MessageBuilder{
		msg: messageImpl{
			receivers: make([]*msg.Receiver, 0),
			options:   make(map[string]any),
		},
	}
}

// WithChannel 设置消息渠道
func (b *MessageBuilder) WithChannel(channel msg.ChannelType) *MessageBuilder {
	b.msg.channel = channel
	return b
}

// WithType 设置消息类型
func (b *MessageBuilder) WithType(msgType msg.MessageType) *MessageBuilder {
	b.msg.msgType = msgType
	return b
}

// WithReceiver 添加接收者
func (b *MessageBuilder) WithReceiver(typ msg.ReceiverType, id, name string, extra any) *MessageBuilder {
	b.msg.receivers = append(b.msg.receivers, &msg.Receiver{
		Type:  typ,
		Id:    id,
		Name:  name,
		Extra: extra,
	})
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

// Build 构建消息
func (b *MessageBuilder) Build() msg.Message {
	return &b.msg
}
