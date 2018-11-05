package engine

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/boris-lenzinger/repeatit/datamodel"
	"github.com/boris-lenzinger/repeatit/internal/tests"
	"github.com/boris-lenzinger/repeatit/parsing"
)

func getGenericInterrogationParameters() datamodel.InterrogationParameters {
	ip := datamodel.NewInterrogationParameters()
	ip.SetLinearMode()
	ip.SetLimit(10)
	ip.SetPauseTime(1 * time.Millisecond)

	return ip
}

func getGenericUnattendedInterrogationParameters() datamodel.InterrogationParameters {
	return getGenericInterrogationParameters()
}

func getGenericInteractiveInterrogationParameters() datamodel.InterrogationParameters {
	ip := getGenericInterrogationParameters()
	ip.SetInteractive()
	return ip
}

// TestAskQuestionsInUnattendedAndRandomMode tests that the random mode
// works. How can it be possible since it is random ? A simple way to check
// this is that we are not linear. A mode complex way would be to get the
// index of each question and check the distribution. This would be quite
// complex test and checking the code show that we use the random function
// if we are not linear. So this strategy should be enough.
func TestAskQuestionsInUnattendedAndRandomMode(t *testing.T) {

	r := strings.NewReader(tests.GetSampleCsvAsStream())
	tpp := datamodel.NewTopicParsingParameters()
	topic, err := parsing.ParseTopic(r, tpp)
	if err != nil {
		t.Fatalf("sample csv must be parsed with no error. Got the following: %v", err)
	}

	pr, pw := io.Pipe()
	defer pw.Close()
	ip := getGenericUnattendedInterrogationParameters()
	ip.SetOutputStream(pw)
	ip.SetRandomMode()

	fmt.Println("    ****************")
	fmt.Println("Test Ask Question in Linear Mode...")
	fmt.Println("    ****************")

	questionsSet := topic.BuildVocabularyQuestionsSet()
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		AskQuestions(questionsSet, ip)
		pw.Close()
	}()

	s := bufio.NewScanner(pr)

	validateRandomOutput(tpp, questionsSet, *s, t, ip.IsReversedMode())
}

// TestAskQuestions tests that, in case of linear run of the questions,
// you get the good questions and good answers that respects the requested
// order.
func TestAskQuestionsInUnattendedMode(t *testing.T) {

	r := strings.NewReader(tests.GetSampleCsvAsStream())
	tpp := datamodel.NewTopicParsingParameters()
	topic, err := parsing.ParseTopic(r, tpp)
	if err != nil {
		t.Fatalf("sample csv must be parsed with no error. Received: %v", err)
	}

	pr, pw := io.Pipe()
	defer pw.Close()
	ip := getGenericUnattendedInterrogationParameters()
	ip.SetOutputStream(pw)

	fmt.Println("    ****************")
	fmt.Println("Test Ask Question in Linear Mode...")
	fmt.Println("    ****************")

	questionsSet := topic.BuildVocabularyQuestionsSet()
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		AskQuestions(questionsSet, ip)
		pw.Close()
	}()

	s := bufio.NewScanner(pr)

	validateOutput(tpp, questionsSet, *s, t, ip.IsReversedMode())
}

// TestAskQuestionsInReverseMode tests that, in case of linear and reverse run
// of the questions, you get the good questions and good answers that respects
// the requested order.
func TestAskQuestionsInReverseAndUnattendedMode(t *testing.T) {

	r := strings.NewReader(tests.GetSampleCsvAsStream())
	tpp := datamodel.NewTopicParsingParameters()
	topic, err := parsing.ParseTopic(r, tpp)
	if err != nil {
		t.Fatalf("parsing sample csv must not return an error. Received: %v", err)
	}

	pr, pw := io.Pipe()
	ip := getGenericUnattendedInterrogationParameters()
	ip.SetOutputStream(pw)
	ip.SetReverseMode()

	questionsSet := topic.BuildVocabularyQuestionsSet()
	go func() {
		defer pw.Close()
		AskQuestions(questionsSet, ip)
	}()

	fmt.Println("    ****************")
	fmt.Println("    Test Ask Question in Linear Reversed Mode...")
	fmt.Println("    ****************")

	s := bufio.NewScanner(pr)

	validateOutput(tpp, questionsSet, *s, t, ip.IsReversedMode())
}

