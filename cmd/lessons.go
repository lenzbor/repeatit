// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/boris-lenzinger/repeatit/datamodel"
	"github.com/boris-lenzinger/repeatit/engine"
	"github.com/boris-lenzinger/repeatit/parsing"
	"github.com/boris-lenzinger/repeatit/tools"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var interactive bool

// lessonsCmd represents the lessons command
var lessonsCmd = &cobra.Command{
	Use:   "lessons [numbers]",
	Short: "Requires repetition for lessons set on the command line",
	Long: `This commands requires to repeat a series of lessons.
The numbers can be dispatched as follow:
  * n:m requires to repeat the lessons n to m
  * n,m requires the lesson n and m
  * you can combine the above syntaxes to generate complex combinations that match your needs
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			tools.NegativeStatus("Please supply lessons number. Check the syntax of the command if you don't know how to set lessons number.")
			os.Exit(1)
		}
		lessonsToLearn := args[0]
		fmt.Printf("lessons required : %s\n", lessonsToLearn)
		fmt.Printf("[lessons] Is it interactive ? %t\n", params.IsInteractive())
		fmt.Printf("[lessons] Path to file to handle: %s\n", params.GetLessonsFile())

		// file existence has already been checked by the root command
		f, err := os.Open(params.GetLessonsFile())
		if err != nil {
			tools.Error(err, fmt.Sprintf("failed to parse the lessons file %q", params.GetLessonsFile()))
			os.Exit(1)
		}
		lessonsRange, err := parsing.ParseNumberSerie(lessonsToLearn)
		if err != nil {
			tools.Error(err, "the arguments passed do not seem to be a list of numbers")
			os.Exit(1)
		}
		lessonNumbers := make([]string, len(lessonsRange))
		for i := 0; i < len(lessonsRange); i++ {
			s := strconv.Itoa(lessonsRange[i])
			if lessonsRange[i] < 10 {
				s = "0" + s
			}
			lessonNumbers[i] = s
		}
		parsingParameters := datamodel.TopicParsingParameters{
			LessonAnnounce:   viper.GetString("announcementForLessons"),
			SentenceAnnounce: viper.GetString("announcementForSentences"),
			QaSep:            viper.GetString("qaSep"),
		}
		topic, err := parsing.ParseTopic(f, parsingParameters)
		if err != nil {
			tools.NegativeStatus(fmt.Sprintf("Parsing of %q has failed due to %v", params.GetLessonsFile(), err))
		}
		tools.Debug(topic.String())
		qa := topic.BuildVocabularyQuestionsSet(lessonNumbers...)
		engine.AskQuestions(qa, params)
	},
}

func init() {
	rootCmd.AddCommand(lessonsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lessonsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lessonsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
