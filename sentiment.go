package sentiment

import (
	"encoding/json"
	"fmt"
	"math"
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
	Words []Score `json:"words"`
	Score float64 `json:"score"`
}

// SentimentOfWord takes in a single word and
// returns the probability that the word is
// classified as positive based on the corpus
func (m *Model) SentimentOfWord(word string) float64 {
	if _, ok := m.Words[word]; len(word) < 3 || !ok {
		return 0.5
	}

	one := m.Words[word].ProbabilityXIsOne * m.PYIsOne

	return one / (one + m.Words[word].ProbabilityXIsZero*m.PYIsZero)
}

// SentimentOfSentence takes in a string of a
// space-delimited sentence and returns 0 if
// the sentence is negative and 1 if the sample
// is positive
func (m *Model) SentimentOfSentence(sentence string) float64 {
	// These are the _probability_ products
	// given in the numerator and denominator
	// that a sentence is positive or negative
	// (1 or 0)
	var one float64 = 1.0
	var zero float64 = 1.0

	w := strings.Split(sentence, " ")
	for _, word := range w {
		if _, ok := m.Words[word]; len(word) < 3 || !ok {
			continue
		}

		one += math.Log(m.Words[word].ProbabilityXIsOne)
		zero += math.Log(m.Words[word].ProbabilityXIsZero)
	}

	one += math.Log(m.PYIsOne)
	zero += math.Log(m.PYIsZero)

	if one > zero {
		return 1.0
	}
	return 0.0
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

// Restore restores a pre-trained model from
// a binary asset this is the preferable method
// of generating a model (use it unless you want
// to train the model again)
func Restore() (*Model, error) {
	data, err := Asset("model.json")
	if err != nil {
		return nil, fmt.Errorf("Could not restore model from binary asset!\n\t%v\n", err)
	}

	model := Model{}
	err = json.Unmarshal(data, &model)
	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal stored model!\n\t%v\n", err)
	}

	return &model, nil
}

// Train takes in a directory path to persist the model
// to, trains the model, and saves the model to
// the given file. After this is run you can
// run the Sentiment... functions effectively.
//
// Note that this must be run from within the project
// directory! To just get the model without re-training
// you should just call "Resore"
func Train(dir string) (*Model, error) {
	err := parseDirToData(pos)
	if err != nil {
		return nil, fmt.Errorf("Count not parse directory < %v > to data!\n\t%v\n", pos, err)
	}

	err = parseDirToData(neg)
	if err != nil {
		return nil, fmt.Errorf("Count not parse directory < %v > to data!\n\t%v\n", neg, err)
	}

	calcProbabilities()

	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("Count not create temp directory!\n\t%v\n", err)
	}

	err = PersistToFile(*model, path.Join(dir, "model.json"))
	if err != nil {
		return nil, fmt.Errorf("Could not persist m.Words to JSON!\n\t%v\n", err)
	}

	return model, nil
}
