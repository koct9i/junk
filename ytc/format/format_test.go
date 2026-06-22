package format

import (
	"bytes"
	"testing"
)

func TestParse(t *testing.T) {
	for _, name := range []string{"json", "yson", "yaml", "yml"} {
		if _, err := Parse(name); err != nil {
			t.Fatalf("Parse(%q) returned error: %v", name, err)
		}
	}

	if _, err := Parse("xml"); err == nil {
		t.Fatal("Parse(\"xml\") succeeded, want error")
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
