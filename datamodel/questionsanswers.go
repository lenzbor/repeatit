package datamodel

// QuestionsAnswers is a datastructure to store questions and their matching
// answers. The answers[i] matches questions[i].
type QuestionsAnswers struct {
	questions []string
	answers   []string
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
