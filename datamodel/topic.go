package datamodel

import (
	"fmt"

	"github.com/boris-lenzinger/repeatit/tools"
)

// Topic represents the list of subsections of the file with the questions
// attached for that section.
type Topic struct {
	// the map listing the vocabulary of the lessons
	// (by number or name of lesson)
	vocabulary map[string]QuestionsAnswers
	// the map listing the sentences by number of lessons
	// or lessons names.
	sentences       map[string]QuestionsAnswers
	vocabularyCount int
	sentencesCount  int
}

// NewTopic creates a new object with initialized fields. A topic is a set
// of questions with a title. The topic is created with an initialized empty
// map of questions/answers
func NewTopic() Topic {
	return Topic{
		vocabulary: make(map[string]QuestionsAnswers),
		sentences:  make(map[string]QuestionsAnswers),
	}
}

// GetVocabularySubsection returns the current list of vocabulary questions
// for a given topic id.
// If there is no associated questions and answers for this topic id, it
// returns a new structure.
func (topic *Topic) GetVocabularySubsection(ID string) QuestionsAnswers {
	qa, ok := topic.vocabulary[ID]
	if !ok {
		qa = NewQA()
		topic.vocabulary[ID] = qa
	}
	return qa
}

// GetSentencesSubsection returns the current list of questions for a
// given topic id.
// If there is no associated questions and answers for this topic id, it
// returns a new structure.
func (topic *Topic) GetSentencesSubsection(ID string) QuestionsAnswers {
	qa, ok := topic.sentences[ID]
	if !ok {
		qa = NewQA()
		topic.sentences[ID] = qa
	}
	return qa
}

// IncreaseSentencesCount increments the number of sentences words.
func (topic *Topic) IncreaseSentencesCount() {
	topic.sentencesCount++
}

// IncreaseVocabularyCount increments the number of vocabulary words.
func (topic *Topic) IncreaseVocabularyCount() {
	topic.vocabularyCount++
}

// SetVocabularySubsection defines (or overrides if it already existed) a subsection
// with a given ID and associates to it a list of questions.
func (topic *Topic) SetVocabularySubsection(ID string, qa QuestionsAnswers) {
	topic.vocabulary[ID] = qa
}

// SetSentencesSubsection defines (or overrides if it already existed) a subsection
// with a given ID and associates to it a list of questions.
func (topic *Topic) SetSentencesSubsection(ID string, qa QuestionsAnswers) {
	topic.sentences[ID] = qa
}

// GetVocabularySubsectionsCount returns the number of vocabulary
// lessons subtopics.
func (topic Topic) GetVocabularySubsectionsCount() int {
	return len(topic.vocabulary)
}

// GetSentencesSubsectionsCount returns the number of sentences subtopics.
func (topic Topic) GetSentencesSubsectionsCount() int {
	return len(topic.sentences)
}

// GetVocabularySubsectionsName returns the list of vocabulary lessons
// subtopics that have been imported.
func (topic Topic) GetVocabularySubsectionsName() []string {
	subsections := []string{}
	if topic.GetVocabularySubsectionsCount() != 0 {
		subsections = make([]string, 0, topic.GetVocabularySubsectionsCount())
		for ID := range topic.vocabulary {
			subsections = append(subsections, ID)
		}
	}
	return subsections
}

// GetSentencesSubsectionsName returns the list of sentences subtopics that
// have been imported.
func (topic Topic) GetSentencesSubsectionsName() []string {
	subsections := []string{}
	if topic.GetSentencesSubsectionsCount() != 0 {
		subsections = make([]string, 0, topic.GetSentencesSubsectionsCount())
		for ID := range topic.sentences {
			subsections = append(subsections, ID)
		}
	}
	return subsections
}

// BuildVocabularyQuestionsSet creates a set of questions based on a Topic.
// We use a
// variadic list of parameters to allow to supply as many as topic on which
// the user wants to be questionned. If she/he supplies nothing, we use the
// the whole topic.
// BuildVocabularyQuestionsSet creates a set of questions based on a Topic.
func (topic Topic) BuildVocabularyQuestionsSet(ids ...string) QuestionsAnswers {
	qa := NewQA()
	var qaForID QuestionsAnswers
	var subsections = ids
	if len(subsections) == 0 {
		fmt.Println("     *** You supplied no subsection, we take them all ***")
		subsections = topic.GetVocabularySubsectionsName()
	}
	for _, ID := range subsections {
		tools.Debug(fmt.Sprintf("Getting vocabulary from section %s", ID))
		qaForID = topic.GetVocabularySubsection(ID)
		tools.Debug(fmt.Sprintf("Found %d entries in the QA section", qaForID.GetCount()))
		qa.Concatenate(qaForID)
	}

	return qa
}

// BuildSentencesQuestionsSet creates a set of questions based on a Topic.
// We use a
// variadic list of parameters to allow to supply as many as topic on which
// the user wants to be questionned. If she/he supplies nothing, we use the
// the whole topic.
func (topic Topic) BuildSentencesQuestionsSet(ids ...string) QuestionsAnswers {
	qa := NewQA()
	var qaForID QuestionsAnswers
	var subsections = ids
	if len(subsections) == 0 {
		fmt.Println("     *** You supplied no subsection, we take them all ***")
		subsections = topic.GetSentencesSubsectionsName()
	}
	for _, ID := range subsections {
		qaForID = topic.GetSentencesSubsection(ID)
		qa.Concatenate(qaForID)
	}

	return qa
}

// GetNumberOfWords returns the number of words stored in the structure
func (topic *Topic) GetNumberOfWords() int {
	return topic.vocabularyCount
}

// GetNumberOfSentences returns the number of sentences stored in the structure
func (topic *Topic) GetNumberOfSentences() int {
	return topic.sentencesCount
}

// String makes a string representation from the object for debug
// purpose.
func (topic *Topic) String() string {
	output := ""
	output += fmt.Sprintf("Number of words in this topic:     %d\n", topic.GetNumberOfWords())
	output += fmt.Sprintf("Number of sentences in this topic: %d\n", topic.GetNumberOfSentences())
	output += fmt.Sprintf("Details abouts words:\n")
	for key := range topic.vocabulary {
		output += fmt.Sprintf("\t* %s has %d words\n", key, topic.GetVocabularySubsection(key).GetCount())
	}
	output += fmt.Sprintf("Details abouts sentences:\n")
	for key := range topic.sentences {
		output += fmt.Sprintf("\t* %s has %d sentences\n", key, topic.GetSentencesSubsection(key).GetCount())
	}

	return output
}
