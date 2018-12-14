package debug

import (
	"fmt"
	"os"
)

const (
	envAssertPanic = "GO_PROMPT_ENABLE_ASSERT"
)

var (
	enableAssert bool
)

func init() {
	if e := os.Getenv(envAssertPanic); e == "true" || e == "1" {
		enableAssert = true
	}
}

// Assert ensures expected condition.
func Assert(cond bool, msg interface{}) {
	if cond {
		return
	}
	if enableAssert {
		panic(msg)
	}
	writeWithSync(2, "[ASSERT] "+toString(msg))
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

// AssertNoError ensures err is nil.
func AssertNoError(err error) {
	if err == nil {
		return
	}
	if enableAssert {
		panic(err)
	}
	writeWithSync(2, "[ASSERT] "+err.Error())
}
