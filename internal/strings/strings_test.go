package strings_test

import (
	"fmt"

	"github.com/confluentinc/go-prompt/internal/strings"
)

func ExampleIndexNotByte() {
	fmt.Println(strings.IndexNotByte("golang", 'g'))
	fmt.Println(strings.IndexNotByte("golang", 'x'))
	fmt.Println(strings.IndexNotByte("gggggg", 'g'))
	// Output:
	// 1
	// 0
	// -1
}

func ExampleLastIndexNotByte() {
	fmt.Println(strings.LastIndexNotByte("golang", 'g'))
	fmt.Println(strings.LastIndexNotByte("golang", 'x'))
	fmt.Println(strings.LastIndexNotByte("gggggg", 'g'))
	// Output:
	// 4
	// 5
	// -1
}

func ExampleIndexNotAny() {
	fmt.Println(strings.IndexNotAny("golang", "glo"))
	fmt.Println(strings.IndexNotAny("golang", "gl"))
	fmt.Println(strings.IndexNotAny("golang", "golang"))
	// Output:
	// 3
	// 1
	// -1
}

func ExampleLastIndexNotAny() {
	fmt.Println(strings.LastIndexNotAny("golang", "agn"))
	fmt.Println(strings.LastIndexNotAny("golang", "an"))
	fmt.Println(strings.LastIndexNotAny("golang", "golang"))
	// Output:
	// 2
	// 5
	// -1
}
