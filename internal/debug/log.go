package debug

import (
	"io/ioutil"
	"log"
	"os"
)

const (
	envEnableLog = "GO_PROMPT_ENABLE_LOG"
	logFileName  = "go-prompt.log"
)

var (
	logfile *os.File
	logger  *log.Logger
)

func init() {
	if e := os.Getenv(envEnableLog); e == "true" || e == "1" {
		var err error
		logfile, err = os.OpenFile(logFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err == nil {
			logger = log.New(logfile, "", log.Llongfile)
			return
		}
	}
	logger = log.New(ioutil.Discard, "", log.Llongfile)
}

// Teardown to close logfile
func Teardown() {
	if logfile == nil {
		return
	}
	_ = logfile.Close()
}

func writeWithSync(calldepth int, msg string) {
	calldepth++
	if logfile == nil {
		return
	}
	_ = logger.Output(calldepth, msg)
	_ = logfile.Sync() // immediately write msg
}

// Log to output message
func Log(msg string) {
	calldepth := 2
	writeWithSync(calldepth, msg)
}
