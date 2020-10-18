package main

import (
	"bytes"
	"testing"

	"github.com/m-mizutani/resqs/pkg/adaptor"
	"github.com/m-mizutani/resqs/pkg/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	srcQueue = "https://sqs.ap-northeast-1.amazonaws.com/1111111111/src-queue"
	dstQueue = "https://sqs.us-east-1.amazonaws.com/222222222/dst-queue"
)

func TestBasicUsage(t *testing.T) {
	sqsClient, newSQS := mock.NewMockSQS()
	adaptors := &adaptor.Adaptors{
		NewSQS: newSQS,
	}
	src := sqsClient.Queues.Get(srcQueue)
	dst := sqsClient.Queues.Get(dstQueue)
	src.Messages = []string{"a", "b", "c"}

	require.NoError(t, newApp(adaptors).
		Run([]string{"resqs", "-s", srcQueue, "-d", dstQueue}))

	t.Run("Should retrieve 4 times (3 messages + 1 empty response)", func(t *testing.T) {
		// 3 + 1 (empty response)
		assert.Equal(t, 4, len(sqsClient.RecvInput))
		for _, input := range sqsClient.RecvInput {
			assert.Equal(t, srcQueue, *input.QueueUrl)
		}
	})

	t.Run("Should send and delete 3 times to dstQueue", func(t *testing.T) {
		// 3
		require.Equal(t, 3, len(sqsClient.SendInput))
		require.Equal(t, 3, len(sqsClient.DeleteInput))
		for _, input := range sqsClient.SendInput {
			assert.Equal(t, dstQueue, *input.QueueUrl)
		}
		for _, input := range sqsClient.DeleteInput {
			assert.Equal(t, srcQueue, *input.QueueUrl)
		}
	})

	t.Run("Should match sent messages with received messages", func(t *testing.T) {
		require.Equal(t, 3, len(dst.Messages))
		assert.Equal(t, "a", dst.Messages[0])
		assert.Equal(t, "b", dst.Messages[1])
		assert.Equal(t, "c", dst.Messages[2])
	})

}

func TestFailureCases(t *testing.T) {
	t.Run("Should fail if missing src queue", func(t *testing.T) {
		app := newApp(&adaptor.Adaptors{})
		app.Writer = &bytes.Buffer{}
		require.Error(t, app.Run([]string{"resqs", "-d", dstQueue}))
	})

	t.Run("Should fail if missing dst queue", func(t *testing.T) {
		app := newApp(&adaptor.Adaptors{})
		app.Writer = &bytes.Buffer{}
		require.Error(t, app.Run([]string{"resqs", "-s", srcQueue}))
	})
}

func TestMessageNumberLimit(t *testing.T) {
	t.Run("Should sent only 4 messages by -m 4 option", func(t *testing.T) {
		sqsClient, newSQS := mock.NewMockSQS()
		adaptors := &adaptor.Adaptors{NewSQS: newSQS}
		src := sqsClient.Queues.Get(srcQueue)
		dst := sqsClient.Queues.Get(dstQueue)
		src.Messages = []string{"a", "b", "c", "d", "e"}

		require.NoError(t, newApp(adaptors).
			Run([]string{"resqs", "-s", srcQueue, "-d", dstQueue, "-m", "4"}))

		require.Equal(t, 4, len(dst.Messages))
		assert.Equal(t, "a", dst.Messages[0])
		assert.Equal(t, "b", dst.Messages[1])
		assert.Equal(t, "c", dst.Messages[2])
		assert.Equal(t, "d", dst.Messages[3])
	})

	t.Run("Should sent all messages if a number of message is lesser than -m option", func(t *testing.T) {
		sqsClient, newSQS := mock.NewMockSQS()
		adaptors := &adaptor.Adaptors{NewSQS: newSQS}
		src := sqsClient.Queues.Get(srcQueue)
		dst := sqsClient.Queues.Get(dstQueue)
		src.Messages = []string{"a", "b", "c", "d", "e"}

		require.NoError(t, newApp(adaptors).
			Run([]string{"resqs", "-s", srcQueue, "-d", dstQueue, "-m", "6"}))

		require.Equal(t, 5, len(dst.Messages))
		assert.Equal(t, "a", dst.Messages[0])
		assert.Equal(t, "b", dst.Messages[1])
		assert.Equal(t, "c", dst.Messages[2])
		assert.Equal(t, "d", dst.Messages[3])
		assert.Equal(t, "e", dst.Messages[4])
	})
}
