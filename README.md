## Sentiment
#### Simple Sentiment Analysis in Golang

Sentiment lets you pass strings into a function and get an estimate of the sentiment of the string (in english) using a very simple probabalistic model. The model is trained off of [this dataset](https://inclass.kaggle.com/c/si650winter11/data) which is a collection of IMDB movie reviews classified by sentiment. The return values of the word sentiment is on the interval [-1,1]. Sentence/text classification can vary but tends to stay close to 0. For sentences, each word probability can be modified by using the words "don't", "not", etc. before the word. Again, this is not some advanced SVM, Neural Net, or Bayesian Classifier.

### Model

This is effectively ((# of times the word x is in a positive) - (# of times the word x is in a negative sentence) / (total # of instances of the word x)) though the modifier could be different. For sentences we are differing only by checking for negation words ('dont', 'not', etc.) and amplification words ('really') to modify overall score. For the most part sentence analysis just adds together the predictions of the words. Words (for both sentence and word based prediction) with a length of less than 3 characters are ignored.

I'm saying h(y|x) instead of P(y|x) because the estimate could be negative as well

```
h(sentiment is positive | a word x) = Σ1{x = x[i]}*y[i] / Σ1{x = x[i]}
```

### Example

You can save the model trained off of the dataset to a json file using the `PersistToFile(filepath string) error` function so you don't have to run the training again, though it only takes about 4 seconds max.

Training (this only needs to be done once, and should take less than 5 seconds):
```go
model, err := sentiment.Train("dir/to/save/to")
if err != nil {
    panic(fmt.Sprintf("Could not persist words to JSON!\n\t%v\n", err))
}
```

Analysis:
```go
// sanitize input
cleaned := sentiment.Clean("MaKe Th1s Into LOWER c-----ase and tak3 out numb3rs, etc.")

// get word sentiment
s := model.SentimentOfWord("love") // greater than 0
s = model.SentimentOfWord("hate") // less than 0

// get sentence sentiment
s = model.SentimentOfSentence(sentiment.Clean("I had a great day!!!")) // greater than 0

// get sentiment analysis (sentiment for every word as well as overall score)
//
// this is the 'plug and play' analysis, pretty much.
analysis = model.SentimentAnalysis("this is a D3rty sentence that will get cleaned prior to being evaluated")
```

### LICENSE - MIT
