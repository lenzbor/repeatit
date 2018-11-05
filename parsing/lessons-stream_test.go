package parsing

import (
	"strconv"
	"strings"
	"testing"

	"github.com/boris-lenzinger/repeatit/internal/tests"
)

// Testing the way to get the data into the topic data structure.
func TestParseStream(t *testing.T) {

	r := strings.NewReader(tests.GetSampleCsvAsStream())
	p := tests.GetTpp()

	topic, err := ParseTopic(r, p)
	if err != nil {
		t.Fatalf("parsing of topic should not raise an error. Get %v", err)
	}
	count := topic.GetVocabularySubsectionsCount()
	if count != 3 {
		t.Errorf("After parsing the stream should result in 3 subtopics. We have counted %d\n", count)
	}

	qa := topic.BuildVocabularyQuestionsSet()
	count = qa.GetCount()
	if count != 6 {
		t.Errorf("We should have a list of 6 questions for the global but we found %d\n", count)
	}

	for i := 1; i <= 3; i++ {
		qa = topic.BuildVocabularyQuestionsSet(strconv.Itoa(i))
		count = qa.GetCount()
		if count != i {
			t.Errorf("We should have a list of %d questions for the global but we found %d\n", i, count)
		}
	}

}
