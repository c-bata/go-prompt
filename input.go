package prompt

// ConsoleParser is an interface to abstract input layer.
type ConsoleParser interface {
	// SetUp should be called before starting input
	SetUp() error
	// TearDown should be called after stopping input
	TearDown() error
	// Read returns byte array.
	Read() ([]byte, error)
}

// InputProcessor is worker goroutine to read user input.
type InputProcessor struct {
	UserInput chan []byte
	Pause     chan bool
	in        ConsoleParser
	pause     bool
}

// NewInputProcessor returns object of InputProcessor.
func NewInputProcessor(in ConsoleParser) *InputProcessor {
	return &InputProcessor{
		in:        in,
		UserInput: make(chan []byte, 128),
		Pause:     make(chan bool),
	}
}
