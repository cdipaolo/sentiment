package sentiment

import "testing"

func init() {
	var err error

	model, err = Train(".")

	//model, err = Restore()
	if err != nil {
		panic(err.Error())
	}
}

func TestPositiveWordSentimentShouldPass1(t *testing.T) {
	t.Parallel()

	w := []string{"happy", "love", "happiness", "humanity", "awesome", "great", "fun", "super", "trust", "fearless", "creative", "dream", "good", "compassion", "joy", "independent", "success"}
	for _, word := range w {
		s := model.SentimentOfWord(Clean(word))
		if s <= 0.5 {
			t.Errorf("Sentiment of < %v > (returned %v) should be greater than 0.5!\n", word, s)
		} else {
			t.Logf("Sentiment of < %v > valid\n\tReturned %v\n", word, s)
		}
	}

	if model.SentimentOfWord("awesome") < model.SentimentOfWord("happy") {
		t.Error("Sentiment of word < 'awesome' > should be greater than sentiment of word < 'happy' >")
	}

	if model.SentimentOfWord("happy") < model.SentimentOfWord("decent") {
		t.Error("Sentiment of word < 'happy' > should be greater than sentiment of word < 'decent' >")
	}
}

func TestNegativeWordSentimentShouldPass1(t *testing.T) {
	t.Parallel()

	w := []string{"not", "resent", "deplorable", "bad", "terrible", "hate", "scary", "terrible", "concerned", "wrong", "rude!!!", "sad", "horrible", "unimpressed", "useless", "offended", "disrespectful"}
	for _, word := range w {
		s := model.SentimentOfWord(Clean(word))
		if s >= 0.5 {
			t.Errorf("Sentiment of < %v > (returned %v) should be less than 0.5!\n", word, s)
		} else {
			t.Logf("Sentiment of < %v > valid\n\tReturned %v\n", word, s)
		}
	}

	if model.SentimentOfWord("hate") < model.SentimentOfWord("resent") {
		t.Error("Sentiment of word < 'hate' > should be greater than sentiment of word < 'resent' >")
	}

	if model.SentimentOfWord("resent") < model.SentimentOfWord("bad") {
		t.Error("Sentiment of word < 'bad' > should be greater than sentiment of word < 'bad' >")
	}
}

func TestPositiveSentenceSentimentShouldPass1(t *testing.T) {
	t.Parallel()

	w := []string{
		"I had an awesome time watching this movie",
		"Sometimes I like to say hello to strangers and it's fun",
		"America needs to support the middle class",
		"Harry Potter is a great movie!",
		"The quest for love is a long one, but it ends in happiness",
		"You are a great person",
		"I love the way you can't talk",
		"You are a caring person",
		"I'm quite ambitious, and this job would be a great opportunity for me!",
		"I'm pretty easy-going.",
		"I find it easy to get along with people",
		"I am very hard-working",
		"I'm very methodical and take care over my work",
	}

	for _, sentence := range w {
		s := model.SentimentOfSentence(Clean(sentence))
		if s <= 0.5 {
			t.Errorf("Sentiment of sentence < %v > (returned %v) should be greater than 0.5!\n", sentence, s)
		} else {
			t.Logf("Sentiment of sentence < %v > is valid.\n\tReturned %v\n", sentence, s)
		}
	}
}

func TestNegativeSentenceSentimentShouldPass1(t *testing.T) {
	t.Parallel()

	w := []string{
		"Jeffery is not a fun guy",
		"I don't enjoy saying hello to strangers",
		"I would compare this person to Donald Trump (ARGH!!!!!) Blind and ignorant!",
		"I'm happy here. I think so, at least.",
		"I hate random people",
		"I don't like your tone right now",
		"I'm not sure you know what you are talking about",
		"The rapture is upon us! Behold!!",
		"I think the growing consensus that China is somehow not a fair player is a bad thing overall",
		"Michelle Bachmann is a total idiot!",
		"How could you say such a thing!",
		"I hate banannas almost as much as I don't love you",
		"Dinner last night sucked",
	}

	for _, sentence := range w {
		s := model.SentimentOfSentence(Clean(sentence))
		if s >= 0.5 {
			t.Errorf("Sentiment of sentence < %v > (returned %v) should be less than 0.5!\n", sentence, s)
		} else {
			t.Logf("Sentiment of sentence < %v > is valid.\n\tReturned %v\n", sentence, s)
		}
	}
}

