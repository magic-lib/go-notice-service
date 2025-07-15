package customrobot

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/magic-lib/go-notice-service/msg"
	"github.com/magic-lib/go-notice-service/msgbuild"
	"github.com/magic-lib/go-plat-curl/curl"
	"github.com/magic-lib/go-plat-utils/cond"
	"github.com/magic-lib/go-plat-utils/conv"
	"github.com/magic-lib/go-plat-utils/templates"
	"github.com/tidwall/sjson"
	"golang.org/x/net/context"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultCustomRobotFeiShuUrlPath = "https://open.feishu.cn/open-apis/bot/v2/hook/"
	tmplPath                        = "tmpl"
	templateSuffix                  = ".json"
)

var (
	allTemplateFiles = make(map[string]string)
)

type feiShuBot struct {
	botUrl string
}

func NewFeiShuBot(botUrl string) *feiShuBot {
	if len(allTemplateFiles) == 0 {
		templateTemp, err := readAllTemplateFiles(tmplPath)
		if err != nil {
			panic(err)
		}
		allTemplateFiles = templateTemp
	}
	return &feiShuBot{
		botUrl: botUrl,
	}
}

func (f *feiShuBot) getFeiShuUrl() (string, error) {
	if f.botUrl == "" {
		return "", fmt.Errorf("invalid bot url: empty, default: %s", defaultCustomRobotFeiShuUrlPath)
	}
	if strings.HasPrefix(f.botUrl, defaultCustomRobotFeiShuUrlPath) {
		return f.botUrl, nil
	}
	if cond.IsUUID(f.botUrl) {
		f.botUrl = defaultCustomRobotFeiShuUrlPath + f.botUrl
		return f.botUrl, nil
	}
	return f.botUrl, fmt.Errorf("invalid bot url: %s", f.botUrl)
}
func (f *feiShuBot) getFeiShuBotToken() (string, error) {
	rawUrl, _ := f.getFeiShuUrl()
	if rawUrl == "" {
		return "", fmt.Errorf("invalid bot url: empty")
	}
	parsedURL, err := url.Parse(rawUrl)
	if err != nil {
		return "", fmt.Errorf("解析URL失败: %v", err)
	}

	path := parsedURL.Path

	// 按 "/" 分割路径，取最后一段即为Token
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		return "", fmt.Errorf("URL路径格式异常")
	}

	token := parts[len(parts)-1]
	if token == "" {
		return "", fmt.Errorf("未找到有效的Token")
	}
	if cond.IsUUID(token) {
		return strings.TrimSpace(token), nil
	}

	return token, fmt.Errorf("URL Token 异常: %s", token)
}
func (f *feiShuBot) getFeiShuSignature(timeNow time.Time) (string, error) {
	token, err := f.getFeiShuBotToken()
	if err != nil {
		return "", err
	}
	signKey := fmt.Sprintf("%d\n%s\n", timeNow.Unix(), token)
	h := hmac.New(sha256.New, []byte(signKey))
	h.Write([]byte{}) // Empty data as in original
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

func (f *feiShuBot) botPost(ctx context.Context, text string) (string, error) {
	botUrl, err := f.getFeiShuUrl()
	if err != nil {
		return "", err
	}
	resp := curl.NewClient().NewRequest(&curl.Request{
		Url:    botUrl,
		Data:   text,
		Method: http.MethodPost,
	}).Submit(ctx)
	if resp.Error != nil {
		return "", resp.Error
	}
	return resp.Response, nil
}

func (f *feiShuBot) SendWithJson(ctx context.Context, jsonString string, messages msg.Message) (string, error) {
	optMap, err := f.getMessageMap(messages)
	if err != nil {
		return "", err
	}

	jsonString, _ = sjson.Set(jsonString, "sign", optMap["sign"])
	jsonString, _ = sjson.Set(jsonString, "timestamp", optMap["timestamp"])
	msgBody, err := templates.Template(jsonString, optMap)

	var data interface{}
	err = json.Unmarshal([]byte(msgBody), &data)
	if err != nil {
		return "", err
	}
	formattedJson := conv.String(data)

	resp, err := f.botPost(ctx, formattedJson)
	if err != nil {
		return "", err
	}
	return resp, nil
}

func (f *feiShuBot) getMessageMap(messages msg.Message) (map[string]any, error) {
	timeNow := time.Now()
	sign, err := f.getFeiShuSignature(timeNow)
	if err != nil {
		return nil, err
	}
	if messages == nil {
		messages = msgbuild.NewMessageBuilder().Build()
	}

	optMap := messages.Options()
	optMap["title"] = messages.Title()
	optMap["sign"] = sign
	optMap["timestamp"] = timeNow.Unix()

	if _, ok := optMap["title_color"]; !ok {
		optMap["title_color"] = "orange"
	}
	if _, ok := optMap["subtitle"]; !ok {
		optMap["subtitle"] = ""
	}
	return optMap, nil
}

func (f *feiShuBot) getContentByTemplateId(templateId string) string {
	if templateId == "" {
		return ""
	}
	fileName := templateId + templateSuffix
	if retStr, ok := allTemplateFiles[fileName]; ok {
		return retStr
	}
	return ""
}

func (f *feiShuBot) SendByTemplateId(ctx context.Context, templateId string, messages msg.Message) (string, error) {
	msgBodyTmpl := f.getContentByTemplateId(templateId)
	if msgBodyTmpl == "" {
		return "", fmt.Errorf("invalid templateId: %s", templateId)
	}
	return f.SendWithJson(ctx, msgBodyTmpl, messages)
}
