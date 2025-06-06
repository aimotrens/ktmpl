package ktmpl

import (
	"fmt"
	"runtime"
)

func Version(ktmplVersion, compileDate string) {
	fmt.Printf("ktmpl %s, compiled at %s with %v\n", ktmplVersion, compileDate, runtime.Version())
}
