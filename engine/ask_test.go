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
)

func getGenericInterrogationParameters() InterrogationParameters {
	ip := InterrogationParameters{
		interactive: false,
		wait:        1 * time.Millisecond,
		limit:       10,
		mode:        linear,
		qachan:      make(chan string),
		command:     make(chan string),
		publisher:   make(chan string),
	}
	return ip
}

func getGenericUnattendedInterrogationParameters() InterrogationParameters {
	ip := getGenericInterrogationParameters()
	ip.interactive = false
	return ip
}

func getGenericInteractiveInterrogationParameters() InterrogationParameters {
	ip := getGenericInterrogationParameters()
	ip.interactive = true
	return ip
}

// TestAskQuestionsInUnattendedAndRandomMode tests that the random mode
// works. How can it be possible since it is random ? A simple way to check
// this is that we are not linear. A mode complex way would be to get the
// index of each question and check the distribution. This would be quite
// complex test and checking the code show that we use the random function
// if we are not linear. So this strategy should be enough.
func TestAskQuestionsInUnattendedAndRandomMode(t *testing.T) {

	r := strings.NewReader(getSampleCsvAsStream())
	tpp := getTpp()
	topic := ParseTopic(r, tpp)

	pr, pw := io.Pipe()
	defer pw.Close()
	ip := getGenericUnattendedInterrogationParameters()
	ip.out = pw
	ip.mode = random

	fmt.Println("    ****************")
	fmt.Println("Test Ask Question in Linear Mode...")
	fmt.Println("    ****************")

	questionsSet := topic.BuildQuestionsSet()
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		AskQuestions(questionsSet, ip)
		pw.Close()
	}()

	s := bufio.NewScanner(pr)

	validateRandomOutput(tpp, questionsSet, *s, t, ip.reversed)
}

// TestAskQuestions tests that, in case of linear run of the questions,
// you get the good questions and good answers that respects the requested
// order.
func TestAskQuestionsInUnattendedMode(t *testing.T) {

	r := strings.NewReader(getSampleCsvAsStream())
	tpp := getTpp()
	topic := ParseTopic(r, tpp)

	pr, pw := io.Pipe()
	defer pw.Close()
	ip := getGenericUnattendedInterrogationParameters()
	ip.out = pw

	fmt.Println("    ****************")
	fmt.Println("Test Ask Question in Linear Mode...")
	fmt.Println("    ****************")

	questionsSet := topic.BuildQuestionsSet()
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		AskQuestions(questionsSet, ip)
		pw.Close()
	}()

	s := bufio.NewScanner(pr)

	validateOutput(tpp, questionsSet, *s, t, ip.reversed)
}

// TestAskQuestionsInReverseMode tests that, in case of linear and reverse run
// of the questions, you get the good questions and good answers that respects
// the requested order.
func TestAskQuestionsInReverseAndUnattendedMode(t *testing.T) {

	r := strings.NewReader(getSampleCsvAsStream())
	tpp := getTpp()
	topic := ParseTopic(r, tpp)

	pr, pw := io.Pipe()
	ip := getGenericUnattendedInterrogationParameters()
	ip.out = pw
	ip.reversed = true

	questionsSet := topic.BuildQuestionsSet()
	go func() {
		defer pw.Close()
		AskQuestions(questionsSet, ip)
	}()

	fmt.Println("    ****************")
	fmt.Println("    Test Ask Question in Linear Reversed Mode...")
	fmt.Println("    ****************")

	s := bufio.NewScanner(pr)

	validateOutput(tpp, questionsSet, *s, t, ip.reversed)
}

//
func validateOutput(tpp TopicParsingParameters, questionsSet QuestionsAnswers, s bufio.Scanner, t *testing.T, reverseMode bool) {

	announcement, _ := regexp.Compile("^" + tpp.TopicAnnounce)
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
		isEmpty = emptyLine.MatchString(s.Text())
		isLoop = loop.MatchString(s.Text())
		isSeparator = separator.MatchString(s.Text())
		isNbOfQ = nbOfQuestions.MatchString(s.Text())
		isLimitReached = limitReached.MatchString(s.Text())
		if !isAnnounce && !isEmpty && !isLoop && !isSeparator && !isNbOfQ && !isLimitReached {
			// default is non reverse mode
			expected = questionsSet.questions[i] + "     --> " + questionsSet.answers[i]
			if reverseMode {
				expected = questionsSet.answers[i] + "     --> " + questionsSet.questions[i]
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
func validateRandomOutput(tpp TopicParsingParameters, questionsSet QuestionsAnswers, s bufio.Scanner, t *testing.T, reverseMode bool) {

	announcement, _ := regexp.Compile("^" + tpp.TopicAnnounce)
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
		isEmpty = emptyLine.MatchString(s.Text())
		isLoop = loop.MatchString(s.Text())
		isSeparator = separator.MatchString(s.Text())
		isNbOfQ = nbOfQuestions.MatchString(s.Text())
		isLimitReached = limitReached.MatchString(s.Text())
		if !isAnnounce && !isEmpty && !isLoop && !isSeparator && !isNbOfQ && !isLimitReached {
			// default is non reverse mode
			expected = questionsSet.questions[i] + "     --> " + questionsSet.answers[i]
			if reverseMode {
				expected = questionsSet.answers[i] + "     --> " + questionsSet.questions[i]
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

	r := strings.NewReader(getSampleCsvAsStream())
	tpp := getTpp()
	topic := ParseTopic(r, tpp)

	pr, pw := io.Pipe()
	userIn, userOut := io.Pipe()
	ip := getGenericInteractiveInterrogationParameters()
	ip.in = userIn
	ip.out = pw

	fmt.Println("    ****************")
	fmt.Println("    Test Ask Question in Linear Mode (pseudo-interactive)...")
	fmt.Println("    ****************")

	questionsSet := topic.BuildQuestionsSet()
	questionsCount := questionsSet.GetCount()

	go func() {
		defer pw.Close()
		AskQuestions(questionsSet, ip)
	}()

	go func() {
		// Simulation of interactive mode: the "user" sends return
		// carriage to command.
		for i := 0; i < ip.limit*questionsCount; i++ {
			fmt.Fprintf(userOut, "\n")
		}
	}()

	s := bufio.NewScanner(pr)

	validateOutput(tpp, questionsSet, *s, t, ip.reversed)

}
