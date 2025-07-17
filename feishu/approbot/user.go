package approbot

import (
	"context"
	"fmt"
	larkcore "github.com/larksuite/oapi-sdk-go/v3/core"
	larkcontact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
)

// https://open.feishu.cn/api-explorer/cli_a8e7721514529013?apiName=batch_get_id&from=op_doc_tab&project=contact&resource=user&version=v3
func (m *feiShuBot) getUserIdInfo(userIdType string, mobiles []string, emails []string) ([]*larkcontact.UserContactInfo, error) {
	if userIdType == "" || (userIdType != "open_id" && userIdType != "user_id" && userIdType != "union_id") {
		return nil, fmt.Errorf("userIdType is error: %s", userIdType)
	}
	if len(mobiles) == 0 && len(emails) == 0 {
		return nil, fmt.Errorf("mobiles and emails are empty")
	}

	getIdReq := larkcontact.NewBatchGetIdUserReqBodyBuilder()
	if len(mobiles) > 0 {
		getIdReq = getIdReq.Mobiles(mobiles)
	}
	if len(emails) > 0 {
		getIdReq = getIdReq.Emails(emails)
	}

	// 创建请求对象
	req := larkcontact.NewBatchGetIdUserReqBuilder().
		UserIdType(userIdType).
		Body(getIdReq.
			IncludeResigned(true).
			Build()).
		Build()

	resp, err := m.client.Contact.V3.User.BatchGetId(context.Background(), req)
	if err != nil {
		return nil, err
	}

	// 服务端错误处理
	if !resp.Success() {
		return nil, fmt.Errorf("logId: %s, error response: \n%s", resp.RequestId(), larkcore.Prettify(resp.CodeError))
	}

	userIdList := make([]*larkcontact.UserContactInfo, 0)
	if len(resp.Data.UserList) == 0 {
		return userIdList, nil
	}

	for _, user := range resp.Data.UserList {
		userId := larkcore.StringValue(user.UserId)
		if userId == "" {
			continue
		}
		userIdList = append(userIdList, user)
	}
	return userIdList, nil
}
func (m *feiShuBot) UserIdMapByMobiles(userIdType string, mobiles []string) ([]*larkcontact.UserContactInfo, error) {
	return m.getUserIdInfo(userIdType, mobiles, nil)
}
func (m *feiShuBot) UserIdMapByEmails(userIdType string, emails []string) ([]*larkcontact.UserContactInfo, error) {
	return m.getUserIdInfo(userIdType, nil, emails)
}
func (m *feiShuBot) UserInfoListByIds(userIdType string, idList []string, options ...larkcore.RequestOptionFunc) ([]*larkcontact.User, error) {
	if userIdType == "" || (userIdType != "open_id" && userIdType != "user_id" && userIdType != "union_id") {
		return nil, fmt.Errorf("userIdType is error: %s", userIdType)
	}
	// 创建请求对象
	req := larkcontact.NewBatchUserReqBuilder().
		UserIds(idList).
		UserIdType(userIdType).
		DepartmentIdType(`department_id`).
		Build()

	// 发起请求
	resp, err := m.client.Contact.V3.User.Batch(context.Background(), req, options...)

	// 处理错误
	if err != nil {
		return nil, err
	}

	// 服务端错误处理
	if !resp.Success() {
		return nil, fmt.Errorf("logId: %s, error response: \n%s", resp.RequestId(), larkcore.Prettify(resp.CodeError))
	}
	return resp.Data.Items, nil
}
