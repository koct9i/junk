package cypress

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/go-logr/logr"
	"gopkg.in/yaml.v3"

	"go.ytsaurus.tech/yt/go/yson"
	ytsdk "go.ytsaurus.tech/yt/go/yt"
	"go.ytsaurus.tech/yt/go/yt/ythttp"

	"github.com/koct9i/junk/ytc/config"
	"github.com/koct9i/junk/ytc/log"
)

func client(ctx context.Context) (ytsdk.Client, error) {
	config := config.Config
	config.Logger = log.NewYTSlog(logr.FromContextAsSlogLogger(ctx))
	return ythttp.NewClient(&config)
}

func printValue(w io.Writer, v any) error {
	switch config.Format {
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
		return fmt.Errorf("unsupported format %q", config.Format)
	}
}

func parseValue(s string) (any, error) {
	var value any
	s = strings.TrimSpace(s)

	switch config.Format {
	case "json":
		if err := json.Unmarshal([]byte(s), &value); err != nil {
			return s, nil
		}
		return foldYSONAttributes(value)
	case "yson":
		if err := yson.Unmarshal([]byte(s), &value); err != nil {
			return nil, err
		}
	case "yaml", "yml":
		if err := yaml.Unmarshal([]byte(s), &value); err != nil {
			return nil, err
		}
		return foldYSONAttributes(value)
	default:
		return nil, fmt.Errorf("unsupported format %q", config.Format)
	}

	return value, nil
}

func foldYSONAttributes(value any) (any, error) {
	switch typed := value.(type) {
	case []any:
		for i, item := range typed {
			folded, err := foldYSONAttributes(item)
			if err != nil {
				return nil, err
			}
			typed[i] = folded
		}
		return typed, nil
	case map[string]any:
		return foldYSONAttributeMap(typed)
	case map[any]any:
		return foldYSONAttributeMap(typed)
	default:
		return value, nil
	}
}

func foldYSONAttributeMap[T comparable](node map[T]any) (any, error) {
	nodeValue, nodeAttrs, isAttributeNode := ysonAttributeFields(node)
	if !isAttributeNode {
		return foldYSONAttributeMapValues(node)
	}

	foldedValue, err := foldYSONAttributes(nodeValue)
	if err != nil {
		return nil, err
	}

	attrs := map[string]any(nil)
	if nodeAttrs != nil {
		attrs, err = stringMap(nodeAttrs)
		if err != nil {
			return nil, fmt.Errorf("$attributes must be an object: %w", err)
		}
		attrs, err = foldYSONAttributeMapValues(attrs)
		if err != nil {
			return nil, err
		}
	}

	return &yson.ValueWithAttrs{Attrs: attrs, Value: foldedValue}, nil
}

func foldYSONAttributeMapValues[T comparable](node map[T]any) (map[string]any, error) {
	converted, err := stringMap(node)
	if err != nil {
		return nil, err
	}

	for key, item := range converted {
		folded, err := foldYSONAttributes(item)
		if err != nil {
			return nil, err
		}
		converted[key] = folded
	}
	return converted, nil
}

func ysonAttributeFields[T comparable](node map[T]any) (any, any, bool) {
	if len(node) > 2 {
		return nil, nil, false
	}

	var value, attrs any
	matches := 0
	for key, item := range node {
		switch any(key) {
		case "$value":
			value = item
			matches++
		case "$attributes":
			attrs = item
			matches++
		}
	}

	return value, attrs, len(node) == matches
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

func readValue(values []string) (any, error) {
	if len(values) > 0 {
		return parseValue(values[0])
	}
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return nil, err
	}
	return parseValue(string(data))
}
