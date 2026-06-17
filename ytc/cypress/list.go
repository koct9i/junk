package cypress

import (
	"context"
	"os"

	"github.com/urfave/cli/v3"
	"go.ytsaurus.tech/yt/go/ypath"
)

func listCommand() *cli.Command {
	return &cli.Command{
		Name:      "list",
		Usage:     "list Cypress directory children",
		ArgsUsage: "<path>",
		Flags:     flags(),
		Action: func(ctx context.Context, c *cli.Command) error {
			if err := requireArgs(c, 1, "list <path>"); err != nil {
				return err
			}
			yc, err := client()
			if err != nil {
				return err
			}

			var result []string
			if err := yc.ListNode(ctx, ypath.Path(c.Args().Get(0)), &result, nil); err != nil {
				return err
			}
			return printValue(os.Stdout, result)
		},
	}
}
