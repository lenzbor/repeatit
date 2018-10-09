package engine

import (
	"fmt"
	"io"
	"sync"

	"github.com/fatih/color"
)

// fanOutChannel reads from the readFrom channel and dispatch the elements
// to the writeTo channel. When reading from the readFrom channel breaks,
// we write to the stopper channel the name of the channel from which we
// cannot read anymore.
func fanOutChannel(wg *sync.WaitGroup, readFrom <-chan string, writeTo chan<- string) {
	defer wg.Done()

	for {
		select {
		case v, ok := <-readFrom:
			if !ok {
				return
			}
			if len(v) != 0 {
				writeTo <- v
			}
		}
	}
}

//
func publishChanToWriter(wg *sync.WaitGroup, readFrom <-chan string, out io.Writer, qCount int, maxLoops int) {
	defer wg.Done()
	itemsRead := 0
	currentLoop := 0
	c := color.New(color.FgBlue).Add(color.Bold)

	fmt.Fprintf(out, "Nb of questions: %d\n", qCount)

	for {
		if itemsRead%(2*qCount) == 0 {
			currentLoop++
			if currentLoop > maxLoops {
				fmt.Fprintf(out, "Limit reached. Exiting. Number of loops set to: %d\n", maxLoops)
				return
			}
			fmt.Fprintf(out, c.Sprintf("Loop (%d/%d)\n", currentLoop, maxLoops))
		}
		select {
		case v, ok := <-readFrom:
			if !ok {
				return
			}
			itemsRead++
			switch {
			case itemsRead%2 == 1:
				fmt.Fprintf(out, v)
				// Questions asked. Must publish the answer now.
			case itemsRead%2 == 0:
				fmt.Fprintf(out, "     --> "+v+"\n")
				fmt.Fprintf(out, "---------------------------\n")
			}
		}
	}
}
