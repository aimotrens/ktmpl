package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/aimotrens/ktmpl/app/templating"
)

const (
	binaryLinux   = "./bin/ktmpl_linux_amd64"
	binaryWindows = "./bin/ktmpl_windows_amd64.exe"
)

func main() {
	buildBinaryFlag := flag.Bool("binary", false, "build binary")
	installFlag := flag.Bool("install", false, "install binary")
	installPathFlag := flag.String("install-path", "", "install path")

	buildDockerFlag := flag.Bool("docker", false, "build docker image")
	imageNameFlag := flag.String("image", "ktmpl", "docker image name")
	imagePushFlag := flag.Bool("push", false, "push docker image")
	withoutKubectlFlag := flag.Bool("without-kubectl", false, "exclude kubectl in docker image")

	flag.Parse()

	if *buildBinaryFlag || *installFlag {
		buildBinary()

		if *installFlag {
			installBinary(*installPathFlag)
		}
	}

	if *buildDockerFlag {
		if *imageNameFlag == "" {
			panic("image name is required")
		}

		buildDocker(*imageNameFlag, *withoutKubectlFlag, *imagePushFlag)
	}
}

func buildBinary() {
	createCmd := func(osName, output string) *exec.Cmd {
		cmd := exec.Command("go", "build",
			"-ldflags", "-X \"main.compileDate="+time.Now().Format("02.01.2006 15:04:05")+"\"",
			"-o", output,
			"./app",
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Env = append(os.Environ(), "GOOS="+osName, "GOARCH=amd64")
		return cmd
	}

	createCmd("linux", binaryLinux).Run()
	createCmd("windows", binaryWindows).Run()
}

func installBinary(installPath string) {
	var src, dst string
	switch runtime.GOOS {
	case "linux":
		src = binaryLinux
		dst = os.ExpandEnv("$HOME/.local/bin/ktmpl")
	case "windows":
		src = binaryWindows
		dst = os.ExpandEnv("$LOCALAPPDATA/Programs/kube-tools/ktmpl.exe")
	default:
		panic("install is not supported on " + runtime.GOOS + " yet")
	}

	if installPath != "" {
		dst = installPath
	}

	_, err := CopyFile(src, dst)
	if err != nil {
		panic(err)
	}
}

func buildDocker(imageNameFlag string, withoutKubectlFlag bool, imagePushFlag bool) {
	dockerfile := `
FROM golang:{{.goversion}}-alpine AS builder
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
		map[string]any{
			"goversion": runtime.Version()[2:],
			"kubectl":   !withoutKubectlFlag,
		},
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

	if imagePushFlag {
		cmd = exec.Command("docker", "push", imageNameFlag)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run()
	}
}

func CopyFile(src, dst string) (written int64, err error) {
	srcFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer dstFile.Close()

	return io.Copy(dstFile, srcFile)
}
