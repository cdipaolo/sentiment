package sentiment

// SplitSentences takes in a rune r and
// returns whether the rune is a sentence
// delimiter ('.', '?', or '!').
//
// It satisfies the interface for
// strings.FieldsFunc()
func SplitSentences(r rune) bool {
	switch r {
	case '.', '?', '!':
		return true
	}
	return false
}