func TestSentimentAnalysisShouldPass1(t *testing.T) {
	t.Parallel()
	transcript := `On the cross to put away sin by the sacrifice of himself told ...so infinite are on this is so great that only the sacrifice of jesus christ god son could pay for the enormously of arson thank god. He said his son to die for your ...you could be  and blameless before this is john macarthur praying you're continuing to be corporate Now let's get a check of traffic with charlie simon's ...into chaos by a traffic ...center got a problem ...five northbound if your past ninety nine ...just went past the airport and you're ...the woodland watch out in the left lane we have an accident at old river road traffic is backed up the vietnam veterans memorial bridge so far and getting slower by the second incident cleared fifty eastbound eldorado hills boulevard that's good news ninety nine southbound shoulders blocked by an accident at ...and ...capital city freeway the business eighty portion it's got its usual stop-and-go happening right about ...street until you get past E street that ...driving arbitrarily simon's seven ten K FI a joins us now for basic gospel with bob davis and richard piper recorded earlier for broadcasted this time on seven ten K ...keeping faith in america.</p> Fellow everybody with richard piper I bought a and this is basic gossip of ...dedicated to helping ...loop ...good ... If you level bible question or ...because we ... ...recall right ...number's eight four three two seven four two eight four three two seven four two we're lives off the air and online bright out you can get ...could answer your question or discuss the issues ...important in your life so we both of those fault line in the toll-free number four three two seven four two we'd love to hear from you right now it's basic gospel everybody now here's richard piper thanks bob almost long we've been the studying the< idea of freedom why it's so important how do you find it what's the source of it I don't think we could talk about freedom enough because everything in`

	analysis := model.SentimentAnalysis(transcript)

	if analysis.Score <= 0.5 {
		t.Errorf("Analysis of transcript should be greater than 0.5\n\treturned %v\n", analysis.Score)
	} else {
		t.Logf("Analysis of transcript was valid\n\treturned %v\n", analysis.Score)
	}

	if analysis.Words == nil {
		t.Errorf("Analysis of transcript returned nil words array!\n")
	} else {
		t.Logf("Analysis of transcript retuned valid word array\n\t returned %+v\n", analysis.Words)
	}
}

// Paul Graham essay on immigration laws is _slightly_ negative
func TestSentimentAnalysisShouldPass2(t *testing.T) {
	t.Parallel()
	transcript := `The anti-immigration people have to invent some explanation to account for all the effort technology companies have expended trying to make immigration easier. So they claim it's because they want to drive down salaries. But if you talk to startups, you find practically every one over a certain size has gone through legal contortions to get programmers into the US, where they then paid them the same as they'd have paid an American. Why would they go to extra trouble to get programmers for the same price? The only explanation is that they're telling the truth: there are just not enough great programmers to go around`

	analysis := model.SentimentAnalysis(transcript)

	if analysis.Score >= 0.5 {
		t.Errorf("Analysis of transcript should be greater than 0.5\n\treturned %v\n", analysis.Score)
	} else {
		t.Logf("Analysis of transcript was valid\n\treturned %v\n", analysis.Score)
	}

	if analysis.Words == nil {
		t.Errorf("Analysis of transcript returned nil words array!\n")
	} else {
		t.Logf("Analysis of transcript retuned valid word array\n\t returned %+v\n", analysis.Words)
	}
}

// Donald Trump snippet from his annoucement presidential speech
func TestAssholeSentimentAnalysisShouldPass1(t *testing.T) {
	t.Parallel()
	transcript := `Thank you. It's true, and these are the best and the finest. When Mexico sends its people, they're not sending their best. They're not sending you. They're not sending you. They're sending people that have lots of problems, and they're bringing those problems with us. They're bringing drugs. They're bringing crime. They're rapists. And some, I assume, are good people. But I speak to border guards and they tell us what we're getting. And it only makes common sense. It only makes common sense. They're sending us not the right people. It's coming from more than Mexico. It's coming from all over South and Latin America, and it's coming probably -- probably -- from the Middle East. But we don't know. Because we have no protection and we have no competence, we don't know what's happening. And it's got to stop and it's got to stop fast.`

	analysis := model.SentimentAnalysis(transcript)

	if analysis.Score >= 0.5 {
		t.Errorf("Analysis of transcript should be greater than 0.5\n\treturned %v\n", analysis.Score)
	} else {
		t.Logf("Analysis of transcript was valid\n\treturned %v\n", analysis.Score)
	}

	if analysis.Words == nil {
		t.Errorf("Analysis of transcript returned nil words array!\n")
	} else {
		t.Logf("Analysis of transcript retuned valid word array\n\t returned %+v\n", analysis.Words)
	}
}
