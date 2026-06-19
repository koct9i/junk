package cypress

import (
	"context"
	"os"

	"github.com/urfave/cli/v3"
	"go.ytsaurus.tech/yt/go/ypath"
)

func List() *cli.Command {
	return &cli.Command{
		Name:  "list",
		Usage: "list Cypress directory children",
		Arguments: []cli.Argument{
			&cli.StringArgs{Name: "path", Min: 1, Max: 1},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			yc, err := client(ctx)
			if err != nil {
				return err
			}

			var result []string
			if err := yc.ListNode(ctx, ypath.Path(c.StringArgs("path")[0]), &result, nil); err != nil {
				return err
			}
			return printValue(os.Stdout, result)
		},
	}
}
