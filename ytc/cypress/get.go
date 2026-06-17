package cypress

import (
	"context"
	"os"

	"github.com/urfave/cli/v3"
	"go.ytsaurus.tech/yt/go/ypath"
)

func getCommand() *cli.Command {
	return &cli.Command{
		Name:      "get",
		Usage:     "get Cypress node value",
		ArgsUsage: "<path>",
		Flags:     flags(),
		Action: func(ctx context.Context, c *cli.Command) error {
			if err := requireArgs(c, 1, "get <path>"); err != nil {
				return err
			}
			yc, err := client()
			if err != nil {
				return err
			}

			var result any
			if err := yc.GetNode(ctx, ypath.Path(c.Args().Get(0)), &result, nil); err != nil {
				return err
			}
			return printValue(os.Stdout, result)
		},
	}
}
