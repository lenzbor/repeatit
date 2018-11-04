package engine

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/boris-lenzinger/repeatit/tools"

	"github.com/boris-lenzinger/repeatit/datamodel"
)

// AskQuestions will question the user on the set of questions. The
// parameter object will supply data to refine the questioning.
func AskQuestions(qa datamodel.QuestionsAnswers, p datamodel.InterrogationParameters) {
	loopsCount, i, idxQuestions := 0, 0, 0

	var wg sync.WaitGroup
	wg.Add(3)
	nbOfQuestions := qa.GetCount()

	if nbOfQuestions == 0 {
		tools.NegativeStatus("Number of questions is zero. Please check your file.")
		os.Exit(1)
	}

	// Handling channels in sub-goroutines
	go fanOutChannel(&wg, p.Qachan, p.Publisher)
	go publishChanToWriter(&wg, p.Publisher, p.GetOutputStream(), nbOfQuestions, p.GetLimit())
	go fanOutChannel(&wg, p.Command, p.Publisher)

	var question, answer string
	s := bufio.NewScanner(p.GetInputStream())
	for {
		if idxQuestions%nbOfQuestions == 0 {
			loopsCount++
			if loopsCount > p.GetLimit() {
				// if the qa chan is closed, then we have to close the others.
				close(p.Qachan)
				close(p.Command)
				break
			}
		}
		if p.IsRandomMode() {
			i = int(rand.Int31n(int32(nbOfQuestions)))
		}
		question = qa.GetQuestion(i)
		answer = qa.GetAnswer(i)
		if p.IsReversedMode() {
			// user has requested Jeopardy like
			question = qa.GetAnswer(i)
			answer = qa.GetQuestion(i)
		}
		tools.Debugf("Pushing question %q to qachan", question)
		p.Qachan <- fmt.Sprintf("%s", question)
		if p.IsInteractive() {
			if s.Scan() {
				p.Command <- s.Text()
			}
		} else {
			time.Sleep(p.GetPauseTime())
		}
		p.Qachan <- fmt.Sprintf("%s", answer)

		if !p.IsRandomMode() {
			i = (i + 1) % nbOfQuestions
		}
		idxQuestions++
	}

	wg.Wait()
}
