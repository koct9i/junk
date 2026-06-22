package format

import (
	"bytes"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	for _, name := range []string{"json", "yson", "yaml", "yml", "<format=text>yson"} {
		if _, err := Parse(name); err != nil {
			t.Fatalf("Parse(%q) returned error: %v", name, err)
		}
	}

	if _, err := Parse("xml"); err == nil {
		t.Fatal("Parse(\"xml\") succeeded, want error")
	}
}

func TestParseParameters(t *testing.T) {
	format, err := Parse("<format=text>yson", Parameter{Name: "format", Value: "pretty"})
	if err != nil {
		t.Fatalf("Parse() returned error: %v", err)
	}

	var buf bytes.Buffer
	if err := format.NewEncoder(&buf).Encode(map[string]any{"key": "value"}); err != nil {
		t.Fatalf("Encode() returned error: %v", err)
	}
	if got := buf.String(); !strings.Contains(got, "\n") {
		t.Fatalf("encoded YSON = %q, want pretty output", got)
	}
}

func TestWithParametersYSONText(t *testing.T) {
	format, err := YSON.WithParameters(Parameter{Name: "format", Value: "text"})
	if err != nil {
		t.Fatalf("WithParameters() returned error: %v", err)
	}

	var buf bytes.Buffer
	if err := format.NewEncoder(&buf).Encode(map[string]any{"key": "value"}); err != nil {
		t.Fatalf("Encode() returned error: %v", err)
	}
	if got := buf.String(); strings.Contains(got, "\n    ") {
		t.Fatalf("encoded YSON = %q, want text output", got)
	}
}

func TestWithParametersRejectsInvalidParameter(t *testing.T) {
	if _, err := YSON.WithParameters(Parameter{Name: "unknown", Value: "text"}); err == nil {
		t.Fatal("WithParameters() succeeded with unknown parameter, want error")
	}
	if _, err := YSON.WithParameters(Parameter{Name: "format", Value: "compact"}); err == nil {
		t.Fatal("WithParameters() succeeded with invalid format, want error")
	}
}

func TestFormatsEncodeDecode(t *testing.T) {
	formats := []Format{JSON, YSON, YAML}
	for _, format := range formats {
		t.Run(format.Name(), func(t *testing.T) {
			var buf bytes.Buffer
			input := map[string]any{"key": "value", "number": int64(42)}
			if err := format.NewEncoder(&buf).Encode(input); err != nil {
				t.Fatalf("Encode() returned error: %v", err)
			}

			var output map[string]any
			if err := format.NewDecoder(&buf).Decode(&output); err != nil {
				t.Fatalf("Decode() returned error: %v", err)
			}
			if output["key"] != "value" {
				t.Fatalf("decoded key = %v, want value", output["key"])
			}
		})
	}
}
