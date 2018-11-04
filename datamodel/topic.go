package datamodel

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/boris-lenzinger/repeatit/tools"
)

// Topic represents the list of subsections of the file with the questions
// attached for that section. Usually, a topic will be a lesson subdivided
// in vocabulary, grammar, sentences, etc.
type Topic struct {
	// The language for the original words
	LearnedLanguage string `json:"learned"`
	// The language for the translation
	NativeLanguage string `json:"native"`
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

// ShowSummary displays on user what is available  in this topic.
func (topic Topic) ShowSummary() {
	tools.WriteInCyan("  Content of the loaded resources\n")
	fmt.Printf("    * Learned: %s\n", topic.LearnedLanguage)
	fmt.Printf("    * Native: %s\n", topic.NativeLanguage)
	fmt.Printf("      - Lessons available: %s", topic.computeLessonsRange())
}

// computeLessonRange returns an easy string representation for the lessons stored
// in a topic. Instead of displaying 1, 2, 3, 4 for instance, it will return 1:4.
// For 1,2,3,4,5,8,9,10 it will return 1:5,8:10
func (topic Topic) computeLessonsRange() string {
	lessonsID := make([]int, len(topic.vocabulary))
	i := 0
	var err error
	for ID := range topic.vocabulary {
		lessonsID[i], err = strconv.Atoi(ID)
		if err != nil {
			tools.Warningf("a lesson is referenced with %q which is not an integer", ID)
		}
		i++
	}
	sort.Ints(lessonsID)
	return computeRangesOnArrayOfInts(lessonsID)
}

func computeRangesOnArrayOfInts(v []int) string {
	rangeOfIDs := ""
	var startOfSuite, latestContiguous int
	// Then we ignore the zeroes
	for i := 0; i < len(v); i++ {
		if v[i] <= 0 {
			continue
		}
		if rangeOfIDs == "" {
			startOfSuite = v[i]
			latestContiguous = startOfSuite
			rangeOfIDs = strconv.Itoa(v[i])
			continue
		}
		if v[i-1] == v[i] {
			// duplicate in lessons number. Ignore
			continue
		}
		// there is a first element. Let's check now if the current element is
		// contiguous with the previous one...
		if v[i-1]+1 == v[i] {
			latestContiguous = v[i]
			// Is contiguous. Check if next is still...
			if i != len(v)-1 {
				continue
			}
			// end is reached
			rangeOfIDs += ":" + strconv.Itoa(latestContiguous)
			break
		}

		// not contiguous
		if startOfSuite == latestContiguous {
			rangeOfIDs += "," + strconv.Itoa(v[i])
		} else {
			rangeOfIDs += ":" + strconv.Itoa(latestContiguous) + "," + strconv.Itoa(v[i])
		}
		startOfSuite = v[i]
		latestContiguous = startOfSuite
	}
	return rangeOfIDs
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
