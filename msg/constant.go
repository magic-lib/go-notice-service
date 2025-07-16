package msg

// ChannelType 消息渠道类型
type ChannelType string
type ChannelKey string //同类型可以注册多个，比如机器人，可以有多个来进行单独发送

const (
	ChannelKeyDefault        ChannelKey  = "default"
	ChannelFeiShu            ChannelType = "feishu"              // 飞书
	ChannelFeiShuAppRoBot    ChannelType = "feishu-app-robot"    // 飞书
	ChannelFeiShuCustomRoBot ChannelType = "feishu-custom-robot" // 飞书
	ChannelWeChat            ChannelType = "wechat"              // 微信
	ChannelDingTalk          ChannelType = "dingtalk"            // 钉钉
	ChannelEmail             ChannelType = "email"               // 邮件
)

// MessageType 消息内容类型
type MessageType string

const (
	MsgTypeText        MessageType = "text"        // 文本消息
	MsgTypePost        MessageType = "post"        // 富文本消息
	MsgTypeImage       MessageType = "image"       // 图片消息
	MsgTypeFile        MessageType = "file"        // 文件消息
	MsgTypeMarkdown    MessageType = "markdown"    // Markdown消息
	MsgTypeInteractive MessageType = "interactive" // 交互式卡片消息
	MsgTypeTemplate    MessageType = "template"    // 模板消息
)

func (mt MessageType) String() string {
	return string(mt)
}

// ReceiverType 接收者类型
type ReceiverType string

const (
	ReceiverOpenId  ReceiverType = "open_id"
	ReceiverUserId  ReceiverType = "user_id"
	ReceiverUnionId ReceiverType = "union_id"
	ReceiverEmail   ReceiverType = "email"
	ReceiverChatId  ReceiverType = "chat_id"
	ReceiverUser    ReceiverType = "user"  // 用户
	ReceiverGroup   ReceiverType = "group" // 群组
	ReceiverChat    ReceiverType = "chat"  // 会话
)

func (rt ReceiverType) String() string {
	return string(rt)
}
