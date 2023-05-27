package tmplext

import "html/template"

func GetTemplateFuncMap() template.FuncMap {
	return template.FuncMap{
		"toYaml": toYaml,
		"indent": indent,
	}
}
