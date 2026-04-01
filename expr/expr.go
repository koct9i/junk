package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/expr-lang/expr"
	"github.com/expr-lang/expr/file"
	"github.com/expr-lang/expr/parser"
	"github.com/expr-lang/expr/ast"
	"github.com/expr-lang/expr/parser/lexer"
	"github.com/urfave/cli/v3"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	app := cli.Command{
		Arguments: []cli.Argument{
			&cli.StringArg{
				Name: "expr",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {
			tokens, err := lexer.Lex(file.NewSource(c.StringArg("expr")))
			if err != nil {
				return err
			}
			for _, token := range tokens {
				fmt.Printf("%s\t%s\n", token.Kind, token.Value)
			}
			tree, err := parser.Parse(c.StringArg("expr"))
			fmt.Println(ast.Dump(tree.Node))
			output, err := expr.Eval(c.StringArg("expr"), nil)
			if err != nil {
				return err
			}
			return json.NewEncoder(os.Stdout).Encode(output)
		},
	}
	err := app.Run(ctx, os.Args)
	stop()
	if err != nil {
		log.Fatal(err)
	}
}
