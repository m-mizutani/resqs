package service

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/m-mizutani/resqs/pkg/adaptor"
	"github.com/m-mizutani/resqs/pkg/errors"
	"github.com/m-mizutani/resqs/pkg/logging"
	"github.com/sirupsen/logrus"
)

var logger = logging.Logger

type RequeueOptions struct {
	Adaptors     *adaptor.Adaptors
	MessageLimit int
}

func extractSQSRegion(url string) (string, error) {
	// QueueURL sample: https://sqs.eu-west-2.amazonaws.com/
	urlParts := strings.Split(url, "/")
	if len(urlParts) < 3 {
		return "", errors.New("Not enough slash in queue URL")
	}
	domainParts := strings.Split(urlParts[2], ".")
	if len(domainParts) != 4 {
		return "", errors.New("Not enough dot in queue URL")
	}

	return domainParts[1], nil
}

func Requeue(srcQueueURL, dstQueueURL string) error {
	return RequeueWithOpt(srcQueueURL, dstQueueURL, &RequeueOptions{})
}

func RequeueWithOpt(srcQueueURL, dstQueueURL string, opt *RequeueOptions) error {
	logger.WithFields(logrus.Fields{
		"src": srcQueueURL,
		"dst": dstQueueURL,
	}).Info("Start requeuing")

	newSQS := opt.Adaptors.NewSQS
	if newSQS == nil {
		newSQS = adaptor.NewSQSClient
	}

	srcRegion, err := extractSQSRegion(srcQueueURL)
	if err != nil {
		return errors.Wrap(err, "Invalid srcQueueURL").With("url", srcQueueURL)
	}
	dstRegion, err := extractSQSRegion(dstQueueURL)
	if err != nil {
		return errors.Wrap(err, "Invalid dstQueueURL").With("url", srcQueueURL)
	}

	srcClient, err := newSQS(srcRegion)
	if err != nil {
		return errors.Wrap(err, "Failed to create new SQS client for source queue")
	}
	dstClient, err := newSQS(dstRegion)
	if err != nil {
		return errors.Wrap(err, "Failed to create new SQS client for destination queue")
	}

	msgCount := 0

	for {
		recvOut, err := srcClient.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl: aws.String(srcQueueURL),
		})
		if err != nil {
			return errors.Wrap(err, "Failed to receive message from src queue")
		}
		if len(recvOut.Messages) == 0 {
			logger.Info("No available message in src queue")
			break
		}
		logger.WithField("len(msg)", len(recvOut.Messages)).Debug("Got message(s) from src queue")

		for _, msg := range recvOut.Messages {
			msgCount++
			if 0 < opt.MessageLimit && opt.MessageLimit < msgCount {
				logger.WithField("count", msgCount).Info("Exit loop by hitting limit")
				return nil
			}

			logger.WithField("message", *msg.Body).Debug("Sending a message to dst queue")

			_, err := dstClient.SendMessage(&sqs.SendMessageInput{
				QueueUrl:    aws.String(dstQueueURL),
				MessageBody: msg.Body,
			})
			if err != nil {
				return errors.Wrap(err, "Failed to send message to dst queue")
			}

			_, err = srcClient.DeleteMessage(&sqs.DeleteMessageInput{
				QueueUrl:      aws.String(srcQueueURL),
				ReceiptHandle: msg.ReceiptHandle,
			})
			if err != nil {
				return errors.Wrap(err, "Failed to delete message from src queue").
					With("handle", *msg.ReceiptHandle)
			}
			logger.WithField("receipt", msg.ReceiptHandle).Debug("Remove a message")
		}
	}

	return nil
}