//
func validateOutput(tpp datamodel.TopicParsingParameters, questionsSet datamodel.QuestionsAnswers, s bufio.Scanner, t *testing.T, reverseMode bool) {
	announcement := regexp.MustCompile("^" + tpp.LessonAnnounce)
	questionsCount := questionsSet.GetCount()

	i := 0
	var (
		isAnnounce     bool
		isEmpty        bool
		isLoop         bool
		isSeparator    bool
		isNbOfQ        bool
		isLimitReached bool
		expected       string
		computed       string
	)
	for s.Scan() {
		isAnnounce = announcement.MatchString(s.Text())
		isEmpty = tests.EmptyLine.MatchString(s.Text())
		isLoop = tests.Loop.MatchString(s.Text())
		isSeparator = tests.Separator.MatchString(s.Text())
		isNbOfQ = tests.NbOfQuestions.MatchString(s.Text())
		isLimitReached = tests.LimitReached.MatchString(s.Text())
		if !isAnnounce && !isEmpty && !isLoop && !isSeparator && !isNbOfQ && !isLimitReached {
			// default is non reverse mode
			expected = questionsSet.GetQuestion(i) + "     --> " + questionsSet.GetAnswer(i)
			if reverseMode {
				expected = questionsSet.GetAnswer(i) + "     --> " + questionsSet.GetQuestion(i)
			}
			computed = s.Text()
			if computed != expected {
				t.Errorf("Check of answers failed. We were expected '%s' but received '%s'\n", expected, computed)
			}
			i = (i + 1) % questionsCount
		}
	}
}

//
func validateRandomOutput(tpp datamodel.TopicParsingParameters, questionsSet datamodel.QuestionsAnswers, s bufio.Scanner, t *testing.T, reverseMode bool) {

	announcement, _ := regexp.Compile("^" + tpp.LessonAnnounce)
	questionsCount := questionsSet.GetCount()
	i := 0
	matchCount := 0
	totalTests := 0
	var (
		isAnnounce     bool
		isEmpty        bool
		isLoop         bool
		isSeparator    bool
		isNbOfQ        bool
		isLimitReached bool
		expected       string
		computed       string
	)
	for s.Scan() {
		isAnnounce = announcement.MatchString(s.Text())
		isEmpty = tests.EmptyLine.MatchString(s.Text())
		isLoop = tests.Loop.MatchString(s.Text())
		isSeparator = tests.Separator.MatchString(s.Text())
		isNbOfQ = tests.NbOfQuestions.MatchString(s.Text())
		isLimitReached = tests.LimitReached.MatchString(s.Text())
		if !isAnnounce && !isEmpty && !isLoop && !isSeparator && !isNbOfQ && !isLimitReached {
			// default is non reverse mode
			expected = questionsSet.GetQuestion(i) + "     --> " + questionsSet.GetQuestion(i)
			if reverseMode {
				expected = questionsSet.GetAnswer(i) + "     --> " + questionsSet.GetQuestion(i)
			}
			computed = s.Text()
			if computed == expected {
				matchCount++
			}
			i = (i + 1) % questionsCount
			totalTests++
		}
	}
	if matchCount == totalTests {
		t.Errorf("The random mode does not work since we match the linear scenario.")
	}
}

// TestAskQuestionsInteractive tests that, in case of linear and interactive
// run of the questions, the user gets the good questions and has to press
// return to get the matching answers and all of this in the requested order.
func TestAskQuestionsInLinearAndInteractiveMode(t *testing.T) {

	r := strings.NewReader(tests.GetSampleCsvAsStream())
	tpp := datamodel.NewTopicParsingParameters()
	topic, err := parsing.ParseTopic(r, tpp)
	if err != nil {
		t.Fatalf("parsing sample csv must not fail. Received: %v", err)
	}

	pr, pw := io.Pipe()
	userIn, userOut := io.Pipe()
	ip := getGenericInteractiveInterrogationParameters()
	ip.SetInputStream(userIn)
	ip.SetOutputStream(pw)

	fmt.Println("    ****************")
	fmt.Println("    Test Ask Question in Linear Mode (pseudo-interactive)...")
	fmt.Println("    ****************")

	questionsSet := topic.BuildVocabularyQuestionsSet()
	questionsCount := questionsSet.GetCount()

	go func() {
		defer pw.Close()
		AskQuestions(questionsSet, ip)
	}()

	go func() {
		// Simulation of interactive mode: the "user" sends return
		// carriage to command.
		for i := 0; i < ip.GetLimit()*questionsCount; i++ {
			fmt.Fprintf(userOut, "\n")
		}
	}()

	s := bufio.NewScanner(pr)

	validateOutput(tpp, questionsSet, *s, t, ip.IsReversedMode())

}
