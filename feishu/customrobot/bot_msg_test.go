package customrobot_test

import (
	"context"
	"fmt"
	"github.com/magic-lib/go-notice-service/feishu/customrobot"
	"github.com/magic-lib/go-notice-service/msg"
	"github.com/magic-lib/go-notice-service/msgbuild"
	"testing"
)

var (
	botToken = "5226bced-fd8c-4ea3-b8db-b50e5f5a7733"
)

func TestSendText(t *testing.T) {
	fsBot := customrobot.NewFeiShuBot(botToken)
	messages := msgbuild.NewMessageBuilder().WithTitle("这是一个标题").
		WithType(msg.MsgTypeText).
		WithContent("<at user_id=\"ou_xxx\">Tom</at> 新更新提醒").Build()

	resp, err := fsBot.Send(context.Background(), messages)
	fmt.Println(resp, err)
}

func TestSendPost(t *testing.T) {
	fsBot := customrobot.NewFeiShuBot(botToken)
	contentList := []any{
		[]any{
			map[string]any{
				"tag":  "text",
				"text": "项目有更新:",
			},
			map[string]any{
				"tag":     "at",
				"user_id": "all",
			},
		},
	}

	messages := msgbuild.NewMessageBuilder().WithTitle("这是一个标题").
		WithType(msg.MsgTypePost).
		WithContent(contentList).Build()
	resp, err := fsBot.Send(context.Background(), messages)
	fmt.Println(resp, err)
}

func TestSendShareChat(t *testing.T) {
	fsBot := customrobot.NewFeiShuBot(botToken)
	messages := msgbuild.NewMessageBuilder().
		WithTemplateId("share_chat_default").
		WithContent("oc_f5b1a7eb27ae2****339ff").Build()
	resp, err := fsBot.Send(context.Background(), messages)
	fmt.Println(resp, err)
}
func TestSendImage(t *testing.T) {
	fsBot := customrobot.NewFeiShuBot(botToken)
	messages := msgbuild.NewMessageBuilder().
		WithTemplateId("image_default").
		WithContent("img_ecffc3b9-8f14-400f-a014-05eca1a4310g").Build()
	resp, err := fsBot.Send(context.Background(), messages)
	fmt.Println(resp, err)
}

func TestSendCardWarning(t *testing.T) {
	fsBot := customrobot.NewFeiShuBot(botToken)
	messages := msgbuild.NewMessageBuilder().WithTitle("这是一个标题").
		WithTemplateId("interactive_warning").
		WithTemplateData(map[string]any{
			"subtitle":     "",
			"title_color":  "green",
			"content_tips": "<font color='orange'>提示句</font>",
		}).
		WithContent("**告警时间：**aaaaaaaa\\n**异常系统：**ccccc").Build()
	resp, err := fsBot.Send(context.Background(), messages)
	fmt.Println(resp, err)
}

// https://open.feishu.cn/cardkit/editor?cardId=AAqIK0h6leWlx
func TestSendCardWarning1(t *testing.T) {
	fsBot := customrobot.NewFeiShuBot(botToken)
	messages := msgbuild.NewMessageBuilder().
		WithType(msg.MsgTypeInteractive).
		WithContent(`{"config":{"update_multi":true},"i18n_elements":{"zh_cn":[{"tag":"markdown","content":"<font color='orange'>请及时处理该异常，避免影响正常业务</font>","text_align":"left","text_size":"normal","icon":{"tag":"standard_icon","token":"warning_outlined","color":"orange"}},{"tag":"markdown","content":"**告警时间：**@warning_time\\n**异常系统：**@warning_system","text_align":"left","text_size":"normal"},{"tag":"markdown","content":"","text_align":"left","text_size":"normal"},{"tag":"table","columns":[{"data_type":"text","name":"module","display_name":"异常模块","horizontal_align":"left","width":"auto"},{"data_type":"text","name":"reason","display_name":"异常原因","horizontal_align":"left","width":"auto"},{"data_type":"number","name":"count","display_name":"异常数量","horizontal_align":"right","width":"auto","format":{"precision":0}}],"rows":[],"row_height":"low","header_style":{"background_style":"none","bold":true,"lines":1},"page_size":5}]},"i18n_header":{"zh_cn":{"title":{"tag":"plain_text","content":"短信服务异常告警"},"subtitle":{"tag":"plain_text","content":""},"template":"orange"}}}`).Build()
	resp, err := fsBot.Send(context.Background(), messages)
	fmt.Println(resp, err)
}
