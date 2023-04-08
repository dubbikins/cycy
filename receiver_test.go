package cycy

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SQSReceiveMessageImpl struct{}

func (dt SQSReceiveMessageImpl) ReceiveMessage(ctx context.Context,
	params *sqs.ReceiveMessageInput,
	optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error) {

	messages := []types.Message{
		{
			MessageId:     aws.String("aws-docs-example-message1-id"),
			ReceiptHandle: aws.String("aws-docs-example-message1-receipt-handle"),
		},
		{
			MessageId:     aws.String("aws-docs-example-message2-id"),
			ReceiptHandle: aws.String("aws-docs-example-message2-receipt-handle"),
		},
	}

	output := &sqs.ReceiveMessageOutput{
		Messages: messages,
	}

	return output, nil
}
func TestStartPolling(t *testing.T) {
	r := NewReceiver(func(m types.Message, r Receiver) {}, func(config *ReceiverConfig) {
		config = &ReceiverConfig{}
	})
	r.StartPolling()

}
