package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"runtime"

	"os"
	"path/filepath"
	"strings"

	"github.com/aimotrens/ktmpl/app/templating"
	"gopkg.in/yaml.v3"
)

var compileDate, ktmplVersion string

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
		fmt.Printf("ktmpl %s, compiled at %s with %v on %v/%v\n", ktmplVersion, compileDate, runtime.Version(), runtime.GOOS, runtime.GOARCH)
		return
	}

	if flag.NArg() == 0 {
		flag.Usage()
		return
	}

	var opendValuesFile *os.File
	var opendOutputFile *os.File
	var err error

	inputTemplate := flag.Arg(0)

	if inputTemplate == "" {
		fmt.Println("No input specified")
		return
	} else if _, err := os.Stat(inputTemplate); os.IsNotExist(err) {
		fmt.Println("Input does not exist")
		return
	}

	if valuesFile != "" {
		if valuesFile == "-" {
			opendValuesFile = os.Stdin
		} else {
			opendValuesFile, err = os.Open(valuesFile)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		defer opendValuesFile.Close()
	}

	if output == "-" {
		opendOutputFile = os.Stdout
	} else {
		opendOutputFile, err = os.Create(output)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	defer opendOutputFile.Close()

	values := make(map[string]any)
	if opendValuesFile != nil {
		err = readValuesFile(values, opendValuesFile)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if addEnv {
		addEnvToValues(values)
	}

	inputFiles, err := getInputFiles(inputTemplate, recursive, []string{valuesFile, output})
	if err != nil {
		fmt.Println(err)
		return
	}

	for i, inputFile := range inputFiles {
		templating.RunTemplateFile(inputFile, values, opendOutputFile)

		if i < len(inputFiles)-1 {
			opendOutputFile.WriteString("\n---\n")
		}
	}
}

func getInputFiles(input string, recursive bool, exclude []string) ([]string, error) {
	isValidInputFile := func(fileinfo fs.FileInfo) bool {
		for _, exclude := range exclude {
			absExclude, _ := filepath.Abs(exclude)
			absFileinfo, _ := filepath.Abs(fileinfo.Name())
			if absExclude == absFileinfo {
				return false
			}
		}

		return fileinfo.Mode().IsRegular() &&
			(strings.HasSuffix(fileinfo.Name(), ".yml") ||
				strings.HasSuffix(fileinfo.Name(), ".yaml"))
	}

	inputStat, err := os.Stat(input)
	if err != nil {
		return nil, err
	}

	if inputStat.IsDir() {
		if recursive {
			inputFiles := []string{}
			err := filepath.Walk(input, func(path string, info fs.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if isValidInputFile(info) {
					inputFiles = append(inputFiles, path)
				}

				return nil
			})

			return inputFiles, err
		} else {
			inputFiles, err := os.ReadDir(input)
			if err != nil {
				return nil, err
			}

			var files []string
			for _, file := range inputFiles {
				if info, err := file.Info(); err == nil && isValidInputFile(info) {
					files = append(files, filepath.Join(input, file.Name()))
				} else if err != nil {
					return nil, err
				}
			}

			return files, nil
		}
	} else if inputStat.Mode().IsRegular() {
		return []string{input}, nil
	} else {
		panic("Input is not a directory or a file")
	}
}

func readValuesFile(values map[string]any, valuesFile *os.File) error {
	valuesBytes, err := io.ReadAll(valuesFile)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(valuesBytes, &values)
}

func addEnvToValues(values map[string]any) {
	envMap := make(map[string]any)
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		envMap[pair[0]] = pair[1]
	}
	values["env"] = envMap
}
