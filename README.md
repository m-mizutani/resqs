# ReSQS [![Travis-CI](https://travis-ci.com/m-mizutani/resqs.svg)](https://travis-ci.org/m-mizutani/resqs) [![Report card](https://goreportcard.com/badge/github.com/m-mizutani/resqs)](https://goreportcard.com/report/github.com/m-mizutani/resqs)

CLI tool to requeue messages in AWS SQS queue. A major use case is to requeue from DLQ (Dead Letter Queue) to original queue.

## Install

```
$ go get github.com/m-mizutani/resqs
```

## Usage

In basic usage, `resqs` retrieves message(s) from https://sqs.ap-northeast-1.amazonaws.com/1111111111/source-queue and sends them to https://sqs.ap-northeast-1.amazonaws.com/1111111111/destination-queue

```
$ resqs -s https://sqs.ap-northeast-1.amazonaws.com/1111111111/source-queue -d https://sqs.ap-northeast-1.amazonaws.com/1111111111/destination-queue
```

### Set limit of message number

```
$ resqs -m 100 -s https://sqs.ap-northeast-1.amazonaws.com/1111111111/source-queue -d https://sqs.ap-northeast-1.amazonaws.com/1111111111/destination-queue
```

`resqs` retrieves and sends only 100 messages by `-m 100` option.

## Test

```
go test .
```

## License

MIT

## Author

Masayoshi Mizutani <mizutani@hey.com>
