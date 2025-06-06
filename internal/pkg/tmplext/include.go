package tmplext

import (
	"os"
	"path/filepath"
)

func createIncludeFunc(runTemplate func(t string) (string, error), asYamlFields bool) func(pattern string) (string, error) {
	return func(pattern string) (string, error) {
		files, err := filepath.Glob(pattern)
		if err != nil {
			return "", err
		}

		var result string
		for _, f := range files {
			fi, err := os.Stat(f)
			if err != nil {
				return "", err
			}

			if fi.IsDir() {
				continue
			}

			b, err := os.ReadFile(f)
			if err != nil {
				return "", err
			}

			res, err := runTemplate(string(b))
			if err != nil {
				return "", err
			}

			if asYamlFields {
				result += fi.Name() + ": |\n"
				result += indent(2, res) + "\n"
			} else {
				result += res
			}
		}

		return result, nil
	}
}
