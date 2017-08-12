# Internals of `go-prompt`

This is a short description of go-prompt implementation.
go-prompt consists of three parts.

1. Input parser
2. Emulate user input with Buffer object.
3. Render buffer object.

### Input Parser

![input-parser animation](./_resources/input-parser.gif)

Input Parser only supports only vt100 compatible console now.

* Set raw mode.
* Read standard input.
* Parse byte array

### Emulate user input

go-prompt contains Buffer class.
It represents input state by handling user input key.

`Buffer` object has text and cursor position.

**TODO prepare the sample of buffer**

```go
package main

import "github.com/c-bata/go-prompt"

func main() {
  b := prompt.NewBuffer()
  ... wip
}
```

### Renderer

`Renderer` object renders a buffer object.

**TODO prepare the sample of brender**

```go
package main
```

the output is below:

**TODO prepare a screen shot**
