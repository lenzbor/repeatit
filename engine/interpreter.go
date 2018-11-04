package engine

import (
	"fmt"
	"os"

	"github.com/boris-lenzinger/repeatit/datamodel"
	"github.com/boris-lenzinger/repeatit/tools"
)

// StartEngine is starting the command interpretor.
func StartEngine(t datamodel.Topic) {
	t.ShowSummary()
	for {
		tools.WriteInCyan(fmt.Sprintf("> "))
		userInput := tools.ReadFromStdin(false)
		switch userInput {
		case "list":
		case "help":
		case "quit":
			fmt.Println("Exiting on user request.")
			os.Exit(1)
		default:
			tools.NegativeStatus(fmt.Sprintf("%q is an invalid command. Use help to get the full list of supported commands", userInput))
		}
	}
}
