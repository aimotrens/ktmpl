package tmplext

import (
	"os"
	"path/filepath"
)

func includeAsYamlFields(pattern string) (string, error) {
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

		result += fi.Name() + ": |\n"
		result += indent(2, string(b)) + "\n"
	}

	return result, nil
}
