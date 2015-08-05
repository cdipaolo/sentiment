## Sentiment
#### Simple, Drop In Sentiment Analysis in Golang
[![GoDoc](https://godoc.org/github.com/cdipaolo/sentiment?status.svg)](https://godoc.org/github.com/cdipaolo/sentiment)

Sentiment lets you pass strings into a function and get an estimate of the sentiment of the string (in english) using a very simple probabalistic model. The model is trained off of [this dataset](https://inclass.kaggle.com/c/si650winter11/data) which is a collection of IMDB movie reviews classified by sentiment. The returned values for single word classification is the corresponding probability that the document is positive on [0,1]. For sentence/document classification the result is returned as a descrete classification on {0,1} corresponding to either negative or positive, respectively.

### Model

Sentiment uses a Naive Bayes classification model for prediction. There are plusses and minuses, but Naive bayes tends to do well for text classification.

### Example

You can save the model trained off of the dataset to a json file using the `PersistToFile(filepath string) error` function so you don't have to run the training again, though it only takes about 4 seconds max.

Training, or Restoring a Pre-Trained Model:
```go
// Train is used within the library, but you should
// usually prefer Restore because it's faster and
// you don't have to be in the project's directory
//
// model, err := sentiment.Train("dir/to/save/to")

model, err := sentiment.Restore()
if err != nil {
    panic(fmt.Sprintf("Could not restore model!\n\t%v\n", err))
}
```

Analysis:
```go
// sanitize input
cleaned := sentiment.Clean("MaKe Th1s Into LOWER c-----ase and tak3 out numb3rs, etc.")

// get word sentiment
s := model.SentimentOfWord("love") // greater than 0.5
s = model.SentimentOfWord("hate") // less than 0.5

// get sentence/document sentiment
s = model.SentimentOfSentence(sentiment.Clean("I had a great day!!!")) // 1

// get sentiment analysis (sentiment for every word as well as overall score)
//
// this is the 'plug and play' analysis, pretty much.
analysis = model.SentimentAnalysis("this is a D3rty sentence that will get cleaned prior to being evaluated")
```

### LICENSE - MIT
