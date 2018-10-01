package datamodel

import "fmt"

// Topic represents the list of subsections of the file with the questions
// attached for that section.
type Topic struct {
	list map[string]QuestionsAnswers
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
