package prompt

// InfoWindow is the interface that needs
// to be implemented by new types of information
// windows
type InfoWindow interface {
	GetLines(int) []string
}
