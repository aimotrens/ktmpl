package templating

import (
	"os"
	"strings"
	"text/template"

	"github.com/aimotrens/ktmpl/app/tmplext"
)

func RunTemplateFile(inputFile string, values map[string]any, outputFile *os.File) {
	templateBytes, err := os.ReadFile(inputFile)
	if err != nil {
		panic(err)
	}
	templateContent := string(templateBytes)

	output, err := RunTemplateInMemory(templateContent, values)
	if err != nil {
		panic(err)
	}
	outputFile.WriteString(output)
}

func RunTemplateInMemory(inputTemplate string, values map[string]any) (string, error) {
	t := template.Must(template.New("template").Funcs(tmplext.GetTemplateFuncMap()).Parse(inputTemplate))

	buf := new(strings.Builder)

	err := t.Execute(buf, values)
	return buf.String(), err
}
