package customrobot_test

import (
	"context"
	"fmt"
	"github.com/magic-lib/go-notice-service/feishu/customrobot"
	"testing"
)

// SendSmsWarning sends SMS warning notification to FeiShu
func TestSendSmsWarning(t *testing.T) {
	fsBot := customrobot.NewFeiShuBot("e4ff48d4-ea6b-45b6-9217-35bc23e8a57f")
	resp, err := fsBot.SendByTemplateId(context.Background(), "card_warning", nil)

	// {"StatusCode":0,"StatusMessage":"success","code":0,"data":{},"msg":"success"}
	fmt.Println(resp, err)
}
