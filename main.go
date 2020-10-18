package main

import (
	"os"

	"github.com/m-mizutani/resqs/pkg/adaptor"
	"github.com/m-mizutani/resqs/pkg/errors"
	"github.com/m-mizutani/resqs/pkg/logging"
	"github.com/m-mizutani/resqs/pkg/service"
	cli "github.com/urfave/cli/v2"
)

var logger = logging.Logger

func main() {
	app := newApp(&adaptor.Adaptors{})
	if err := app.Run(os.Args); err != nil {
		entry := logger.WithError(err)
		if e, ok := err.(*errors.Error); ok {
			for k, v := range e.Values {
				entry = entry.WithField(k, v)
			}
		}
		entry.Fatal("Exit with error")
	}
}

func newApp(adaptors *adaptor.Adaptors) *cli.App {
	var srcQueue, dstQueue string
	var logLevel string
	opt := &service.RequeueOptions{
		Adaptors: adaptors,
	}

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "src-queue",
				Aliases:     []string{"s"},
				Destination: &srcQueue,
				Required:    true,
				Usage:       "Source Queue URL",
			},
			&cli.StringFlag{
				Name:        "dst-queue",
				Aliases:     []string{"d"},
				Destination: &dstQueue,
				Required:    true,
				Usage:       "Destination Queue URL",
			},
			&cli.StringFlag{
				Name:        "log-level",
				Aliases:     []string{"l"},
				Destination: &logLevel,
				Usage:       "Log Level [DEBUG|INFO|WARN|ERROR]",
				Value:       "INFO",
			},
			&cli.IntFlag{
				Name:        "message-limit",
				Aliases:     []string{"m"},
				Destination: &opt.MessageLimit,
				Usage:       "Limit of message",
			},
		},
		Action: func(c *cli.Context) error {
			logging.SetLogLevel(logLevel)

			if err := service.RequeueWithOpt(srcQueue, dstQueue, opt); err != nil {
				return err
			}
			return nil
		},
	}

	return app
}
