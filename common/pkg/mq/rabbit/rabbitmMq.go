package rabbit

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"github.com/zeromicro/go-queue/rabbitmq"
)

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
	// 队列名称
	QueueName string
	// 交换机
	Exchange string
	// routing Key
	RoutingKey string
	//MQ链接字符串
	Mqurl string
}

// 创建结构体实例
func NewRabbitMQ(queueName string, exchange string, rabbitConf rabbitmq.RabbitConf, routingKey string) (mq *RabbitMQ, err error) {
	rabbitMQ := RabbitMQ{
		QueueName:  queueName,
		Exchange:   exchange,
		RoutingKey: routingKey,
		Mqurl:      getRabbitURL(rabbitConf),
	}

	//创建rabbitmq连接
	rabbitMQ.Conn, err = amqp.Dial(rabbitMQ.Mqurl)
	if err != nil {
		return nil, err
	}

	//创建Channel
	rabbitMQ.Channel, err = rabbitMQ.Conn.Channel()
	if err != nil {
		return nil, err
	}

	return &rabbitMQ, err

}

// 使用map结构发送延迟消息
func (mq *RabbitMQ) DirectSendDelayMessageQueue(messages []map[string]interface{}, delaySec int64) (err error) {
	var mapMessageToStringMessage []string

	for _, message := range messages {
		textMessage, _ := json.Marshal(message)
		mapMessageToStringMessage = append(mapMessageToStringMessage, string(textMessage))
	}

	return mq.DirectSendDelayMessageQueueByTextMessage(mapMessageToStringMessage, delaySec)
}

// 直接发送消息
func (mq *RabbitMQ) DirectSendMessageQueueByTextMessage(messages []string) (err error) {
	for _, message := range messages {
		var err = mq.Channel.Publish("", mq.QueueName, false, false, amqp.Publishing{
			ContentType:  "application/json",
			Body:         []byte(message),
			Headers:      map[string]interface{}{},
			DeliveryMode: 2,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// 使用string text message结构发送延迟消息
func (mq *RabbitMQ) DirectSendDelayMessageQueueByTextMessage(messages []string, delaySec int64) (err error) {
	exchangeName := mq.QueueName + "_Delay"
	_, err = mq.Channel.QueueDeclare(mq.QueueName, true, false, false, false, map[string]interface{}{})
	if err != nil {
		return err
	}

	err = mq.Channel.ExchangeDeclare(exchangeName, "x-delayed-message", true, false, false, false, map[string]interface{}{
		"x-delayed-types": "direct",
	})
	if err != nil {
		return err
	}

	err = mq.Channel.QueueBind(mq.QueueName, mq.QueueName, exchangeName, false, map[string]interface{}{})
	if err != nil {
		return err
	}

	for _, message := range messages {
		var err = mq.Channel.Publish(exchangeName, mq.QueueName, false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(message),
			Headers: map[string]interface{}{
				"x-delay": int(delaySec * 1000),
			},
			DeliveryMode: 2,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// 广播消息
func (mq RabbitMQ) Broadcast(exchangeName string, message string) error {
	err := mq.Channel.ExchangeDeclare(exchangeName, "fanout", false, false, false, false, map[string]interface{}{})
	if err != nil {
		return err
	}
	err = mq.Channel.Publish(exchangeName, "rk_"+exchangeName, false, false, amqp.Publishing{ContentType: "application/json", Body: []byte(message)})
	if err != nil {
		return err
	}
	return nil
}
func (mq RabbitMQ) SendByRouting(exchangeName string, routingKey string, message string) error {
	err := mq.Channel.ExchangeDeclare(exchangeName, "direct", false, false, false, false, map[string]interface{}{})
	if err != nil {
		return err
	}
	err = mq.Channel.Publish(exchangeName, routingKey, false, false, amqp.Publishing{ContentType: "application/json", Body: []byte(message)})
	if err != nil {
		return err
	}
	return nil
}

// 释放资源,建议NewRabbitMQ获取实例后 配合defer使用
func (mq *RabbitMQ) ReleaseRes() {
	err := mq.Conn.Close()
	if err != nil {
		return
	}
	err = mq.Channel.Close()
	if err != nil {
		return
	}
}
