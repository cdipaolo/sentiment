package sentiment

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cdipaolo/goml/base"
	"github.com/cdipaolo/goml/text"
)

// TrainEnglishModel takes in a path to the expected
// IMDB datasets, and a map of models to add the model
// to. It'll return any errors if there were any.
func TrainEnglishModel(modelMap Models) error {
	pos, posErr := filepath.Abs("datasets/train/pos")
	neg, negErr := filepath.Abs("datasets/train/neg")
	if posErr != nil || negErr != nil {
		return fmt.Errorf("Error getting the IMDB English review dataset from expected paths. Are you in the project directory???\n\tNegative Sample Error: %v\n\tPositive Sample Error: %v\n", negErr, posErr)
	}

	var ct int

	stream := make(chan base.TextDatapoint, 1000)
	errors := make(chan error, 100)
	model := text.NewNaiveBayes(stream, 2, base.OnlyWords)

	go model.OnlineLearn(errors)

	walk := func(dirpath string) error {
		var class uint8
		if strings.Contains(dirpath, "pos") {
			class = uint8(1)
		}

		return filepath.Walk(dirpath, func(path string, info os.FileInfo, err error) error {
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
				line := strings.TrimSuffix(scanner.Text(), "\n")

				ct++
				if ct%500 == 0 {
					print(".")
				}

				stream <- base.TextDatapoint{
					X: line,
					Y: class,
				}
			}

			return nil
		})
	}

	now := time.Now()
	fmt.Printf("Starting munging from < %v, %v > to data at %v\n", pos, neg, now)

	posErr = walk(pos)
	negErr = walk(neg)
	if posErr != nil || negErr != nil {
		return fmt.Errorf("Error training english sentiment model!\n\tPositive Error: %v\n\tNegative Error: %v\n", posErr, negErr)
	}

	close(stream)

	for {
		err, more := <-errors
		if more {
			fmt.Printf("Error passed: %v\n", err)
		} else {
			// training is done!
			break
		}
	}

	modelMap[English] = model

	fmt.Printf("\nFinished munging from < %v, %v > to data\n\tdelta: %v\n\taverage time per line: %v\n", pos, neg, time.Now().Sub(now), time.Now().Sub(now)/time.Duration(int64(ct)))

	return nil
}
