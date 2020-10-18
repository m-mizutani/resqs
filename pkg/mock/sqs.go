package mock

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/m-mizutani/resqs/pkg/adaptor"
)

// SQSClient is mock of adaptor.SQSClient
type SQSClient struct {
	SendInput   []*sqs.SendMessageInput
	RecvInput   []*sqs.ReceiveMessageInput
	DeleteInput []*sqs.DeleteMessageInput
	Queues      queueMap
}

type queueMap map[string]*queue

func (x queueMap) Get(url string) *queue {
	q, ok := x[url]
	if !ok {
		q = &queue{}
		x[url] = q
	}
	return q
}

type queue struct {
	Messages   []string
	MessagePtr int
}

// NewMockSQS returns mock.SQSClient and factory to return the SQSClient
func NewMockSQS() (*SQSClient, adaptor.SQSClientFactory) {
	client := &SQSClient{
		Queues: make(queueMap),
	}
	return client, func(region string) (adaptor.SQSClient, error) {
		return client, nil
	}
}

// SendMessage stores input to SendInput and message to Messages
func (x *SQSClient) SendMessage(input *sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	x.SendInput = append(x.SendInput, input)
	q := x.Queues.Get(*input.QueueUrl)
	q.Messages = append(q.Messages, *input.MessageBody)
	return &sqs.SendMessageOutput{}, nil
}

// ReceiveMessage stores input to RecvInput and returns message from Messages
func (x *SQSClient) ReceiveMessage(input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	// NOTE: This behaviour is not enough to emulate actual SQS
	x.RecvInput = append(x.RecvInput, input)
	q := x.Queues.Get(*input.QueueUrl)
	if q.MessagePtr < len(q.Messages) {
		output := &sqs.ReceiveMessageOutput{
			Messages: []*sqs.Message{
				{
					Body: aws.String(q.Messages[q.MessagePtr]),
				},
			},
		}
		q.MessagePtr++
		return output, nil
	}

	return &sqs.ReceiveMessageOutput{}, nil
}

// DeleteMessage stores input to DeleteInput
func (x *SQSClient) DeleteMessage(input *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	x.DeleteInput = append(x.DeleteInput, input)
	return &sqs.DeleteMessageOutput{}, nil
}
