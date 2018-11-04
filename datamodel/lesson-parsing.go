package datamodel

const (
	// LessonDelimiter is the string that is searched in the vocabulary file to
	// delimit the lessons. The lesson number should be right after the
	// delimiter, on the same line.
	LessonDelimiter = "### Lesson "

	// DefaultQaSep is the default character used to separate questions
	// from answers in the CSV file
	DefaultQaSep = ";"

	// Sentences is the string that is searched in the vocabulary file to
	// delimit the sentences. The lesson number of the sentences should be
	// right after the delimiter, on the same line.
	Sentences = "### Sentences Lesson"
)

// TopicParsingParameters is a data structure that helps to parse the lines that
// define the different sections.
type TopicParsingParameters struct {
	// LessonAnnounce is the string that is used to announce the lesson
	// section in the csv file.
	// The text after this string will be considered as the ID of the topic.
	LessonAnnounce string
	// SentenceAnnounce is the string that is used to announce the Sentence
	// section in the csv file.
	// The text after this string will be considered as the ID of the topic.
	SentenceAnnounce string
	// QaSep is the separator on the line between the question and the answer in
	// the csv file. If this separator is found multiple times on the line, the
	// first one is considered as the separator.
	QaSep string
}

// NewTopicParsingParameters returns a set of values that makes it possible
// to parse the CSV file containing the questions/answers.
func NewTopicParsingParameters() TopicParsingParameters {
	return TopicParsingParameters{
		LessonAnnounce:   LessonDelimiter,
		SentenceAnnounce: Sentences,
		QaSep:            DefaultQaSep,
	}
}
