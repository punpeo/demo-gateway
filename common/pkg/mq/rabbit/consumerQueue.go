package rabbit

import "github.com/zeromicro/go-queue/rabbitmq"

type ConsumerConfig struct {
	ListenerQueues []rabbitmq.ConsumerConf
	Consumers      int `json:",default=1"`
}
