package engine

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/boris-lenzinger/repeatit/parsing"

	"github.com/boris-lenzinger/repeatit/datamodel"
	"github.com/boris-lenzinger/repeatit/tools"
)

// StartEngine is starting the command interpretor.
func StartEngine(t datamodel.Topic) {
	t.ShowSummary()
	rand.Seed(time.Now().UTC().UnixNano())

loop:
	for {
		tools.WriteInCyan(fmt.Sprintf("> "))
		userInput := tools.ReadFromStdin(true)
		switch {
		case userInput == "":
			continue
		case strings.HasPrefix(userInput, "select"):
			selected := strings.TrimPrefix(userInput, "select")
			selected = strings.TrimPrefix(selected, " ")
			selectedLessons, err := parsing.ParseNumberSerie(selected)
			if err != nil {
				tools.NegativeStatus(fmt.Sprintf("error while parsing the list of lessons: %v", err))
				continue
			}
			numberOfDigitsForLessonsIDs := len(strconv.Itoa(t.GetVocabularySubsectionsCount()))
			qa := t.BuildVocabularyQuestionsSet(tools.ConvertToStringsArray(selectedLessons, numberOfDigitsForLessonsIDs)...)
			interrogParams := datamodel.NewInterrogationParameters()
			AskQuestions(qa, interrogParams)
			fmt.Printf("Session is over...\n")
		case userInput == "list":
			fmt.Printf("Lessons available: %s\n", t.ComputeLessonsRange())
		case userInput == "help":
		case userInput == "quit":
			fmt.Println("Exiting on user request.")
			break loop
		default:
			tools.NegativeStatus(fmt.Sprintf("%q is an invalid command. Use help to get the full list of supported commands", userInput))
		}
	}
	os.Exit(1)
}
