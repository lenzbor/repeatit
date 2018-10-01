package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

// Parse is parsing a list of strings to build a set of parameters
// for the AskQuestion function.
func Parse(args ...string) (InterrogationParameters, error) {
	p := InterrogationParameters{
		interactive: false,
		wait:        2 * time.Second,
		mode:        random,
		in:          os.Stdin,
		out:         os.Stdout,
		subsections: "",
		limit:       1,
		qachan:      make(chan string),
		command:     make(chan string),
		publisher:   make(chan string),
	}
	for i, opt := range args {
		switch opt {
		case "-i":
			p.interactive = true
		case "-t":
			value, err := strconv.Atoi(args[i+1])
			if err != nil {
				return p, fmt.Errorf("the time you set (%s) is not an integer. Please set the time in milliseconds", args[i+1])
			}
			p.wait = time.Duration(value) * time.Millisecond
		case "-m":
			// The other mode is the default so we have nothing to do.
			if args[i+1] == "linear" {
				p.mode = linear
			}
		case "-s":
			p.mode = summary
		case "-l":
			p.subsections = args[i+1]
		case "-r":
			p.reversed = true
		}
	}
	return p, nil
}

// ParseTopic is reading the data source and transforms it to a topic
// structure.
func ParseTopic(r io.Reader, p TopicParsingParameters) Topic {
	// Reading the file line by line
	s := bufio.NewScanner(r)

	lines := make([]string, 50)
	for s.Scan() {
		lines = append(lines, s.Text())
	}

	topic := NewTopic()
	var subsectionID string
	qaSubsection := NewQA()
	for i := 0; i < len(lines); i++ {
		input := lines[i]
		// Ignore empty lines
		if len(input) > 0 {
			split := strings.Split(input, p.QaSep)
			switch len(split) {
			case 1:
				if strings.HasPrefix(input, p.TopicAnnounce) {
					subsectionID = strings.TrimPrefix(input, p.TopicAnnounce)
					qaSubsection = topic.GetSubsection(subsectionID)
				}
			default:
				// Question is in split[0] while answer in in split[1]. It may happen
				// the answer contains the separator so we have to join the different
				// elements.
				qaSubsection.AddEntry(split[0], strings.Join(split[1:], p.QaSep))
				topic.SetSubsection(subsectionID, qaSubsection)
			}
		}
	}
	return topic
}
