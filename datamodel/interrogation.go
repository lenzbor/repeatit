package datamodel

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// InterrogationMode makes it possible to change the way to be questionned.
type InterrogationMode int

const (

	// Linear configures the engine to ask questions in the same order as they
	// appear in the file
	Linear InterrogationMode = iota
	// Random configures the engine to ask questions in a random order
	Random
	// Summary requires to show the list of subsections so the user has
	// a clear view of what the available content is.
	Summary

	// DefaultInterrogationPause is the default pause time between 2
	// questions. Default value is 2 seconds.
	DefaultInterrogationPause = 2 * time.Second

	// DefaultLoopCount is the default number of loops during an interrogation.
	DefaultLoopCount = 10
)

// InterrogationParameters is a datastructure that contains the parameters required
// by the user for the questionning.
type InterrogationParameters struct {
	interactive bool
	wait        time.Duration
	mode        InterrogationMode
	// Default is to use io.Stdin. Allows to send command to the engine
	in io.Reader
	// The place where the questions are written to
	out io.Writer
	// the list of selected subsections chosen for the questioning
	subsections string
	// Limit is the number of times the list is repeated during interrogation.
	// Default is defined in the constant DefaultLoopCount. This can be overriden
	// on the command line or in the configuration file.
	limit int
	// Requires that questions becomes answers and answers becomes questions
	reversed bool
	// Experimental. Channel to receive questions and answers
	Qachan chan string
	// Experimental. Channel to receive commands
	Command chan string
	// Experimental. Channel to publish to the output. This channel collects all that needs to be displayed to the user.
	Publisher chan string
	// Absolute path to the lesson file to use
	lessonsFile string
}

// NewInterrogationParameters creates a default instance of the
// parameters that will be used to do the interrogation.
// Here are the values returned:
//   * it is not interactive
//   * wait time between 2 questions is 2 seconds
//   * the reader is io.Stdin . This allows to send commands to the engine
//   * the output is io.Stdout
//   * the number of loops is 10
//   * the interrogation is not in Jeopardy mode
func NewInterrogationParameters() InterrogationParameters {
	loopCount := DefaultLoopCount
	configuredLoopCount := viper.GetInt("limit")
	if configuredLoopCount != 0 {
		loopCount = configuredLoopCount
	}
	return InterrogationParameters{
		interactive: false,
		wait:        DefaultInterrogationPause,
		mode:        Random,
		in:          os.Stdin,
		out:         os.Stdout,
		limit:       loopCount,
		reversed:    false,
		lessonsFile: "NoFileDefined",
		Qachan:      make(chan string),
		Publisher:   make(chan string),
		Command:     make(chan string),
	}
}

// IsInteractive tells if the interrogation requires a human interaction
// or not.
func (p *InterrogationParameters) IsInteractive() bool {
	return p.interactive
}

// SetInteractive changes the parameters for questions to interactive so the
// user has to hit Return key to pass to next question.
func (p *InterrogationParameters) SetInteractive() {
	p.interactive = true
}

// IsSummaryMode tells if the parameters require to have a summary of the subsections.
func (p *InterrogationParameters) IsSummaryMode() bool {
	return p.mode == Summary
}

// IsRandomMode tells if the questioning is done randomly or linearly
func (p *InterrogationParameters) IsRandomMode() bool {
	return p.mode == Random
}

// GetInputStream gets the Reader from where we read the user input.
func (p *InterrogationParameters) GetInputStream() io.Reader {
	return p.in
}

// GetOutputStream gets the Writer where questions will be written to.
func (p *InterrogationParameters) GetOutputStream() io.Writer {
	return p.out
}

// IsReversedMode tells if the user wants that the left column are now answers and right column(s) are the questions
func (p *InterrogationParameters) IsReversedMode() bool {
	return p.reversed
}

// GetListOfSubsections returns a string array containing all the subsections selected by
// the end user.
func (p *InterrogationParameters) GetListOfSubsections() []string {
	if len(p.subsections) == 0 {
		return nil
	}
	return strings.Split(p.subsections, ",")
}

// GetLimit returns the number of loops for lessons to learn.
func (p *InterrogationParameters) GetLimit() int {
	return p.limit
}

// GetLessonsFile returns the absolute path to the lessons file that is used
// to ask questions to the user.
func (p *InterrogationParameters) GetLessonsFile() string {
	return p.lessonsFile
}

// SetLessonsFile assigns the path to lessons file to parse and use as
// the source of questions/answers.
func (p *InterrogationParameters) SetLessonsFile(path string) {
	p.lessonsFile = path
}

// GetPauseTime returns the pause between each question in milliseconds.
func (p *InterrogationParameters) GetPauseTime() time.Duration {
	return p.wait
}
