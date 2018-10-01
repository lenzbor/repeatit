package datamodel

import (
	"io"
	"strings"
	"time"
)

type interrogationMode int

const (
	// Lesson is the string that is searched in the vocabulary file to
	// delimit the lessons. The lesson number should be right after the
	// delimiter, on the same line.
	Lesson = "### Lesson "
	// Sentences is the string that is searched in the vocabulary file to
	// delimit the sentences. The lesson number of the sentences should be
	// right after the delimiter, on the same line.
	Sentences = "### Sentences "

	linear  interrogationMode = iota // will ask questions in the same order as the file
	random                           // will ask questions in a random order
	summary                          // ask to show the list of subsections
)

// TopicParsingParameters is a data structure that helps to parse the lines that
// split the different sections.
type TopicParsingParameters struct {
	// topicAnnounce is the string that is used to announce the section in the
	// csv file. For instance, '### Lesson '
	// The text after this string will be considered as the ID of the topic.
	TopicAnnounce string
	// QaSep is the separator on the line between the question and the answer in
	// the csv file. If this separator is found multiple times on the line, the
	// first one is considered as the separator.
	QaSep string
}

// InterrogationParameters is a datastructure that contains the parameters required
// by the user for the questionning.
type InterrogationParameters struct {
	interactive bool
	wait        time.Duration     // Default is to wait 2 seconds
	mode        interrogationMode // Default is random.
	in          io.Reader         // Default is to use io.Stdin. Allows to send command to the engine
	out         io.Writer         // The place where the questions are written to
	subsections string            // the list of selected subsections chosen for the questioning
	limit       int               // Limit is the number of times the list is repeated during interrogation. Default is 10
	reversed    bool              // Requires that questions becomes answers and answers becomes questions
	qachan      chan string       // Experimental. Channel to receive questions and answers
	command     chan string       // Experimental. Channel to receive commands
	publisher   chan string       // Experimental. Channel to publish to the output. This channel collects all that needs to be put to the user.
}

// IsSummaryMode tells if the parameters require to have a summary of the subsections.
func (p InterrogationParameters) IsSummaryMode() bool {
	return p.mode == summary
}

// GetOutputStream gets the Writer where questions will be written to.
func (p InterrogationParameters) GetOutputStream() io.Writer {
	return p.out
}

// IsReversedMode tells if the user wants that the left column are now answers and right column(s) are the questions
func (p InterrogationParameters) IsReversedMode() bool {
	return p.reversed
}

// GetListOfSubsections returns a string array containing all the subsections selected by
// the end user.
func (p InterrogationParameters) GetListOfSubsections() []string {
	if len(p.subsections) == 0 {
		return nil
	}
	return strings.Split(p.subsections, ",")
}
