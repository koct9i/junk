// Package format provides a common interface for YT data formats.
package format

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"go.ytsaurus.tech/yt/go/yson"
	"gopkg.in/yaml.v3"
)

// Encoder writes a single value in a data format.
type Encoder interface {
	Encode(v any) error
}

// Decoder reads a single value in a data format.
type Decoder interface {
	Decode(v any) error
}

// Format constructs encoders and decoders for a data format.
type Format interface {
	Name() string
	NewEncoder(w io.Writer) Encoder
	NewDecoder(r io.Reader) Decoder
}

type factory struct {
	name       string
	newEncoder func(io.Writer) Encoder
	newDecoder func(io.Reader) Decoder
}

func (f factory) Name() string { return f.name }

func (f factory) NewEncoder(w io.Writer) Encoder { return f.newEncoder(w) }

func (f factory) NewDecoder(r io.Reader) Decoder { return f.newDecoder(r) }

var (
	JSON Format = factory{
		name: "json",
		newEncoder: func(w io.Writer) Encoder {
			encoder := json.NewEncoder(w)
			encoder.SetIndent("", "  ")
			return encoder
		},
		newDecoder: func(r io.Reader) Decoder { return json.NewDecoder(r) },
	}

	YSON Format = factory{
		name:       "yson",
		newEncoder: func(w io.Writer) Encoder { return ysonEncoder{w: w} },
		newDecoder: func(r io.Reader) Decoder { return yson.NewDecoder(r) },
	}

	YAML Format = factory{
		name:       "yaml",
		newEncoder: func(w io.Writer) Encoder { return yaml.NewEncoder(w) },
		newDecoder: func(r io.Reader) Decoder { return yaml.NewDecoder(r) },
	}
)

// Parse returns the Format selected by name.
func Parse(name string) (Format, error) {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "json":
		return JSON, nil
	case "yson":
		return YSON, nil
	case "yaml", "yml":
		return YAML, nil
	default:
		return nil, fmt.Errorf("unsupported format %q", name)
	}
}

type ysonEncoder struct {
	w io.Writer
}

func (e ysonEncoder) Encode(v any) error {
	data, err := yson.MarshalFormat(v, yson.FormatPretty)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(e.w, string(data))
	return err
}
