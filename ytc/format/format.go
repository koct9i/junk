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

// Parameter configures a data format encoder or decoder.
type Parameter struct {
	Name  string
	Value any
}

// Format constructs encoders and decoders for a data format.
type Format interface {
	Name() string
	NewEncoder(w io.Writer) Encoder
	NewDecoder(r io.Reader) Decoder
	WithParameters(...Parameter) (Format, error)
}

type jsonFormat struct {
	parameters []Parameter
	format     string
}

func (f *jsonFormat) Name() string { return "json" }

func (f *jsonFormat) NewEncoder(w io.Writer) Encoder {
	encoder := json.NewEncoder(w)
	if f.format == "pretty" {
		encoder.SetIndent("", "  ")
	}
	return encoder
}

func (f *jsonFormat) NewDecoder(r io.Reader) Decoder { return json.NewDecoder(r) }

func (f *jsonFormat) WithParameters(parameters ...Parameter) (Format, error) {
	clone := *f
	if err := clone.applyParameters(parameters...); err != nil {
		return nil, err
	}
	return &clone, nil
}

func (f *jsonFormat) applyParameters(parameters ...Parameter) error {
	f.parameters = append([]Parameter{}, f.parameters...)
	for _, parameter := range parameters {
		format, err := parseFormatParameter(parameter)
		if err != nil {
			return err
		}
		f.format = format
		f.parameters = append(f.parameters, parameter)
	}
	return nil
}

type ysonFormat struct {
	parameters []Parameter
	format     yson.Format
}

func (f *ysonFormat) Name() string { return "yson" }

func (f *ysonFormat) NewEncoder(w io.Writer) Encoder {
	return ysonEncoder{w: w, format: f.format}
}

func (f *ysonFormat) NewDecoder(r io.Reader) Decoder { return yson.NewDecoder(r) }

func (f *ysonFormat) WithParameters(parameters ...Parameter) (Format, error) {
	clone := *f
	if err := clone.applyParameters(parameters...); err != nil {
		return nil, err
	}
	return &clone, nil
}

func (f *ysonFormat) applyParameters(parameters ...Parameter) error {
	f.parameters = append([]Parameter{}, f.parameters...)
	for _, parameter := range parameters {
		format, err := parseFormatParameter(parameter)
		if err != nil {
			return err
		}
		switch format {
		case "binary":
			f.format = yson.FormatBinary
		case "text":
			f.format = yson.FormatText
		case "pretty":
			f.format = yson.FormatPretty
		}
		f.parameters = append(f.parameters, parameter)
	}
	return nil
}

type yamlFormat struct {
	parameters []Parameter
	format     string
}

func (f *yamlFormat) Name() string { return "yaml" }

func (f *yamlFormat) NewEncoder(w io.Writer) Encoder { return yaml.NewEncoder(w) }

func (f *yamlFormat) NewDecoder(r io.Reader) Decoder { return yaml.NewDecoder(r) }

func (f *yamlFormat) WithParameters(parameters ...Parameter) (Format, error) {
	clone := *f
	if err := clone.applyParameters(parameters...); err != nil {
		return nil, err
	}
	return &clone, nil
}

func (f *yamlFormat) applyParameters(parameters ...Parameter) error {
	f.parameters = append([]Parameter{}, f.parameters...)
	for _, parameter := range parameters {
		format, err := parseFormatParameter(parameter)
		if err != nil {
			return err
		}
		f.format = format
		f.parameters = append(f.parameters, parameter)
	}
	return nil
}

func parseFormatParameter(parameter Parameter) (string, error) {
	switch strings.ToLower(strings.TrimSpace(parameter.Name)) {
	case "format":
		value, ok := parameter.Value.(string)
		if !ok {
			return "", fmt.Errorf("parameter %q must be a string", parameter.Name)
		}
		value = strings.ToLower(strings.TrimSpace(value))
		switch value {
		case "binary", "text", "pretty":
			return value, nil
		default:
			return "", fmt.Errorf("unsupported format parameter value %q", parameter.Value)
		}
	default:
		return "", fmt.Errorf("unsupported parameter %q", parameter.Name)
	}
}

var (
	JSON Format = &jsonFormat{format: "pretty"}
	YSON Format = &ysonFormat{format: yson.FormatPretty}
	YAML Format = &yamlFormat{format: "text"}
)

// Parse returns the Format selected by name and applies optional parameters.
func Parse(name string, parameters ...Parameter) (Format, error) {
	formatName, nameParameters, err := parseName(name)
	if err != nil {
		return nil, err
	}
	parameters = append(nameParameters, parameters...)

	var format Format
	switch strings.ToLower(strings.TrimSpace(formatName)) {
	case "json":
		format = JSON
	case "yson":
		format = YSON
	case "yaml", "yml":
		format = YAML
	default:
		return nil, fmt.Errorf("unsupported format %q", name)
	}
	return format.WithParameters(parameters...)
}

type formatName struct {
	Parameters map[string]any `yson:",attrs"`
	Name       string         `yson:",value"`
}

func parseName(name string) (string, []Parameter, error) {
	var parsed formatName
	if err := yson.Unmarshal([]byte(strings.TrimSpace(name)), &parsed); err != nil {
		return "", nil, fmt.Errorf("invalid format name %q: %w", name, err)
	}

	parameters := make([]Parameter, 0, len(parsed.Parameters))
	for parameterName, parameterValue := range parsed.Parameters {
		parameters = append(parameters, Parameter{Name: parameterName, Value: parameterValue})
	}
	return parsed.Name, parameters, nil
}

type ysonEncoder struct {
	w      io.Writer
	format yson.Format
}

func (e ysonEncoder) Encode(v any) error {
	data, err := yson.MarshalFormat(v, e.format)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(e.w, string(data))
	return err
}
