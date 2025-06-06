package ktmpl_test

import (
	"os"
	"testing"

	"github.com/aimotrens/ktmpl/internal/app/ktmpl"
)

func TestRunDefault(t *testing.T) {
	valuesFile := "_values.yml"
	templateFile := "_template.yml"
	outputFile := "_output.yml"

	testCases := []struct {
		name     string
		values   string
		template string
		output   string
	}{
		{
			name:     "Simple valid template",
			values:   "foo: 1",
			template: "exampleFoo: {{ .foo }}",
			output:   "exampleFoo: 1",
		},
		{
			name:     "Template with multiple values",
			values:   "foo: 1\nbar: 2",
			template: "exampleFoo: {{ .foo }}\nexampleBar: {{ .bar }}",
			output:   "exampleFoo: 1\nexampleBar: 2",
		},
		{
			name:     "Template with nested values",
			values:   "foo:\n  bar: 2",
			template: "exampleFooBar: {{ .foo.bar }}",
			output:   "exampleFooBar: 2",
		},
		{
			name:     "Template with missing value",
			values:   "foo: 1",
			template: "exampleFoo: {{ .foo }}\nexampleBar: {{ .bar }}",
			output:   "exampleFoo: 1\nexampleBar: <no value>",
		},
		{
			name:     "Template with substr function",
			values:   "foo: Hello",
			template: "exampleFoo: {{ substr 0 3 .foo }}",
			output:   "exampleFoo: Hel",
		},
		{
			name:     "Template with iterate function",
			values:   "",
			template: "exampleIterate: {{ range $i := iterate 1 5 }}{{$i}} {{end}}",
			output:   "exampleIterate: 1 2 3 4 5 ",
		},
		{
			name:     "Template with includeAsYamlFields function",
			values:   "foo: 1",
			template: "exampleInclude:\n  {{ includeAsYamlFields \"_values.yml\" }}",
			output:   "exampleInclude:\n  _values.yml: |\n  foo: 1\n",
		},
		{
			name:     "Template with include function",
			values:   "foo: 1",
			template: "exampleInclude: {{ include \"_values.yml\" }}",
			output:   "exampleInclude: foo: 1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create temporary files for values and template
			if err := os.WriteFile(valuesFile, []byte(tc.values), 0644); err != nil {
				t.Fatalf("Failed to write values file: %v", err)
			}
			if err := os.WriteFile(templateFile, []byte(tc.template), 0644); err != nil {
				t.Fatalf("Failed to write template file: %v", err)
			}
			defer os.Remove(valuesFile)
			defer os.Remove(templateFile)

			// Run the ktmpl command
			ktmpl.Run(valuesFile, templateFile, outputFile, false, false)
			defer os.Remove(outputFile)

			// Read the output file
			output, err := os.ReadFile(outputFile)
			if err != nil {
				t.Fatalf("Failed to read output file: %v", err)
			}
			if string(output) != tc.output {
				t.Errorf("Expected output %q, got %q", tc.output, string(output))
			}
		})
	}
}
