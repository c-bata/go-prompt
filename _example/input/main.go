package main

import (
	"fmt"

	"github.com/c-bata/go-prompt"
)

func main() {
	fruits := []string{"Apple", "Banana", "Bilberry", "Coconuts"}

	fmt.Println("What fruits do you like?")
	f := prompt.Choose("> ", fruits)
	fmt.Println("You like " + f)
}
