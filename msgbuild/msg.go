package msgbuild

import "github.com/magic-lib/go-notice-service/msg"
import "context"

// messageImpl 消息接口的基础实现
type messageImpl struct {
	channel      msg.ChannelType
	channelKey   msg.ChannelKey
	msgType      msg.MessageType
	receivers    []*msg.Receiver
	title        string
	content      any
	options      map[string]any
	templateId   string
	templateData map[string]any
}

func (m *messageImpl) TemplateId() string                         { return m.templateId }
func (m *messageImpl) TemplateData() map[string]any               { return m.templateData }
func (m *messageImpl) Channel() (msg.ChannelType, msg.ChannelKey) { return m.channel, m.channelKey }
func (m *messageImpl) MsgType() msg.MessageType                   { return m.msgType }
func (m *messageImpl) Receivers() []*msg.Receiver                 { return m.receivers }
func (m *messageImpl) Title() string                              { return m.title }
func (m *messageImpl) Content() any                               { return m.content }
func (m *messageImpl) Options() map[string]any                    { return m.options }
func (m *messageImpl) Validate() error                            { /* 实现验证逻辑 */ return nil }
func (m *messageImpl) Send(ctx context.Context) (string, error) {
	return "", nil
}
