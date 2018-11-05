package datamodel

import (
	"encoding/json"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
)

// Language is the modelization of a structure that stores data to learn
// a language.
// The structure was built based on a book so it may not be relevant for
// other learning material.
type Language struct {
	// Meta references the metadata related to this language resource
	Meta Metadata `json:"metadata"`
	// Content contains the content extracted from the resource
	Content LanguageResources `json:"content"`
}

// GetLessonsCount returns the number of lessons for this resource
func (l Language) GetLessonsCount() int {
	return l.Content.GetLessonsCount()
}

// LoadLessons reads from a stream the json structure that represents
// the resources of the language to learn.
// If the stream reports an error in the parsing or if reading the stream
// fails, an error is reported.
func LoadLessons(r io.Reader) (Language, error) {
	ba, err := ioutil.ReadAll(r)
	if err != nil {
		return Language{}, errors.Wrap(err, "failed to read the language resource from the json streamed representation")
	}
	output := Language{}
	err = json.Unmarshal(ba, &output)
	if err != nil {
		return Language{}, errors.Wrap(err, "failed to transform the data source to a data structure")
	}
	return output, nil
}

// LanguageResources describes the different elements useful to learn a
// language (words, sentences, grammar, ...). Those elements are grouped
// into  lessons datastructure.
type LanguageResources struct {
	// Lessons is the list of the available lessons in this resource
	Lessons []Lesson `json:"lessons"`
}

// GetLessonsCount returns the number of lessons in the language resources.
func (lr *LanguageResources) GetLessonsCount() int {
	return len(lr.Lessons)
}

// NewLanguageResources initializes an empty language resource.
func NewLanguageResources() LanguageResources {
	return LanguageResources{
		Lessons: []Lesson{},
	}
}

// CreateNewLesson creates a new lesson in the list of the lessons.
// The method returns the ID of the lesson so one can enrich the lesson.
func (lr *LanguageResources) CreateNewLesson() int {
	ID := lr.getIDForNewLesson()
	l := NewLesson(ID)
	lr.Lessons = append(lr.Lessons, l)
	return ID
}

// getIDForNewLesson returns an ID for a new lesson. The lesson is not created
// at that step.
func (lr *LanguageResources) getIDForNewLesson() int {
	ID := 0
	if len(lr.Lessons) > 0 {
		for i := 0; i < len(lr.Lessons); i++ {
			if lr.Lessons[i].ID > ID {
				ID = lr.Lessons[i].ID
			}
		}
	}
	return ID + 1
}

// Lesson describes the content of a single lesson.
type Lesson struct {
	// ID is the unique identifier of the lesson
	ID int `json:"id"`
	// Title is the title of the lesson so one can check what it is about
	Title Resource `json:"title"`
	// Still need to add the vocabulary, the grammary and the sentences
	Vocabulary []Resource `json:"vocabulary"`
	Sentences  []Resource `json:"sentences"`
}

// NewLesson is the default constructor for a lesson. All fields, except ID,
// are set to empty values.
func NewLesson(ID int) Lesson {
	return Lesson{
		ID:         ID,
		Title:      Resource{},
		Vocabulary: []Resource{},
		Sentences:  []Resource{},
	}
}

// Resource contains a string in the language to learn and its translation
// in the native language (native and language to learn are based on the
// Metadata structure.
type Resource struct {
	// Learning is the resource in the language you want to learn
	Learning string `json:"learn"`
	// Native is the resource in the language in your language
	Native string `json:"native"`
}

// Metadata is the data that describes the learning material.
type Metadata struct {
	// Learning is the language that is being learnt
	Learning string `json:"learning"`
	// Native is the native language of the learner
	Native string `json:"native"`
	// Book is the book from which the content was extracted
	Book Book `json:"book"`
}

// Book is a simplist description of a book that is used as a
// shortcut to have a datastructure in the metadata part.
type Book struct {
	// Title is the title of the book
	Title string `json:"title"`
	// Authors are the authors of the book
	Authors []Author `json:"authors"`
	// ISBN is the universal ID of the book that makes it possible
	// to search for it among books databases.
	ISBN string `json:"isbn"`
}

// Author is a datastructure to fill data about a book
type Author struct {
	// Firstname of the author
	Firstname string `json:"firstname"`
	// Lastname of the author
	Lastname string `json:"lastname"`
}
