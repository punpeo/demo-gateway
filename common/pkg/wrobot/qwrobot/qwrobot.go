package qwrobot

import (
	"bytes"
	"encoding/json"
	"net/http"
)

var contentMap = make(map[int64]string)

type Robot struct {
	Url      string
	Msgtype  string
	MarkDown MarkDown
	Text     Text
}
type Message struct {
	MsgType  string
	Text     Text
	MarkDown MarkDown
}

type Text struct {
	Content string
}
type MarkDown struct {
	Content string
}

func (r *Robot) SendMessageWithMarkDown() error {

	body := map[string]interface{}{
		"msgtype": "markdown",
		"markdown": map[string]string{
			"content": r.MarkDown.Content,
		},
	}
	// 发送消息到 Webhook
	messageBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	requestBody := bytes.NewReader(messageBytes)
	resp, err := http.Post(r.Url, "application/json", requestBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
func (r *Robot) SendMessageWithText(url, message string) error {
	webhookURL := url
	body := map[string]interface{}{
		"msgtype": "text",
		"text": map[string]string{
			"content": message,
		},
	}
	// 发送消息到 Webhook
	messageBytes, err := json.Marshal(body)
	if err != nil {
		return err
	}

	requestBody := bytes.NewReader(messageBytes)
	resp, err := http.Post(webhookURL, "application/json", requestBody)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
