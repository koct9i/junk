package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/go-logr/logr"

	"github.com/urfave/cli/v3"

	"github.com/koct9i/junk/ytc/config"
	"github.com/koct9i/junk/ytc/cypress"
	"github.com/koct9i/junk/ytc/log"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var logger logr.Logger

	var timeout time.Duration
	var ctxCancel func()

	command := cli.Command{
		Flags: append([]cli.Flag{
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
		}, config.Flags()...),
		Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
			if timeout > 0 {
				ctx, ctxCancel = context.WithTimeout(ctx, timeout)
			}
			logger = log.NewLogger(os.Stderr)
			ctx = logr.NewContext(ctx, logger)
			if err := config.LoadConfig(); err != nil {
				return ctx, err
			}
			return ctx, nil
		},
		After: func(ctx context.Context, c *cli.Command) error {
			if ctxCancel != nil {
				ctxCancel()
			}
			return nil
		},
		Commands: []*cli.Command{
			cypress.Get(),
			cypress.Set(),
			cypress.List(),
		},
	}
	args := append(strings.Split(os.Args[0], "__"), os.Args[1:]...)
	if err := command.Run(ctx, args); err != nil {
		logger.Error(err, "Failed")
	}
}
