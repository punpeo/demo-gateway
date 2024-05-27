package rabbit

type SenderConfig struct {
	ExchangeName string
	QueueName    string
	RouterKey    string
	ContentType  string
}
