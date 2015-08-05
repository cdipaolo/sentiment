package sentiment

import "encoding/json"

// Restore restores a pre-trained models from
// a binary asset this is the preferable method
// of generating a model (use it unless you want
// to train the model again)
//
// This basically wraps RestoreModels.
func Restore() (*Model, error) {
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
	err := json.Unmarshal(bytes, &model)
	if err != nil {
		return nil, err
	}
	return models, nil
}

// Train takes in a directory path to persist the model
// to, trains the model, and saves the model to
// the given file. After this is run you can
// run the SentimentXXX functions effectively.
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
