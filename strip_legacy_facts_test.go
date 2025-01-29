package main

import (
	"bytes"
	"encoding/json"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestStripLegacyFac(t *testing.T) {
	in, err := os.Open("testdata/input-facts.json")
	if err != nil {
		t.Fatalf("read input data: %v", err)
	}
	defer in.Close()

	want, err := os.ReadFile("testdata/output-facts.json")
	if err != nil {
		t.Fatalf("read wanted output data: %v", err)
	}

	var got bytes.Buffer
	if err := stripLegacyFacts(in, &got); err != nil {
		t.Fatal(err)
	}

	transformJSON := cmp.FilterValues(func(x, y []byte) bool {
		return json.Valid(x) && json.Valid(y)
	}, cmp.Transformer("ParseJSON", func(in []byte) (out interface{}) {
		if err := json.Unmarshal(in, &out); err != nil {
			panic(err) // should never occur given previous filter to ensure valid JSON
		}
		return out
	}))

	if diff := cmp.Diff(want, got.Bytes(), transformJSON); diff != "" {
		t.Errorf("stripLegacyFacts mismatch (-want +got):\n%s", diff)
	}
}
