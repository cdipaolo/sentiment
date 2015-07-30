package sentiment

import (
	"fmt"
	"os"
	"path"
	"strings"

	"golang.org/x/text/transform"
)

type Score struct {
	Word  string  `json:"word"`
	Score float64 `json:"score"`
}

type Analysis struct {
	Words []Score `json:"m.Words"`
	Score float64 `json:"score"`
}

// SentimentOfWord takes in a single word and
// returns the sentiment of that word (negative
// would mean negative sentiment, positive would
// mean positive sentiment)
func (m *Model) SentimentOfWord(word string) float64 {
	if m.Words[word].Count < 6 || len(word) < 3 {
		return 0.0
	}

	return m.Words[word].Probability
}

// SentimentOfSentence takes in a string of a
// space-delimited sentence and returns a
// weighted sum of their probabilities
func (m *Model) SentimentOfSentence(sentence string) float64 {
	var sum float64

	w := strings.Split(sentence, " ")
	for i, word := range w {
		s := m.SentimentOfWord(word)

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
func (m *Model) SentimentAnalysis(sentence string) *Analysis {
	analysis := &Analysis{
		Words: []Score{},
		Score: 0.0,
	}

	sentence = Clean(sentence)

	w := strings.Split(sentence, " ")
	for _, word := range w {
		analysis.Words = append(analysis.Words, Score{
			Word:  word,
			Score: m.SentimentOfWord(word),
		})
	}

	analysis.Score = m.SentimentOfSentence(sentence)

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

// Train takes in a directory path to persist the model
// to, trains the model, and saves the model to
// the given file. After this is run you can
// run the Sentiment... functions effectively.
func Train(dir string) (*Model, error) {
	err := parseDirToData(pos, 1)
	if err != nil {
		return nil, fmt.Errorf("Count not parse directory < %v > to data!\n\t%v\n", pos, err)
	}

	err = parseDirToData(neg, -1)
	if err != nil {
		return nil, fmt.Errorf("Count not parse directory < %v > to data!\n\t%v\n", neg, err)
	}

	calcProbabilities()

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("Count not create temp directory!\n\t%v\n", err)
	}

	err = PersistToFile(words, path.Join(dir, "m.Words.json"))
	if err != nil {
		return nil, fmt.Errorf("Could not persist m.Words to JSON!\n\t%v\n", err)
	}

	return &Model{Words: words}, nil
}
