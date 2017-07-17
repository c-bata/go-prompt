package prompt

type Color int

const (
	DefaultColor Color = iota

	// Low intensity
	Black
	DarkRed
	DarkGreen
	Brown
	DarkBlue
	Purple
	Cyan
	LightGray

	// High intensity
	DarkGray
	Red
	Green
	Yellow
	Blue
	Fuchsia
	Turquoise
	White
)
