package prompt

import (
	"bytes"
	"fmt"
	"os"

	"github.com/noborus/ov/oviewer"
	"golang.org/x/term"
)

func pager_print(b []byte, config *oviewer.Config) {
	// set console redirect (https://github.com/noborus/ov/blob/879f450a792e56646d038d726ca0aa1559e4ae3b/main.go)
	if term.IsTerminal(int(os.Stdout.Fd())) {
		tmpStdout := os.Stdout
		os.Stdout = nil
		defer func() {
			os.Stdout = tmpStdout
		}()
	} else {
		oviewer.STDOUTPIPE = os.Stdout
	}
	if term.IsTerminal(int(os.Stderr.Fd())) {
		tmpStderr := os.Stderr
		os.Stderr = nil
		defer func() {
			os.Stderr = tmpStderr
		}()
	} else {
		oviewer.STDERRPIPE = os.Stderr
	}

	// print (https://github.com/noborus/ov/blob/879f450a792e56646d038d726ca0aa1559e4ae3b/main.go#L147)
	ov, err := oviewer.NewRoot(bytes.NewReader(b))
	if config != nil {
		ov.SetConfig(*config)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		return
	}
	if err := ov.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		return
	}
	ov.WriteOriginal()
}
