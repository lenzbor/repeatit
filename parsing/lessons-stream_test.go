package parsing

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/boris-lenzinger/repeatit/datamodel"
)

func getSampleCsvAsStream() string {
	content := fmt.Sprintf(`
%s 1
1_Question 1;1_Answer 1

%s 2
2_Question 1;2_Answer 1
2_Question 2;2_Answer 2

%s 3
3_Question 1;3_Answer 1
3_Question 2;3_Answer 2
3_Question 3;3_Answer 3
	`, datamodel.Lesson, datamodel.Lesson, datamodel.Lesson)

	return content
}

func getTpp() datamodel.TopicParsingParameters {
	return datamodel.TopicParsingParameters{
		TopicAnnounce: fmt.Sprintf("%s ", datamodel.Lesson),
		QaSep:         ";",
	}
}

// Testing the way to get the data into the topic data structure.
func TestParseStream(t *testing.T) {

	r := strings.NewReader(getSampleCsvAsStream())
	p := getTpp()

	topic := ParseTopic(r, p)
	count := topic.GetSubsectionsCount()
	if count != 3 {
		t.Errorf("After parsing the stream should result in 3 subtopics. We have counted %d\n", count)
	}

	qa := topic.BuildQuestionsSet()
	count = qa.GetCount()
	if count != 6 {
		t.Errorf("We should have a list of 6 questions for the global but we found %d\n", count)
	}

	for i := 1; i <= 3; i++ {
		qa = topic.BuildQuestionsSet(strconv.Itoa(i))
		count = qa.GetCount()
		if count != i {
			t.Errorf("We should have a list of %d questions for the global but we found %d\n", i, count)
		}
	}

}
