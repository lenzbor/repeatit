package main

import (
	"strconv"
	"strings"
	"testing"
	"time"
)

// TestParsing validates the parsing of the command line.
func TestParsingEmptyParameters(t *testing.T) {
	p, err := Parse()
	if err != nil {
		t.Errorf("Parsing should not fail with empty parameters")
	}
	if p.interactive {
		t.Errorf("Default is to be in non interactive. But the parameters says the contrary.")
	}
	if p.wait != 2*time.Second {
		t.Errorf("Default is to wait for 2 seconds. But the current value is %v.\n", p.wait)
	}
}

// TestParsingNonEmptyParameters checks that passing an interactive mode and a different
// time are supported.
func TestParsingNonEmptyParameters(t *testing.T) {
	wt := 1500
	arguments := []string{"-i", "-t", strconv.Itoa(wt)}
	p, err := Parse(arguments[:]...)
	if err != nil {
		t.Errorf("A valid list of parameters must not trigger a parsing error.")
	}
	if !p.interactive {
		t.Errorf("The parameter -i was not detected.")
	}
	if p.wait != time.Duration(wt)*time.Millisecond {
		t.Errorf("Failed to detect wait time as %dms. Found %v instead.\n", wt, p.wait)
	}
}

// TestParsingSelectedTopics checks that the option -l (picking specific topics)
func TestParsingSelectedTopics(t *testing.T) {
	selected := "Topic 1,Topic 2"
	arguments := []string{"-l", selected}
	p, err := Parse(arguments[:]...)
	if err != nil {
		t.Errorf("Parsing detects list of selected topics as an error")
	}
	if p.subsections != selected {
		t.Errorf("Parsing failed to set the list of selected topics to the string that was passed as parameter.")
	}
	listAsArray := p.GetListOfSubsections()
	if len(listAsArray) != 2 {
		t.Errorf("Retrieving the list of selected topics should have reported 2 elements but we received %d\n", len(listAsArray))
	}
}

// TestNoSelectedTopicsReturnsNil checks that when the user sets no
// specific topics, the array in nil.
func TestNoSelectedTopicsReturnsNil(t *testing.T) {
	arguments := []string{}
	p, err := Parse(arguments[:]...)
	if err != nil {
		t.Errorf("Passing no argument make the parsing fail")
	}
	if p.GetListOfSubsections() != nil {
		t.Errorf("No argument passed but the parameters holds a non nil list of selected topics as '%v'", p.GetListOfSubsections())
	}
}

// TestParsingReverseMode checks that reverse mode is detected and works.
func TestParsingReverseMode(t *testing.T) {
	arguments := []string{"-r"}
	p, err := Parse(arguments[:]...)
	if err != nil {
		t.Errorf("Parsing detects reverse mode as an error")
	}
	if p.IsReversedMode() != true {
		t.Errorf("Parsing failed to set the reverse mode.")
	}
}

// TestParsingSummaryMode checks that the feature about the parameter summary works fine.
func TestParsingSummaryMode(t *testing.T) {
	arguments := []string{"-s"}
	p, err := Parse(arguments[:]...)
	if err != nil {
		t.Errorf("Parsing detects summary mode as an error")
	}
	if p.mode != summary {
		t.Errorf("Parsing does not set the mode to summary when the option is set")
	}
	if !p.IsSummaryMode() {
		t.Errorf("Parsing does set mode to summary but the method IsSummaryMode fails to report it.")
	}
}

func TestDetectingLinearMode(t *testing.T) {
	arguments := []string{"-m", "linear"}
	p, err := Parse(arguments[:]...)
	if err != nil {
		t.Errorf("Parsing detects linear mode as an error")
	}
	if p.mode != linear {
		t.Errorf("Parsing does not set the mode to linear when the option is set")
	}
}

func TestErrorParsing(t *testing.T) {
	arguments := []string{"-t", "15aaa"}
	_, err := Parse(arguments[:]...)
	if err == nil {
		t.Errorf("We do not detect when a time is not an integer.")
	}
}

func getSampleCsvAsStream() string {
	content := `
### Lesson 1
1_Question 1;1_Answer 1

### Lesson 2
2_Question 1;2_Answer 1
2_Question 2;2_Answer 2

### Lesson 3
3_Question 1;3_Answer 1
3_Question 2;3_Answer 2
3_Question 3;3_Answer 3
	`

	return content
}

func getTpp() TopicParsingParameters {
	return TopicParsingParameters{
		TopicAnnounce: "### Lesson ",
		QaSep:         ";",
	}
}

// Testing the way to get the data into the topic data structure.
func TestParseStream(t *testing.T) {

	r := strings.NewReader(getSampleCsvAsStream())
	p := getTpp()

	topic := ParseTopic(r, p)
	count := topic.GetSubsectionsCount()
	if count != 3 {
		t.Errorf("After parsing the stream should result in 3 subtopics. We have counted %d\n", count)
	}

	qa := topic.BuildQuestionsSet()
	count = qa.GetCount()
	if count != 6 {
		t.Errorf("We should have a list of 6 questions for the global but we found %d\n", count)
	}

	for i := 1; i <= 3; i++ {
		qa = topic.BuildQuestionsSet(strconv.Itoa(i))
		count = qa.GetCount()
		if count != i {
			t.Errorf("We should have a list of %d questions for the global but we found %d\n", i, count)
		}
	}

}
