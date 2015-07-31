package sentiment

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/text/transform"

	//"github.com/cdipaolo/goml/cluster"
)

var (
	// model is the model to be trained
	model *Model

	dict map[string]int64

	// filepaths (assume user is in
	// current directory)
	pos string
	neg string

	sanitize transform.Transformer
)

func init() {
	model = &Model{
		Words: make(map[string]Word),
	}

	dict = make(map[string]int64)

	var posErr, negErr error

	pos, posErr = filepath.Abs("datasets/train/pos")
	neg, negErr = filepath.Abs("datasets/train/neg")
	if posErr != nil || negErr != nil {
		panic(fmt.Sprintf("Cannot find absolute paths from relative paths.\n\tNegative Path Error: %v\n\tPositive Path Error: %v\n", negErr, posErr))
	}

	sanitize = transform.RemoveFunc(func(r rune) bool {
		switch {
		case r >= 'A' && r <= 'Z':
			return false
		case r >= 'a' && r <= 'z':
			return false
		case r >= '0' && r <= '1':
			return false
		case r == ' ':
			return false
		case r == '\t':
			return false
		default:
			return true
		}
	})
}

// Model encapsulates a the data required
// to run sentiment prediction
type Model struct {
	// Word mappings within the trained model
	Words map[string]Word

	// Gives the number of times
	// an example was either positive
	// or negative, or either
	CountPositive int64
	CountNegative int64
	Count         int64

	// The number of words in our
	// dictionary (pun sort of
	// intended :)
	DictSize int64

	// Gives the pre-calculated
	// probability that Y is 1
	// or Y is 0
	PYIsOne  float64
	PYIsZero float64
}

// Word represents the amount of times
// a word has been seen within the data
// and its associated modifier
type Word struct {
	// Gives the number of times X is found
	// in a positive example, a negative
	// example, or any example
	CountPositive int64
	CountNegative int64

	// Gives the probability that
	// an x is found given y is
	// either 0, or y is 1
	ProbabilityXIsOne  float64
	ProbabilityXIsZero float64
}

func calcProbabilities() {
	now := time.Now()
	fmt.Printf("Starting calculating probabilities at %v\n", now)

	model.DictSize = int64(len(model.Words))

	model.PYIsOne = float64(model.CountPositive) / float64(model.Count)
	model.PYIsZero = float64(model.CountNegative) / float64(model.Count)

	for i := range model.Words {
		tmp := model.Words[i]
		tmp.ProbabilityXIsOne = (float64(tmp.CountPositive) + 1) / (float64(model.CountPositive) + float64(model.DictSize))
		tmp.ProbabilityXIsZero = (float64(tmp.CountNegative) + 1) / (float64(model.CountNegative) + float64(model.DictSize))
		model.Words[i] = tmp
	}

	fmt.Printf("Ended calculating probabilities at %v\n", time.Now().Sub(now))
}

// parseDirToData takes in a directory path and walks through the directory
// to find data (in the expected format given by the dataset in
// datasets/train/* and datasets/test/*.) Because positive and
// negative examples are in separate directories, modifier is
// predetermined and will be the value of 'y' for each example.
func parseDirToData(dirpath string) error {
	now := time.Now()
	fmt.Printf("Starting munging from < %v > to data at %v\n", dirpath, now)

	var ct int

	positive := strings.Contains(dirpath, "pos")

	err := filepath.Walk(dirpath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// read file
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		scanner.Split(bufio.ScanLines)

		// only need one line
		if scanner.Scan() {
			scanText := scanner.Text()
			line := parseLineToText(scanText)
			text := strings.Split(line, " ")

			ct++
			if ct%500 == 0 {
				print(".")
			}

			for _, word := range text {
				if len(word) < 3 {
					continue
				}

				if _, ok := model.Words[word]; ok {
					tmp := model.Words[word]
					if positive {
						tmp.CountPositive++
					} else {
						tmp.CountNegative++
					}

					model.Words[word] = tmp
				} else if positive {
					model.Words[word] = Word{
						CountPositive: int64(1),
					}
				} else {
					model.Words[word] = Word{
						CountNegative: int64(1),
					}
				}
			}
		}

		if positive {
			model.CountPositive++
		} else {
			model.CountNegative++
		}

		model.Count++
		return nil
	})

	if ct == 0 {
		return fmt.Errorf("ERROR: could not find directory!\n\t%v\n", err)
	}
	fmt.Printf("\nFinished munging from < %v > to data\n\tdelta: %v\n\taverage time per line: %v\n", dirpath, time.Now().Sub(now), time.Now().Sub(now)/time.Duration(int64(ct)))

	return err
}

// addmodel.WordsToGlobalMap adds model.Words contained within
// a dictionary file to memory
func addWordsToGlobalMap(dictpath string) {
	f, err := os.Open(dictpath)
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()

	// create word bank
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	now := time.Now()
	fmt.Printf("Starting adding model.Words from < %v > to word bank at %v\n", dictpath, now)

	for scanner.Scan() {
		// remove all punctuation and convert to lower case
		text := parseLineToText(scanner.Text())

		if len(text) < 3 {
			continue
		}

		// add word to map if the word is present
		if _, there := dict[text]; !there {
			dict[text] = int64(len(dict))
		}
	}
	fmt.Printf("Finished adding model.Words from < %v > to bank\n\tdelta: %v\n", dictpath, time.Now().Sub(now))
}

// parseLineToText takes in a string, converts
// it to lowercase, removes anything not in a-z,
// and returns that
func parseLineToText(s string) string {
	sanitized, _, _ := transform.String(sanitize, s)

	return strings.TrimSuffix(strings.ToLower(sanitized), "\n")
}

// PersistToFile persists a model.Words summary map to
// a filepath, returning any errors
func PersistToFile(m Model, path string) error {
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

// RestoreFromFile resores a word summary map
// (probably saved with PersistToFile) into
// memory, returning any errors
func RestoreFromFile(path string) (*Model, error) {
	if path == "" {
		return nil, fmt.Errorf("ERROR: you just tried to restore your model from a file with no path! That's a no-no. Try it with a valid filepath")
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	model := Model{}

	err = json.Unmarshal(bytes, &model)
	if err != nil {
		return nil, err
	}

	return &model, nil
}
