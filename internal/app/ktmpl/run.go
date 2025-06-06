package ktmpl

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/aimotrens/ktmpl/internal/pkg/templating"
	"gopkg.in/yaml.v3"
)

func Run(valuesFile, inputTemplate, output string, addEnv bool, recursive bool) {
	var opendValuesFile *os.File
	var opendOutputFile *os.File
	var err error

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

				if info.IsDir() {
					if c, err := containsIgnoreFile(path); err != nil {
						return err
					} else if c {
						return filepath.SkipDir
					}
				}

				if isValidInputFile(info) {
					inputFiles = append(inputFiles, path)
				}

				return nil
			})

			return inputFiles, err
		} else {
			if c, err := containsIgnoreFile(input); err != nil {
				return nil, err
			} else if c {
				return []string{}, nil
			}

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

func containsIgnoreFile(dir string) (bool, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return false, err
	}

	return slices.ContainsFunc(files, func(file fs.DirEntry) bool {
		return file.Name() == ".ktmpl_ignore_dir"
	}), nil
}
