package cycy

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Receiver struct {
	Client       SQSReceiveMessageAPI //accept this interface for test mocking
	events       chan types.Message
	errors       chan error
	quit         chan struct{}
	cfg          ReceiverConfig
	handlers     MessageHandler
	errorHandler ErrorHandler
}
type MessageHandler interface {
	HandleMessage(types.Message, Receiver)
}
type ErrorHandler interface {
	HandleError(error, Receiver)
}

type MessageHandlerFunc func(types.Message, Receiver)

func (h MessageHandlerFunc) HandleMessage(msg types.Message, r Receiver) {
	h(msg, r)
}

type ErrorHandlerFunc func(error, Receiver)

func (h ErrorHandlerFunc) HandleError(err error, r Receiver) {
	h(err, r)
}

func NewReceiver(handler MessageHandler, opt ...func(config *ReceiverConfig)) *Receiver {
	cycy_cfg := &ReceiverConfig{}
	for _, cb := range opt {
		cb(cycy_cfg)
	}
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}
	client := sqs.NewFromConfig(cfg)
	rcvr := &Receiver{
		Client: client,
		cfg:    *cycy_cfg,
	}
	return rcvr
}

func (rcvr *Receiver) StartPolling() {

	//setup the handlers to process the messages
	go func() {
		var err error
		for err == nil {
			select {
			case msg := <-rcvr.events:
				log.Println("Message ID:     " + *msg.MessageId)
				log.Println("Message Handle: " + *msg.ReceiptHandle)
			case err = <-rcvr.errors:

			}
		}
	}()

	var sleep_time time.Duration
	for {
		input := &sqs.ReceiveMessageInput{
			MessageAttributeNames: []string{
				string(types.QueueAttributeNameAll),
			},
			QueueUrl:            aws.String(rcvr.cfg.SQS.URL),
			MaxNumberOfMessages: rcvr.cfg.Polling.MaxNumberOfMessages,
			VisibilityTimeout:   int32(rcvr.cfg.Timeout),
		}
		result, err := GetMessages(context.TODO(), rcvr.Client, input)
		if err != nil {
			log.Println("Got an error receiving messages:")
			log.Println(err)
			rcvr.errors <- err
			return
		}
		if result.Messages != nil {
			log.Printf("Received %d messages.\n", len(result.Messages))
			sleep_time = rcvr.cfg.Polling.Sleep
			for _, msg := range result.Messages {
				rcvr.events <- msg
			}

		} else {
			log.Println("No messages found")
			if sleep_time < rcvr.cfg.Polling.MaxSleepDuration {
				sleep_time += rcvr.cfg.Polling.DeadPollingSleepDelta
			}
		}
		log.Printf("Sleeping for: %s\n", sleep_time)
		time.Sleep(sleep_time)
	}
	//https://sqs.us-east-2.amazonaws.com/893635279084/quantum-data-publisher.fifo
}

// implemented by *sqs.Client from the aws go v2 sdk
type SQSReceiveMessageAPI interface {
	ReceiveMessage(ctx context.Context,
		params *sqs.ReceiveMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)
}

func GetMessages(c context.Context, api SQSReceiveMessageAPI, input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	return api.ReceiveMessage(c, input)
}
