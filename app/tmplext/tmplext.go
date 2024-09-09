package tmplext

import "html/template"

func GetTemplateFuncMap(runTemplate func(string) (string, error)) template.FuncMap {
	return template.FuncMap{
		"toYaml":              toYaml,
		"indent":              indent,
		"substr":              substr,
		"iterate":             iterate,
		"format":              format,
		"endsWith":            endsWith,
		"startsWith":          startsWith,
		"contains":            contains,
		"include":             createIncludeFunc(runTemplate, false),
		"includeAsYamlFields": createIncludeFunc(runTemplate, true),
	}
}
