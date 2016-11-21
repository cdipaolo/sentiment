package sentiment

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/cdipaolo/goml/base"
	"github.com/cdipaolo/goml/text"
)

const (
	// TempDirectory is the default temporary
	// directory for persisting models to disk
	TempDirectory string = "/tmp/.sentiment"
)

// Restore restores a pre-trained models from
// a binary asset this is the preferable method
// of generating a model (use it unless you want
// to train the model again)
//
// This basically wraps RestoreModels.
func Restore() (Models, error) {
	data, err := Asset("model.json")
	if err != nil {
		return nil, fmt.Errorf("Could not restore model from binary asset!\n\t%v\n", err)
	}

	return RestoreModels(data)
}

// RestoreModels takes in a byte of
// a (presumably) map[Language]LanguageModel
// and marshals it into a usable model that
// you can use to run regular, language
// specific sentiment analysis
func RestoreModels(bytes []byte) (Models, error) {
	models := make(Models)
	err := json.Unmarshal(bytes, &models)
	if err != nil {
		return nil, err
	}

	tokenizer := &text.SimpleTokenizer{
		SplitOn: " ",
	}

	for i := range models {
		models[i].UpdateSanitize(base.OnlyWords)
		models[i].UpdateTokenizer(tokenizer)
	}

	return models, nil
}

// PersistToFile persists a Models struct to
// a filepath, returning any errors
func PersistToFile(m Models, path string) error {
	if path == "" {
		return fmt.Errorf("ERROR: you just tried to persist your model to a file with no path!! That's a no-no. Try it with a valid filepath")
	}

	bytes, err := json.Marshal(m)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, bytes, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

// Train takes in a directory path to persist the model
// to, trains the model, and saves the model to
// the given file. After this is run you can
// run the SentimentXXX functions effectively.
//
// Note that this must be run from within the project
// directory! To just get the model without re-training
// you should just call "Resore"
func Train() (Models, error) {
	models := make(Models)
	err := TrainEnglishModel(models)
	if err != nil {
		return nil, fmt.Errorf("Count not train English sentiment model!\n\t%v\n", err)
	}

	err = os.MkdirAll(TempDirectory, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("Count not create temp directory!\n\t%v\n", err)
	}

	err = PersistToFile(models, path.Join(TempDirectory, "model.json"))
	if err != nil {
		return nil, fmt.Errorf("Could not persist m.Words to JSON!\n\t%v\n", err)
	}

	return models, nil
}
