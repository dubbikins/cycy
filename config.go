package cycy

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Config struct {
	Receiver *ReceiverConfig
	Sender   *SenderConfig
}

func NewConfig() *Config {
	cfg := &Config{}

	cfg.Receiver.MessageHandler = MessageHandlerFunc(func(m types.Message, r Receiver) {

	})
	cfg.Receiver.ErrorHandler = ErrorHandlerFunc(func(err error, r Receiver) {

	})
	return cfg
}

type SQSConfig struct {
	Name   string `env:"QUEUE_NAME"`
	Region string `env:"QUEUE_REGION"`
	ARN    string `env:"QUEUE_ARN"`
	URL    string `env:"QUEUE_URL"`
}

type S3Config struct {
	Name   string `env:"S3_NAME"`
	Region string `env:"S3_REGION"`
	ARN    string `env:"S3_ARN"`
}

type ReceiverConfig struct {
	SQS            *SQSConfig
	Timeout        int `env:"RECEIVER_TIMEOUT;default=5"`
	Polling        *PollingConfig
	MessageHandler MessageHandler
	ErrorHandler   ErrorHandler
}
type SenderConfig struct {
	S3 *S3Config
}

type PollingConfig struct {
	MaxNumberOfMessages   int32         `env:"RECEIVER_POLLING_MAX_MESSAGES;default=25"`
	Sleep                 time.Duration `env:"RECEIVER_POLLING_SLEEP;default=10s"`
	DeadPollingSleepDelta time.Duration `env:"RECEIVER_POLLING_SLEEP;default=0s"`
	MaxSleepDuration      time.Duration `env:"RECEIVER_POLLING_SLEEP;default=5m"`
}
