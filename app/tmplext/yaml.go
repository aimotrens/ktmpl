package tmplext

import (
	"bytes"

	"gopkg.in/yaml.v3"
)

func toYaml(obj any) string {
	dst := bytes.NewBuffer([]byte{})
	encoder := yaml.NewEncoder(dst)
	encoder.SetIndent(2)
	encoder.Encode(obj)
	encoder.Close()

	return dst.String()
}
