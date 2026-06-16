package cypress

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/urfave/cli/v3"
	"go.ytsaurus.tech/yt/go/yson"
	ytsdk "go.ytsaurus.tech/yt/go/yt"
	"go.ytsaurus.tech/yt/go/yt/ythttp"
	"gopkg.in/yaml.v3"
)

var (
	proxy  string
	format string
)

// Commands returns the minimal Cypress command set.
func Commands() []*cli.Command {
	return []*cli.Command{
		getCommand(),
		setCommand(),
		listCommand(),
	}
}

func flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "proxy",
			Usage:       "YT HTTP proxy; defaults to YT_PROXY when empty",
			Sources:     cli.EnvVars("YT_PROXY"),
			Destination: &proxy,
		},
		&cli.StringFlag{
			Name:        "format",
			Aliases:     []string{"f"},
			Usage:       "data format: json, yson, or yaml",
			Value:       "json",
			Destination: &format,
		},
	}
}

func client() (ytsdk.Client, error) {
	return ythttp.NewClient(&ytsdk.Config{Proxy: proxy})
}

func requireArgs(c *cli.Command, n int, usage string) error {
	if c.Args().Len() != n {
		return fmt.Errorf("usage: %s", usage)
	}
	return nil
}

func printValue(w io.Writer, v any) error {
	switch normalizedFormat() {
	case "json":
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		return encoder.Encode(v)
	case "yson":
		data, err := yson.MarshalFormat(v, yson.FormatPretty)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(w, string(data))
		return err
	case "yaml", "yml":
		data, err := yaml.Marshal(v)
		if err != nil {
			return err
		}
		_, err = w.Write(data)
		return err
	default:
		return fmt.Errorf("unsupported format %q", format)
	}
}

func parseValue(s string) (any, error) {
	var value any
	s = strings.TrimSpace(s)

	switch normalizedFormat() {
	case "json":
		if err := json.Unmarshal([]byte(s), &value); err != nil {
			return s, nil
		}
		return foldYSONNodes(value)
	case "yson":
		if err := yson.Unmarshal([]byte(s), &value); err != nil {
			return nil, err
		}
	case "yaml", "yml":
		if err := yaml.Unmarshal([]byte(s), &value); err != nil {
			return nil, err
		}
		return foldYSONNodes(value)
	default:
		return nil, fmt.Errorf("unsupported format %q", format)
	}

	return value, nil
}

func foldYSONNodes(value any) (any, error) {
	switch typed := value.(type) {
	case []any:
		for i, item := range typed {
			folded, err := foldYSONNodes(item)
			if err != nil {
				return nil, err
			}
			typed[i] = folded
		}
		return typed, nil
	case map[string]any:
		_, hasValue := typed["$value"]
		_, hasAttrs := typed["$attributes"]
		if hasValue || hasAttrs {
			return foldYSONNode(typed, hasValue, hasAttrs)
		}
		for key, item := range typed {
			folded, err := foldYSONNodes(item)
			if err != nil {
				return nil, err
			}
			typed[key] = folded
		}
		return typed, nil
	case map[any]any:
		converted := make(map[string]any, len(typed))
		for key, item := range typed {
			keyString, ok := key.(string)
			if !ok {
				return nil, fmt.Errorf("unsupported non-string map key %v", key)
			}
			converted[keyString] = item
		}
		return foldYSONNodes(converted)
	default:
		return value, nil
	}
}

func foldYSONNode(node map[string]any, hasValue, hasAttrs bool) (any, error) {
	foldedValue, err := foldYSONNodeValue(node, hasValue)
	if err != nil {
		return nil, err
	}

	attrs := map[string]any(nil)
	if hasAttrs {
		attrsMap, err := stringMap(node["$attributes"])
		if err != nil {
			return nil, fmt.Errorf("$attributes must be an object: %w", err)
		}
		attrs = make(map[string]any, len(attrsMap))
		for key, item := range attrsMap {
			folded, err := foldYSONNodes(item)
			if err != nil {
				return nil, err
			}
			attrs[key] = folded
		}
	}

	return &yson.ValueWithAttrs{Attrs: attrs, Value: foldedValue}, nil
}

func foldYSONNodeValue(node map[string]any, hasValue bool) (any, error) {
	if hasValue {
		for key := range node {
			if key != "$value" && key != "$attributes" {
				return nil, fmt.Errorf("unexpected key %q alongside $value", key)
			}
		}
		return foldYSONNodes(node["$value"])
	}

	value := make(map[string]any, len(node)-1)
	for key, item := range node {
		if key == "$attributes" {
			continue
		}
		value[key] = item
	}
	if len(value) == 0 {
		return nil, nil
	}
	return foldYSONNodes(value)
}

func stringMap(value any) (map[string]any, error) {
	switch typed := value.(type) {
	case map[string]any:
		return typed, nil
	case map[any]any:
		converted := make(map[string]any, len(typed))
		for key, item := range typed {
			keyString, ok := key.(string)
			if !ok {
				return nil, fmt.Errorf("unsupported non-string map key %v", key)
			}
			converted[keyString] = item
		}
		return converted, nil
	default:
		return nil, fmt.Errorf("got %T", value)
	}
}

func readValue(c *cli.Command) (any, error) {
	if c.Args().Len() >= 2 {
		return parseValue(c.Args().Get(1))
	}
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, err
	}
	return parseValue(string(data))
}

func normalizedFormat() string {
	return strings.ToLower(strings.TrimSpace(format))
}
