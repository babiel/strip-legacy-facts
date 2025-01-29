// Command strip-legacy-facts removes legacy facts from a Puppet fact set.
//
// It reads JSON on stdin and prints JSON to stdout.
package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

//go:embed facter-schema.yaml
var facterSchemaYAML []byte

type fact struct {
	Hidden  bool   `yaml:"hidden"`
	Pattern string `yaml:"pattern"`
}

func legacyFactPatterns() []*regexp.Regexp {
	var schema map[string]fact

	if err := yaml.Unmarshal(facterSchemaYAML, &schema); err != nil {
		panic(err)
	}

	var patterns []*regexp.Regexp

	for n, f := range schema {
		if !f.Hidden {
			continue
		}

		p := f.Pattern
		if p == "" {
			p = `^` + regexp.QuoteMeta(n) + `$`
		}

		patterns = append(patterns, regexp.MustCompile(p))
	}

	return patterns

}

func main() {

	in := os.Stdin
	out := os.Stdout

	if err := stripLegacyFacts(in, out); err != nil {
		fmt.Fprintf(os.Stderr, "strip-legacy-facts error: %v", err)
		os.Exit(2)
	}
}

func stripLegacyFacts(in io.Reader, out io.Writer) error {
	patterns := legacyFactPatterns()

	var facts map[string]any

	if err := json.NewDecoder(in).Decode(&facts); err != nil {
		return err
	}

	for k := range facts {
		for _, re := range patterns {
			if re.MatchString(k) {
				delete(facts, k)
				break
			}
		}
	}

	enc := json.NewEncoder(out)
	enc.SetIndent("", "  ")
	return enc.Encode(facts)
}
