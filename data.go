package sentiment

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/text/transform"

	//"github.com/cdipaolo/goml/cluster"
)

var (
	// words is the main dictionary
	// we're using to base our features
	// from
	words map[string]Word
	dict  map[string]int
	count int

	// filepaths
	pos string
	neg string

	sanitize transform.Transformer
)

func init() {
	words = make(map[string]Word)
	dict = make(map[string]int)
	count = 0

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

	// repeatable!
	rand.Seed(42)
}

// Model encapsulates a the data required
// to run sentiment prediction
type Model struct {
	Words map[string]Word
}

// Word represents the amount of times
// a word has been seen within the data
// and its associated modifier
type Word struct {
	Val         int64
	Count       int64
	Probability float64
}

func calcProbabilities() {
	now := time.Now()
	fmt.Printf("Starting calculating probabilities at %v\n", now)

	for i := range words {
		tmp := words[i]
		tmp.Probability = float64(tmp.Val) / float64(tmp.Count)
		words[i] = tmp
	}

	fmt.Printf("Ended calculating probabilities at %v\n", time.Now().Sub(now))
}

// parseDirToData takes in a directory path and walks through the directory
// to find data (in the expected format given by the dataset in
// datasets/train/* and datasets/test/*.) Because positive and
// negative examples are in separate directories, modifier is
// predetermined and will be the value of 'y' for each example.
func parseDirToData(dirpath string, modifier int64) error {
	now := time.Now()
	fmt.Printf("Starting munging from < %v > to data at %v\n", dirpath, now)

	var ct int
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

				if _, ok := words[word]; ok {
					tmp := words[word]
					tmp.Val += modifier
					tmp.Count++

					words[word] = tmp
				} else {
					words[word] = Word{
						Val:   modifier,
						Count: 1,
					}
				}
			}
		}

		return nil
	})

	if ct == 0 {
		return fmt.Errorf("ERROR: could not find directory!\n\t%v\n", err)
	}
	fmt.Printf("\nFinished munging from < %v > to data\n\tdelta: %v\n\taverage time per line: %v\n", dirpath, time.Now().Sub(now), time.Now().Sub(now)/time.Duration(int64(ct)))

	return err
}

// addWordsToGlobalMap adds words contained within
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
	fmt.Printf("Starting adding words from < %v > to word bank at %v\n", dictpath, now)

	for scanner.Scan() {
		// remove all punctuation and convert to lower case
		text := parseLineToText(scanner.Text())

		if len(text) < 3 {
			continue
		}

		// add word to map if the word is present
		if _, there := dict[text]; !there {
			dict[text] = count
			count++
		}
	}
	fmt.Printf("Finished adding words from < %v > to bank\n\tdelta: %v\n", dictpath, time.Now().Sub(now))
}

// parseLineToText takes in a string, converts
// it to lowercase, removes anything not in a-z,
// and returns that
func parseLineToText(s string) string {
	sanitized, _, _ := transform.String(sanitize, s)

	return strings.TrimSuffix(strings.ToLower(sanitized), "\n")
}

// PersistToFile persists a words summary map to
// a filepath, returning any errors
func PersistToFile(words map[string]Word, path string) error {
	if path == "" {
		return fmt.Errorf("ERROR: you just tried to persist your model to a file with no path!! That's a no-no. Try it with a valid filepath")
	}

	bytes, err := json.Marshal(words)
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
func RestoreFromFile(path string) (map[string]Word, error) {
	if path == "" {
		return nil, fmt.Errorf("ERROR: you just tried to restore your model from a file with no path! That's a no-no. Try it with a valid filepath")
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	words := make(map[string]Word)

	err = json.Unmarshal(bytes, &words)
	if err != nil {
		return nil, err
	}

	return words, nil
}
