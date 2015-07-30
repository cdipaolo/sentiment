package sentiment

import (
	"strings"

	"golang.org/x/text/transform"
)

type Score struct {
	Word  string  `json:"word"`
	Score float64 `json:"score"`
}

type Analysis struct {
	Words []Score `json:"words"`
	Score float64 `json:"score"`
}

// SentimentOfWord takes in a single word and
// returns the sentiment of that word (negative
// would mean negative sentiment, positive would
// mean positive sentiment)
func SentimentOfWord(word string) float64 {
	if words[word].Count < 6 || len(word) < 3 {
		return 0.0
	}

	return words[word].Probability
}

// SentimentOfSentence takes in a string of a
// space-delimited sentence and returns a
// weighted sum of their probabilities
func SentimentOfSentence(sentence string) float64 {
	var sum float64

	w := strings.Split(sentence, " ")
	for i, word := range w {
		s := SentimentOfWord(word)

		if i > 0 && s > 0 && (w[i-1] == "not" || w[i-1] == "no" || w[i-1] == "never" || w[i-1] == "dont") {
			s *= -1
		}

		if i > 0 && (w[i-1] == "really") {
			s *= 1.4
		}

		sum += s
	}

	return sum
}

// SentimentAnalysis takes in a (possibly 'dirty')
// sentence (or any block of text,) cleans the
// text, finds the sentiment of each word in the
// text, finds the sentiment of the sentence as
// a whole, adn returns an Analysis struct
func SentimentAnalysis(sentence string) *Analysis {
	analysis := &Analysis{
		Words: []Score{},
		Score: 0.0,
	}

	sentence = Clean(sentence)

	w := strings.Split(sentence, " ")
	for _, word := range w {
		analysis.Words = append(analysis.Words, Score{
			Word:  word,
			Score: SentimentOfWord(word),
		})
	}

	analysis.Score = SentimentOfSentence(sentence)

	return analysis
}

// Clean takes in a string (it says sentence
// but it doesn't _have_ to be a sentence) and
// deletes punctionation, foreign characters,
// etc., returning an all lowercase, sanitized
// string. Beware that this ignores and errors
// within the transformation to be a
// deterministic, drop in solution
func Clean(sentence string) string {
	sanitized, _, _ := transform.String(sanitize, sentence)

	return strings.ToLower(sanitized)
}
