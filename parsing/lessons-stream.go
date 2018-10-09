package parsing

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/boris-lenzinger/repeatit/datamodel"
	"github.com/boris-lenzinger/repeatit/tools"
)

// ParseTopic is reading the data source and transforms it to a topic
// structure.
func ParseTopic(r io.Reader, p datamodel.TopicParsingParameters) datamodel.Topic {
	// Reading the file line by line
	s := bufio.NewScanner(r)

	lines := []string{}
	for s.Scan() {
		lines = append(lines, s.Text())
	}

	topic := datamodel.NewTopic()
	var subsectionID string
	qaSubsection := datamodel.NewQA()
	var isVocabularySection, isSentencesSection bool
	for i := 0; i < len(lines); i++ {
		input := lines[i]
		// Ignore empty lines
		if len(input) > 0 {
			split := strings.Split(input, p.QaSep)
			switch len(split) {
			case 1:
				if strings.HasPrefix(input, p.LessonAnnounce) {
					tools.Debug(fmt.Sprintf("Found vocabulary delimiter: %s", input))
					subsectionID = strings.TrimPrefix(input, p.LessonAnnounce)
					qaSubsection = topic.GetVocabularySubsection(subsectionID)
					isVocabularySection = true
					isSentencesSection = false
				} else if strings.HasPrefix(input, p.SentenceAnnounce) {
					tools.Debug(fmt.Sprintf("Found sentences delimiter: %s", input))
					subsectionID = strings.TrimPrefix(input, p.SentenceAnnounce)
					qaSubsection = topic.GetSentencesSubsection(subsectionID)
					isVocabularySection = false
					isSentencesSection = true
				}
			default:
				// Question is in split[0] while answer in in split[1]. It may happen
				// the answer contains the separator so we have to join the different
				// elements.
				tools.Debug(fmt.Sprintf("Adding entry %s", split[0]))
				qaSubsection.AddEntry(split[0], strings.Join(split[1:], p.QaSep))
				if isVocabularySection {
					topic.SetVocabularySubsection(subsectionID, qaSubsection)
					topic.IncreaseVocabularyCount()
				} else if isSentencesSection {
					topic.SetSentencesSubsection(subsectionID, qaSubsection)
					topic.IncreaseSentencesCount()
				}
				tools.Debug(fmt.Sprintf("Number of vocabulary sections: %d", topic.GetVocabularySubsectionsCount()))
				tools.Debug(fmt.Sprintf("Total number of vocabulary words: %d", topic.GetNumberOfWords()))
				tools.Debug(fmt.Sprintf("Number of sentences sections: %d", topic.GetVocabularySubsectionsCount()))
				tools.Debug(fmt.Sprintf("Total number of sentences: %d", topic.GetNumberOfSentences()))
			}
		}
	}
	return topic
}
