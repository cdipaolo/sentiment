package sentiment

import (
	"github.com/cdipaolo/goml"
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
)

// Models holds a map from language keys
// to sentiment classifiers.
type Models map[string]base.OnlineTextModel

// LanguageModel holds a goml text model
// to be trained for sentiment. It just
// encapsulates the data stream and
// error stream
type LanguageModel struct {
	stream chan base.TextDatapoint
	errors chan error

	Model base.OnlineTextModel `json:"model"`
}
