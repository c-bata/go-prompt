package main

import (
	"fmt"

	"github.com/c-bata/go-prompt-toolkit/prompt"
)

func main() {
	w := prompt.GetWinSize()
	fmt.Printf("Row %d, Col %d\n", w.Row, w.Col)
}
