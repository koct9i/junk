package cypress

import (
	"context"

	"github.com/urfave/cli/v3"
	"go.ytsaurus.tech/yt/go/ypath"
)

func Set() *cli.Command {
	return &cli.Command{
		Name:  "set",
		Usage: "set Cypress node value from argument or stdin",
		Arguments: []cli.Argument{
			&cli.StringArgs{Name: "path", Min: 1, Max: 1},
			&cli.StringArgs{Name: "value", Min: 0, Max: 1},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			yc, err := client(ctx)
			if err != nil {
				return err
			}
			value, err := readValue(c.StringArgs("value"))
			if err != nil {
				return err
			}
			return yc.SetNode(ctx, ypath.Path(c.StringArgs("path")[0]), value, nil)
		},
	}
}
