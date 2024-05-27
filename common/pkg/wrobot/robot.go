package wrobot

import (
	"study/common/pkg/wrobot/qwrobot"
	"time"
)

func SendMessageWithMarkDown(url, message string) error {
	robot := qwrobot.Robot{
		Url:     url,
		Msgtype: "markdown",
		MarkDown: qwrobot.MarkDown{
			Content: "**刘夏提醒，您有一条消息需要您处理**\n" + "时间：" + time.Now().Format("2006-01-02 15:04:05") + "\n" + "内容：" + message,
		},
	}
	return robot.SendMessageWithMarkDown()
}
func SendMessageWithText(url, message string) error {
	robot := qwrobot.Robot{
		Url:     url,
		Msgtype: "text",
		Text: qwrobot.Text{
			Content: message,
		},
	}
	return robot.SendMessageWithText(url, message)
}
