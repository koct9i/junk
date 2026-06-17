package cypress

import (
	"context"
	"os"

	"github.com/urfave/cli/v3"
	"go.ytsaurus.tech/yt/go/ypath"
)

func getCommand() *cli.Command {
	return &cli.Command{
		Name:  "get",
		Usage: "get Cypress node value",
		Arguments: []cli.Argument{
			&cli.StringArgs{Name: "path", Min: 1, Max: 1},
		},
		Flags: flags(),
		Action: func(ctx context.Context, c *cli.Command) error {
			yc, err := client()
			if err != nil {
				return err
			}

			var result any
			if err := yc.GetNode(ctx, ypath.Path(c.StringArgs("path")[0]), &result, nil); err != nil {
				return err
			}
			return printValue(os.Stdout, result)
		},
	}
}
