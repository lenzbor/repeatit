package datamodel

import (
	"regexp"
	"testing"
)

var (
	emptyLine, _     = regexp.Compile("^\\s*$")
	loop, _          = regexp.Compile("^Loop\\s{1,}\\([0-9]{1,}/[0-9]{1,}\\)$")
	separator, _     = regexp.Compile("^-{1,}")
	nbOfQuestions, _ = regexp.Compile("^Nb of questions:\\s[0-9]{1,}")
	limitReached, _  = regexp.Compile("^Limit reached. Exiting. Number of loops set to:\\s[0-9]{1,}")
)

// TestAddEntry is testing not only the AddEntry function but the GetCount
// too and the initialization of a QA structure.
func TestAddEntry(t *testing.T) {
	qa := NewQA()
	if qa.GetCount() != 0 {
		t.Errorf("A freshly created QA structure should not contain any element. But the GetCount function reports %d\n", qa.GetCount())
	}
	qa.AddEntry("question-1", "answer-1")
	if qa.GetCount() != 1 {
		t.Errorf("Expected 1 entry but received %d\n", qa.GetCount())
	}
	qa.AddEntry("question-2", "answer-2")
	if qa.GetCount() != 2 {
		t.Errorf("Expected 2 entry but received %d\n", qa.GetCount())
	}
}

// TestConcatenate check that adding a set of questions/answers to another
// set is working fine.
func TestConcatenate(t *testing.T) {
	qa := NewQA()
	qa.AddEntry("question", "answer")

	otherQa := NewQA()
	otherQa.AddEntry("q1", "a1")
	otherQa.AddEntry("q2", "a2")
	qa.Concatenate(otherQa)

	count := qa.GetCount()
	if count != 3 {
		t.Errorf("Concatenate does not work find. We are expecting 3 but get %d\n", count)
	}
}

// TestNewTopic valides the construction of a topic.
func TestNewTopic(t *testing.T) {
	topic := NewTopic()
	if topic.vocabulary == nil {
		t.Errorf("A new topic should not have its vocabulary empty.")
	}
	if topic.sentences == nil {
		t.Errorf("A new topic should not have its sentences empty.")
	}
	count := topic.GetVocabularySubsectionsCount()
	if count != 0 {
		t.Errorf("Was expecting 0 but received a count of %d\n", count)
	}
	count = topic.GetSentencesSubsectionsCount()
	if count != 0 {
		t.Errorf("Was expecting 0 but received a count of %d\n", count)
	}
}
