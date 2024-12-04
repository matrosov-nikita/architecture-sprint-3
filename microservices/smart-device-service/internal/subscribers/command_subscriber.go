package subscribers

import "github.com/segmentio/kafka-go"

const (
	deviceCommandsTopic = "device_commands"
	groupID             = "device-consumer-group"
)

type CommandSubscriber struct {
	kafkaReader *kafka.Reader
}

func NewCommandSubscriber(kafkaReader *kafka.Reader) *CommandSubscriber {
	return &CommandSubscriber{kafkaReader: kafkaReader}
}
