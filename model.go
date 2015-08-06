package sentiment

import (
	"github.com/cdipaolo/goml/text"
)

// Language is a language code used
// for differentiating sentiment models
type Language string

// Constants hold the Twitter language
// codes that will correspond to models.
// Obviously all of these won't be used
// initially, but they're here for
// ease of extention. US English is being
// lumped with UK English.
const (
	English            Language = "en"
	Spanish                     = "es"
	French                      = "fr"
	German                      = "de"
	Italian                     = "it"
	Arabic                      = "ar"
	Japanese                    = "ja"
	Indonesian                  = "id"
	Portugese                   = "pt"
	Korean                      = "ko"
	Turkish                     = "tr"
	Russian                     = "ru"
	Dutch                       = "nl"
	Filipino                    = "fil"
	Malay                       = "msa"
	ChineseTraditional          = "zh-tw"
	ChineseSimplified           = "zh-cn"
	Hindi                       = "hi"
	Norwegian                   = "no"
	Swedish                     = "sv"
	Finnish                     = "fi"
	Danish                      = "da"
	Polish                      = "pl"
	Hungarian                   = "hu"
	Farsi                       = "fa"
	Hebrew                      = "he"
	Urdu                        = "ur"
	Thai                        = "th"
	NoLanguage                  = ""
)

// Models holds a map from language keys
// to sentiment classifiers.
type Models map[Language]*text.NaiveBayes

// Score holds the score of a
// singular word (differs from
// SentenceScore only in param
// names and JSON marshaling, not
// actualy types)
type Score struct {
	Word  string `json:"word"`
	Score uint8  `json:"score"`
}

// SentenceScore holds the score
// of a document, which could be
// (and probably is) a sentence
type SentenceScore struct {
	Sentence string `json:"sentence"`
	Score    uint8  `json:"score"`
}

// Analysis returns the analysis
// of a document, splitting it into
// total sentiment, individual sentence
// sentiment, and individual word
// sentiment, along with the language
// code
type Analysis struct {
	Language  Language        `json:"lang"`
	Words     []Score         `json:"words"`
	Sentences []SentenceScore `json:"sentences,omitempty"`
	Score     uint8           `json:"score"`
}
