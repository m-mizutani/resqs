package adaptor

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/m-mizutani/resqs/pkg/errors"
)

// SQSClientFactory is interface SQSClient constructor
type SQSClientFactory func(region string) (SQSClient, error)

// SQSClient is interface of AWS SDK SQS
type SQSClient interface {
	SendMessage(*sqs.SendMessageInput) (*sqs.SendMessageOutput, error)
	ReceiveMessage(input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error)
	DeleteMessage(input *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error)
}

// NewSQSClient creates actual AWS SQS SDK client
func NewSQSClient(region string) (SQSClient, error) {
	ssn, err := session.NewSession(&aws.Config{Region: aws.String(region)})
	if err != nil {
		return nil, errors.Wrap(err, "Failed session.NewSession for NewSQSClient")
	}
	return sqs.New(ssn), nil
}

/*
type DryrunSQSClient struct {
	client *sqs.SQS
}

func NewDryrunSQSClient(region string) SQSClient {
	ssn := session.New(&aws.Config{Region: aws.String(region)})
	return &DryrunSQSClient{
		client: sqs.New(ssn),
	}
}

func (x *DryrunSQSClient) SendMessage(*sqs.SendMessageInput) (*sqs.SendMessageOutput, error) {
	return &sqs.SendMessageOutput{}, nil
}
func (x *DryrunSQSClient) ReceiveMessage(input *sqs.ReceiveMessageInput) (*sqs.ReceiveMessageOutput, error) {
	return x.client.D
}
func (x *DryrunSQSClient) DeleteMessage(input *sqs.DeleteMessageInput) (*sqs.DeleteMessageOutput, error) {
	return nil, nil
}
*/
