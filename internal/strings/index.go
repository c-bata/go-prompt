package strings

// Numeric type that represents an index
// of a single byte in a string, array or slice.
type ByteIndex int

// Numeric type that represents an index
// of a single rune in a string, array or slice.
type RuneIndex int

// Numeric type that represents the visible
// width of characters in a string as seen in a terminal emulator.
type StringWidth int

// Numeric type that represents the amount
// of bytes in a string, array or slice.
type ByteCount = ByteIndex

// Numeric type that represents the amount
// of runes in a string, array or slice.
type RuneCount = RuneIndex
