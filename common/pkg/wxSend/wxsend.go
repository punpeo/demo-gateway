package wxSend

import (
	errorx "github.com/punpeo/punpeo-lib/rest/jerror"
	"github.com/wxpusher/wxpusher-sdk-go"
	"github.com/wxpusher/wxpusher-sdk-go/model"
)

type SendMsg struct {
	Uid       string `json:"uid"`
	TopicId   int    `json:"topicId"`
	MessageId int    `json:"messageId"`
	Code      int    `json:"code"`
	Status    string `json:"status"`
	Message   string `json:"message"`
}

// GetWxUser 获取用户id列表
func GetWxUser(appToken string, page, pageSize int) ([]string, error) {
	result, err := wxpusher.QueryWxUser(appToken, page, pageSize)
	if err != nil {
		return nil, err
	}
	uniqueIDs := make(map[string]bool)
	for _, v := range result.Records {
		uniqueIDs[v.UId] = true
	}

	ids := make([]string, 0, len(uniqueIDs))
	for id := range uniqueIDs {
		ids = append(ids, id)
	}
	return ids, nil
}

// SendByIds  用户
func SendByIds(appToken string, content string, uid []string) (err error) {
	message := model.Message{
		AppToken:    appToken,
		ContentType: 1,
	}
	msg := message.SetContent(content).AddUId(uid[0], uid...)
	msgArr, err := wxpusher.SendMessage(msg)
	if err != nil {
		return err
	}
	var sendMsg []SendMsg
	for _, v := range msgArr {
		if v.Code == 1000 {
			sendMsg = append(sendMsg, SendMsg{
				Uid:       v.Uid,
				TopicId:   v.TopicId,
				MessageId: v.MessageId,
				Code:      v.Code,
				Status:    v.Status,
				Message:   content,
			})
		} else {

			return errorx.NewDefaultError(v.Status)
		}
	}
	return nil
}
func SendByTopic(appToken string, content string, topicId []int) (err error) {
	message := model.Message{
		AppToken:    appToken,
		ContentType: 1,
	}
	msg := message.SetContent(content).AddTopicId(topicId[0], topicId...)
	msgArr, err := wxpusher.SendMessage(msg)
	if err != nil {
		return err
	}
	var sendMsg []SendMsg
	for _, v := range msgArr {
		if v.Code == 1000 {
			sendMsg = append(sendMsg, SendMsg{
				Uid:       v.Uid,
				TopicId:   v.TopicId,
				MessageId: v.MessageId,
				Code:      v.Code,
				Status:    v.Status,
				Message:   content,
			})
		} else {
			return errorx.NewDefaultError(v.Status)
		}
	}
	return nil
}
