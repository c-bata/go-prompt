package prompt

type InfoWindow interface {
	GetLines(int) []string
}
