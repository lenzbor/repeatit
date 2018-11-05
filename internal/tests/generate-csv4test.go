package tests

import (
	"fmt"

	"github.com/boris-lenzinger/repeatit/datamodel"
)

// GetSampleCsvAsStream returns a string that can be used to make tests on the
// parsing of the CSVs.
func GetSampleCsvAsStream() string {
	content := fmt.Sprintf(`#native;learnt

%s 1
1_Question 1;1_Answer 1

%s 2
2_Question 1;2_Answer 1
2_Question 2;2_Answer 2

%s 3
3_Question 1;3_Answer 1
3_Question 2;3_Answer 2
3_Question 3;3_Answer 3
	`, datamodel.LessonDelimiter, datamodel.LessonDelimiter, datamodel.LessonDelimiter)

	return content
}

// GetTpp returns topic parsing parameters for testing purpose.
func GetTpp() datamodel.TopicParsingParameters {
	return datamodel.TopicParsingParameters{
		LessonAnnounce:   datamodel.LessonDelimiter,
		SentenceAnnounce: datamodel.SentencesDelimiter,
		QaSep:            datamodel.DefaultQaSep,
	}
}
