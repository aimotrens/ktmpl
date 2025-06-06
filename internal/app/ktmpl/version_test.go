package ktmpl_test

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/aimotrens/ktmpl/internal/app/ktmpl"
)

func TestVersion(t *testing.T) {
	expectedVersion := "v1.0.0"
	expectedCompileDate := "2023-10-01"

	output := captureStdout(func() {
		ktmpl.Version(expectedVersion, expectedCompileDate)
	})

	if !strings.Contains(output, expectedVersion) {
		t.Errorf("Expected version %s in output, got: %s", expectedVersion, output)
	}

	if !strings.Contains(output, expectedCompileDate) {
		t.Errorf("Expected compile date %s in output, got: %s", expectedCompileDate, output)
	}
}

func captureStdout(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}
