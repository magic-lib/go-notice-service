package msg

import "context"

// Receiver 接收者信息
type Receiver struct {
	Type  ReceiverType `json:"type"`  // 接收者类型
	Id    string       `json:"id"`    // 接收者ID
	Name  string       `json:"name"`  // 接收者名称（可选）
	Extra any          `json:"extra"` // 额外信息（如部门ID、租户ID等）
}

// Message 消息接口，定义发送消息的基本行为
type Message interface {
	Channel() (ChannelType, ChannelKey)       // 消息渠道
	MsgType() MessageType                     // MsgType 消息类型（文本、图片、卡片等）
	Receivers() []*Receiver                   // Receivers 消息接收者（用户ID、群组ID等）
	Title() string                            // 消息标题
	Content() any                             // 消息内容
	Options() map[string]any                  // 可选参数（如超时时间、优先级等）
	Validate() error                          // 验证消息有效性
	Send(ctx context.Context) (string, error) // 执行发送
}

// MessageTemplate 模板消息接口
type MessageTemplate interface {
	Message
	TemplateId() string           // 模板ID
	TemplateData() map[string]any // 模板变量
}
