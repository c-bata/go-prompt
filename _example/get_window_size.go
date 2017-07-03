package main

import (
	"fmt"
	"syscall"

	"github.com/c-bata/go-prompt-toolkit/prompt"
)

func main() {
	w := prompt.GetWinSize(syscall.Stdin)
	fmt.Printf("Row %d, Col %d\n", w.Row, w.Col)
}
