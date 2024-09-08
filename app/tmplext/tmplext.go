package tmplext

import "html/template"

func GetTemplateFuncMap() template.FuncMap {
	return template.FuncMap{
		"toYaml":              toYaml,
		"indent":              indent,
		"substr":              substr,
		"iterate":             iterate,
		"format":              format,
		"endsWith":            endsWith,
		"startsWith":          startsWith,
		"contains":            contains,
		"includeAsYamlFields": includeAsYamlFields,
	}
}
