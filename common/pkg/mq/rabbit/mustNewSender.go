package rabbit

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"github.com/zeromicro/go-queue/rabbitmq"
	"github.com/zeromicro/go-zero/core/logx"
)

type RabbitMqSender struct {
	conn        *amqp.Connection
	channel     *amqp.Channel
	ContentType string
}

func getRabbitURL(rabbitConf rabbitmq.RabbitConf) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/%s", rabbitConf.Username, rabbitConf.Password,
		rabbitConf.Host, rabbitConf.Port, rabbitConf.VHost)
}

var MqSender rabbitmq.Sender

func MustNewSender(rabbitMqConf rabbitmq.RabbitSenderConf) (rabbitmq.Sender, error) {
	if MqSender == nil {
		sender, err2 := InitMqConnect(rabbitMqConf)
		if err2 != nil {
			return nil, err2
		}
		MqSender = sender
	}
	if s, ok := MqSender.(*RabbitMqSender); ok {
		if s.conn.IsClosed() {
			sender, err2 := InitMqConnect(rabbitMqConf)
			if err2 != nil {
				return nil, err2
			}
			MqSender = sender
		}
	}
	return MqSender, nil
}

func InitMqConnect(rabbitMqConf rabbitmq.RabbitSenderConf) (*RabbitMqSender, error) {
	sender := &RabbitMqSender{ContentType: rabbitMqConf.ContentType}
	conn, err := amqp.Dial(getRabbitURL(rabbitMqConf.RabbitConf))
	if err != nil {
		logx.Error(fmt.Sprintf("rabbitSender_conn_err : %v", err))
		return nil, err
	}

	sender.conn = conn
	channel, err := sender.conn.Channel()
	if err != nil {
		logx.Error(fmt.Sprintf("rabbitSender_channel_err : %v", err))
		return nil, err
	}

	sender.channel = channel
	return sender, nil
}

func (q *RabbitMqSender) Send(exchange string, routeKey string, msg []byte) error {
	return q.channel.Publish(
		exchange,
		routeKey,
		false,
		false,
		amqp.Publishing{
			ContentType: q.ContentType,
			Body:        msg,
		},
	)
}

func Publish(body any, rabbitConf rabbitmq.RabbitConf, senderQueue SenderConfig) error {
	defer func() {
		if err := recover(); err != nil {
			logx.Info(fmt.Sprintf(" rabbitSender_publish_err body:%v,err:%v", body, err))
		}
	}()

	conf := rabbitmq.RabbitSenderConf{RabbitConf: rabbitConf, ContentType: senderQueue.ContentType}
	sender, err := MustNewSender(conf)
	if err != nil {
		return err
	}

	msg, err := json.Marshal(body)
	if err != nil {
		return err
	}

	err = sender.Send(senderQueue.ExchangeName, senderQueue.RouterKey, msg)

	return err
}

func (q *RabbitMqSender) Close() error {
	return q.conn.Close()
}
