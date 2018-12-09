package debug

import (
	"fmt"
	"os"
)

const (
	envAssertPanic = "GO_PROMPT_ENABLE_ASSERT"
)

var (
	assertFunc = assertLog
)

func init() {
	enableAssert := os.Getenv(envAssertPanic)
	if enableAssert == "true" || enableAssert == "1" {
		assertFunc = assertPanic
	}
}

func assertPanic(msg interface{}) {
	panic(msg)
}

func assertLog(msg interface{}) {
	calldepth := 3
	writeWithSync(calldepth, "[ASSERT] "+toString(msg))
}

// Assert raise panic or write log if cond is false.
func Assert(cond bool, msg interface{}) {
	if cond {
		return
	}
	assertFunc(msg)
}

func toString(v interface{}) string {
	switch a := v.(type) {
	case func() string:
		return a()
	case string:
		return a
	case fmt.Stringer:
		return a.String()
	default:
		return fmt.Sprintf("unexpected type, %t", v)
	}
}
