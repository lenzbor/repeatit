package main

import (
	"fmt"
	"io"
	"strings"
	"time"
)

const (
	// Lesson is the string that is searched in the vocabulary file to
	// delimit the lessons. The lesson number should be right after the
	// delimiter, on the same line.
	Lesson = "### Lesson "
	// Sentences is the string that is searched in the vocabulary file to
	// delimit the sentences. The lesson number of the sentences should be
	// right after the delimiter, on the same line.
	Sentences = "### Sentences "
)

// QuestionsAnswers is a datastructure to store questions and their matching
// answers. The answers[i] matches questions[i].
type QuestionsAnswers struct {
	questions []string
	answers   []string
}

// Topic represents the list of subsections of the file with the questions
// attached for that section.
type Topic struct {
	list map[string]QuestionsAnswers
}

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

type interrogationMode int

const (
	linear  interrogationMode = iota // will ask questions in the same order as the file
	random                           // will ask questions in a random order
	summary                          // ask to show the list of subsections
)

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

// NewQA builds an empty set of questions/answers.
func NewQA() QuestionsAnswers {
	return QuestionsAnswers{}
}

// GetCount returns the number of entries for the questions.
func (qa QuestionsAnswers) GetCount() int {
	size := 0
	if qa.questions != nil {
		size = len(qa.questions)
	}
	return size
}

// NewTopic creates a new topic. A topic is a set of questions
// with a title.
func NewTopic() Topic {
	return Topic{
		list: make(map[string]QuestionsAnswers),
	}
}

// GetSubsection returns the current list of questions for a given topic id.
// If there is no associated questions and answers for this topic id, it
// returns a new structure.
func (topic *Topic) GetSubsection(id string) QuestionsAnswers {
	qa := topic.list[id]
	if qa.questions == nil {
		qa = NewQA()
		topic.list[id] = qa
	}
	return qa
}

// SetSubsection defines a subsection with a given id and associates
// to it a list of questions.
func (topic *Topic) SetSubsection(id string, qa QuestionsAnswers) {
	topic.list[id] = qa
}

// GetSubsectionsCount returns the number of subtopics.
func (topic Topic) GetSubsectionsCount() int {
	size := 0
	if topic.list != nil {
		size = len(topic.list)
	}
	return size
}

// GetSubsectionsName returns the list of subtopics that have been imported.
func (topic Topic) GetSubsectionsName() []string {
	subsections := []string{}
	if topic.GetSubsectionsCount() != 0 {
		subsections = make([]string, 0, len(topic.list))
		for id := range topic.list {
			subsections = append(subsections, id)
		}
	}
	return subsections
}

// AddEntry adds a set of question/answer to the already existing set.
func (qa *QuestionsAnswers) AddEntry(q string, a string) {
	qa.questions = append(qa.questions, q)
	qa.answers = append(qa.answers, a)
}

// Concatenate adds the entries of the parameter to an existing QA set.
func (qa *QuestionsAnswers) Concatenate(qaToAdd ...QuestionsAnswers) {
	var count int
	for _, toAdd := range qaToAdd {
		count = toAdd.GetCount()
		if count > 0 {
			qa.questions = append(qa.questions, toAdd.questions...)
			qa.answers = append(qa.answers, toAdd.answers...)
		}
	}
}

// BuildQuestionsSet creates a set of questions based on a Topic. We use a
// variadic list of parameters to allow to supply as many as topic on which
// the user wants to be questionned. If she/he supplies nothing, we use the
// the whole topic.
func (topic Topic) BuildQuestionsSet(ids ...string) QuestionsAnswers {
	qa := NewQA()
	var qaForID QuestionsAnswers
	var subsections = ids
	if len(subsections) == 0 {
		fmt.Println("     *** You supplied no subsection, we take them all ***")
		subsections = topic.GetSubsectionsName()
	}
	for _, ID := range subsections {
		qaForID = topic.GetSubsection(ID)
		qa.Concatenate(qaForID)
	}

	return qa
}
