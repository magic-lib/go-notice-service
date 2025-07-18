package approbot_test

import (
	"context"
	"fmt"
	"github.com/magic-lib/go-notice-service/feishu/approbot"
	"github.com/magic-lib/go-notice-service/msg"
	"github.com/magic-lib/go-notice-service/msgbuild"
	"testing"
)

var (
	appId     = ""
	appSecret = ""
	openId    = ""
)

func TestSendText(t *testing.T) {
	fsBot := approbot.NewFeiShuAppBot(appId, appSecret)
	messages := msgbuild.NewMessageBuilder().WithTitle("这是一个标题").
		WithType(msg.MsgTypeText).
		WithReceiver(msg.ReceiverOpenId, openId).
		WithContent("<at user_id=\"ou_xxx\">Tom</at> 新更新提醒").Build()

	resp, err := fsBot.Send(context.Background(), messages)
	fmt.Println(resp, err)
}

func TestSendPost(t *testing.T) {
	fsBot := approbot.NewFeiShuAppBot(appId, appSecret)
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
		WithReceiver(msg.ReceiverOpenId, openId).
		WithContent(contentList).Build()
	resp, err := fsBot.Send(context.Background(), messages)
	fmt.Println(resp, err)
}

func TestSendShareChat(t *testing.T) {
	fsBot := approbot.NewFeiShuAppBot(appId, appSecret)
	messages := msgbuild.NewMessageBuilder().
		WithType("share_chat").
		WithReceiver(msg.ReceiverOpenId, openId).
		WithContent("oc_f5b1a7eb27ae2****339ff").Build()
	resp, err := fsBot.Send(context.Background(), messages)
	fmt.Println(resp, err)
}

func TestSendImage(t *testing.T) {
	fsBot := approbot.NewFeiShuAppBot(appId, appSecret)
	messages := msgbuild.NewMessageBuilder().
		WithType("image").
		WithReceiver(msg.ReceiverOpenId, openId).
		WithContent("img_ecffc3b9-8f14-400f-a014-05eca1a4310g").Build()
	resp, err := fsBot.Send(context.Background(), messages)
	fmt.Println(resp, err)
}
func TestSendCardWarning1(t *testing.T) {
	fsBot := approbot.NewFeiShuAppBot(appId, appSecret)
	messages := msgbuild.NewMessageBuilder().
		WithType(msg.MsgTypeInteractive).
		WithReceiver(msg.ReceiverOpenId, openId).
		WithContent(`{"config":{"update_multi":true},"i18n_elements":{"zh_cn":[{"tag":"markdown","content":"<font color='orange'>请及时处理该异常，避免影响正常业务</font>","text_align":"left","text_size":"normal","icon":{"tag":"standard_icon","token":"warning_outlined","color":"orange"}},{"tag":"markdown","content":"**告警时间：**@warning_time\\n**异常系统：**@warning_system","text_align":"left","text_size":"normal"},{"tag":"markdown","content":"","text_align":"left","text_size":"normal"},{"tag":"table","columns":[{"data_type":"text","name":"module","display_name":"异常模块","horizontal_align":"left","width":"auto"},{"data_type":"text","name":"reason","display_name":"异常原因","horizontal_align":"left","width":"auto"},{"data_type":"number","name":"count","display_name":"异常数量","horizontal_align":"right","width":"auto","format":{"precision":0}}],"rows":[],"row_height":"low","header_style":{"background_style":"none","bold":true,"lines":1},"page_size":5}]},"i18n_header":{"zh_cn":{"title":{"tag":"plain_text","content":"短信服务异常告警"},"subtitle":{"tag":"plain_text","content":""},"template":"orange"}}}`).Build()
	resp, err := fsBot.Send(context.Background(), messages)
	fmt.Println(resp, err)
}

func TestSendCardWarning(t *testing.T) {
	fsBot := approbot.NewFeiShuAppBot(appId, appSecret)
	messages := msgbuild.NewMessageBuilder().WithTitle("这是一个标题").
		WithType(msg.MsgTypeInteractive).
		WithTemplateId("interactive_warning").
		WithTemplateData(map[string]any{
			"subtitle":    "",
			"title_color": "green",
		}).
		WithReceiver(msg.ReceiverUnionId, openId).
		WithContent("**告警时间：**aaaaaaaa\\n**异常系统：**ccccc").Build()
	resp, err := fsBot.Send(context.Background(), messages)
	fmt.Println(resp, err)
}
func TestSendCardWarning3(t *testing.T) {
	fsBot := approbot.NewFeiShuAppBot(appId, appSecret)
	contentStr := `[
      {
        "tag": "markdown",
        "content": "头部文件",
        "text_align": "left",
        "text_size": "normal",
        "icon": {
          "tag": "standard_icon",
          "token": "warning_outlined",
          "color": "red"
        }
      },
      {
        "tag": "markdown",
        "content": "文件内容",
        "text_align": "left",
        "text_size": "normal_v2"
      }
    ]`

	messages := msgbuild.NewMessageBuilder().WithTitle("这是一个标题").
		WithType(msg.MsgTypeInteractive).
		WithTemplateId("interactive_card").
		WithTemplateData(map[string]any{
			"subtitle":     "",
			"title_color":  "red",
			"content_tips": "<font color='orange'>提示句</font>",
		}).
		WithReceiver(msg.ReceiverUserId, openId).
		WithContent(contentStr).Build()
	resp, err := fsBot.Send(context.Background(), messages)
	fmt.Println(resp, err)
}
func TestSendCardWarning12(t *testing.T) {
	fsBot := approbot.NewFeiShuAppBot(appId, appSecret)
	userList, err := fsBot.UserIdMapByMobiles("open_id", []string{""})
	fmt.Println(userList, err)
}
