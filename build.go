package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/aimotrens/ktmpl/app/templating"
)

func main() {
	buildBinaryFlag := flag.Bool("binary", false, "build binary")
	installFlag := flag.Bool("install", false, "install binary")

	buildDockerFlag := flag.Bool("docker", false, "build docker image")
	imageNameFlag := flag.String("image", "ktmpl", "docker image name")
	withoutKubectlFlag := flag.Bool("without-kubectl", false, "exclude kubectl in docker image")

	flag.Parse()

	if *buildBinaryFlag {
		buildBinary()

		if *installFlag {
			installBinary()
		}
	}

	if *buildDockerFlag {
		if *imageNameFlag == "" {
			panic("image name is required")
		}

		buildDocker(*imageNameFlag, *withoutKubectlFlag)
	}

}

func buildBinary() {
	createCmd := func(osName, output string) *exec.Cmd {
		cmd := exec.Command("go", "build", "-o", output, "./app")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = append(os.Environ(), "GOOS="+osName, "GOARCH=amd64")
		return cmd
	}

	createCmd("linux", "bin/ktmpl_linux_amd64").Run()
	createCmd("windows", "bin/ktmpl_windows_amd64.exe").Run()
}

func installBinary() {
	var cmd *exec.Cmd
	if runtime.GOOS == "linux" {
		cmd = exec.Command("cp", "./bin/ktmpl_linux_amd64", os.ExpandEnv("$HOME/.local/bin/ktmpl"))
	} else if runtime.GOOS == "windows" {
		cmd = exec.Command("cp", "./bin/ktmpl_windows_amd64.exe", os.ExpandEnv("%LOCALAPPDATA%/Programs/ktmpl.exe"))
	} else {
		panic("install is not supported on " + runtime.GOOS + " yet")
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}

func buildDocker(imageNameFlag string, withoutKubectlFlag bool) {
	dockerfile := `
FROM golang:1.20-alpine AS builder
WORKDIR /src
ADD . .
RUN go build -o ./bin/ktmpl ./app

# ---
{{if .kubectl }}
FROM debian:latest as kubectl_downloader

RUN apt-get update && apt-get install -y ca-certificates curl && rm -rf /var/lib/apt/lists/*
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
RUN chmod +x ./kubectl

# ---
{{- end}}

FROM debian:latest

COPY --from=builder /src/bin/ktmpl /usr/bin/ktmpl

{{- if .kubectl}}
COPY --from=kubectl_downloader /kubectl /usr/bin/kubectl
{{- end}}	
`

	renderedDockerFile, err := templating.RunTemplateInMemory(
		dockerfile,
		map[string]any{"kubectl": !withoutKubectlFlag},
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(renderedDockerFile)

	cmd := exec.Command("docker", "build", "-t", imageNameFlag, "-f", "-", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = strings.NewReader(renderedDockerFile)
	cmd.Run()
}
