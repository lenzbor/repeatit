package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// AskQuestions will question the user on the set of questions. The
// parameter object will supply data to refine the questioning.
func AskQuestions(qa QuestionsAnswers, p InterrogationParameters) {
	fullLoop, i, j := 0, 0, 0

	var wg sync.WaitGroup
	wg.Add(3)
	nbOfQuestions := qa.GetCount()

	go fanOutChannel(&wg, p.qachan, p.publisher)
	go publishChanToWriter(&wg, p.publisher, p.GetOutputStream(), nbOfQuestions, p.limit)
	go fanOutChannel(&wg, p.command, p.publisher)

	var question, answer string
	s := bufio.NewScanner(p.in)
	for {
		if j%nbOfQuestions == 0 {
			fullLoop++
			if fullLoop > p.limit {
				// if the qa chan is closed, then we have to close the others.
				close(p.qachan)
				close(p.command)
				break
			}
		}
		if p.mode == random {
			i = int(rand.Int31n(int32(nbOfQuestions)))
		}
		question = qa.questions[i]
		answer = qa.answers[i]
		if p.IsReversedMode() {
			question = qa.answers[i]
			answer = qa.questions[i]
		}
		p.qachan <- fmt.Sprintf("%s", question)
		if p.interactive {
			if s.Scan() {
				p.command <- s.Text()
			}
		} else {
			time.Sleep(p.wait)
		}
		p.qachan <- fmt.Sprintf("%s", answer)

		if p.mode == linear {
			i = (i + 1) % nbOfQuestions
		}
		j++
	}

	wg.Wait()
}
