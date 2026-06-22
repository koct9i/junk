package cypress

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/go-logr/logr"

	"go.ytsaurus.tech/yt/go/yson"
	ytsdk "go.ytsaurus.tech/yt/go/yt"
	"go.ytsaurus.tech/yt/go/yt/ythttp"

	"github.com/koct9i/junk/ytc/config"
	dataformat "github.com/koct9i/junk/ytc/format"
	"github.com/koct9i/junk/ytc/log"
)

func client(ctx context.Context) (ytsdk.Client, error) {
	config := config.Config
	config.Logger = log.NewYTSlog(logr.FromContextAsSlogLogger(ctx))
	return ythttp.NewClient(&config)
}

func printValue(w io.Writer, v any) error {
	format, err := dataformat.Parse(config.Format)
	if err != nil {
		return err
	}
	return format.NewEncoder(w).Encode(v)
}

func parseValue(s string) (any, error) {
	format, err := dataformat.Parse(config.Format)
	if err != nil {
		return nil, err
	}

	var value any
	s = strings.TrimSpace(s)
	if err := format.NewDecoder(strings.NewReader(s)).Decode(&value); err != nil {
		if format == dataformat.JSON {
			return s, nil
		}
		return nil, err
	}
	if format == dataformat.YSON {
		return value, nil
	}
	return foldYSONAttributes(value)
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
