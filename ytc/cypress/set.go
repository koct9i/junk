package cypress

import (
	"context"

	"github.com/urfave/cli/v3"
	"go.ytsaurus.tech/yt/go/ypath"
)

func setCommand() *cli.Command {
	return &cli.Command{
		Name:      "set",
		Usage:     "set Cypress node value from JSON argument or stdin",
		ArgsUsage: "<path> [json-value]",
		Flags:     flags(),
		Action: func(ctx context.Context, c *cli.Command) error {
			if c.Args().Len() < 1 || c.Args().Len() > 2 {
				return requireArgs(c, 2, "set <path> [json-value]")
			}
			yc, err := client()
			if err != nil {
				return err
			}
			value, err := readValue(c)
			if err != nil {
				return err
			}
			return yc.SetNode(ctx, ypath.Path(c.Args().Get(0)), value, nil)
		},
	}
}
