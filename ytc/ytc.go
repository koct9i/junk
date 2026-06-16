package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/go-logr/logr"

	"github.com/urfave/cli/v3"

	"github.com/koct9i/junk/ytc/log"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var logger logr.Logger

	var timeout time.Duration
	var ctxCancel func()

	command := cli.Command{
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "log-verbosity",
				Aliases:     []string{"v"},
				Destination: &log.Verbosity,
			},
			&cli.BoolWithInverseFlag{
				Name:        "log-pretty",
				Value:       true,
				Destination: &log.Pretty,
			},
			&cli.DurationFlag{
				Name:        "timeout",
				Aliases:     []string{"t"},
				Destination: &timeout,
			},
		},
		Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
			if timeout > 0 {
				ctx, ctxCancel = context.WithTimeout(ctx, timeout)
			}
			logger = log.NewLogger(os.Stdout)
			ctx = logr.NewContext(ctx, logger)
			return ctx, nil
		},
		After: func(ctx context.Context, c *cli.Command) error {
			if ctxCancel != nil {
				ctxCancel()
			}
			return nil
		},
		Commands: []*cli.Command{},
	}
	args := append(strings.Split(os.Args[0], "__"), os.Args[1:]...)
	if err := command.Run(ctx, args); err != nil {
		logger.Error(err, "Failed")
	}
}
