package main

import (
	"flag"
	"fmt"

	"github.com/aimotrens/ktmpl/internal/app/ktmpl"
)

var (
	compileDate  = "unknown"
	ktmplVersion = "vX.X.X"
)

var usage = `Usage:
ktmpl [options] <templates>
    <templates> can be a file or a directory.
    If a directory is specified, all .yml/.yaml files in that directory will be processed.

Examples:
    ktmpl -i values.yml template.yml          # process template.yml with values.yml
    ktmpl -e -i values.yml template.yml       # process template.yml with values.yml and environment variables
    ktmpl -e -o output.yml template.yml       # process template.yml with environment variables and write to output.yml

Options:
    -v, --version       print version information
    -r, --recursive     recurse into subdirectories
    -o, --output        output file
    -i, --values        values file (as YAML)
    -e, --env           add environment variables to values

Template functions (last argument is pipeline variable):
    indent(spaces int, s string) string             # indent string with spaces
    substr(start, length int, s string) string      # get substring
    iterate(from, to int) []int                     # create slice of integers
    format(format string, obj any) string           # format string
    toYaml(obj any) string                          # convert object to YAML
    endsWith(suffix, s string) bool                 # check if string ends with suffix
    startsWith(prefix, s string) bool               # check if string starts with prefix
    contains(substring, s string) bool              # check if string contains substring
    include(globPattern string) string              # include files concatinated
    includeAsYamlFields(globPattern string) string  # include files as YAML fields (especially useful for K8s config maps)
`

func main() {
	flag.Usage = func() { fmt.Print(usage) }

	var version, recursive, addEnv bool
	var output, valuesFile string

	flag.BoolVar(&version, "v", false, "print version information")
	flag.BoolVar(&version, "version", false, "print version information")

	flag.BoolVar(&recursive, "r", false, "recurse into subdirectories")
	flag.BoolVar(&recursive, "recursive", false, "recurse into subdirectories")

	flag.StringVar(&output, "o", "-", "output file")
	flag.StringVar(&output, "output", "-", "output file")

	flag.StringVar(&valuesFile, "i", "", "values file (as YAML)")
	flag.StringVar(&valuesFile, "values", "", "values file (as YAML)")

	flag.BoolVar(&addEnv, "e", false, "add environment variables to values")
	flag.BoolVar(&addEnv, "env", false, "add environment variables to values")

	flag.Parse()

	if version {
		ktmpl.Version(ktmplVersion, compileDate)
		return
	}

	if flag.NArg() == 0 {
		ktmpl.Usage()
		return
	}

	ktmpl.Run(valuesFile, output, addEnv, recursive)
}
