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
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var isInteractive bool
var isLinear bool
var waitTime int
var isSummary bool
var topics []string
var isReversed bool
var loopCount int

const defaultTopicsAnnouncer = "### "
const defaultQaSeparator = ";"
const defaultWaitTime = 2000
const defaultLoopCount = 1

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "repeatit",
	Short: "A simple tool to help you to memorize vocabulary or text with a simple method: REPETITION !!",
	Long: `This tool helps you to memorize different kind of things.

The main idea is that repetition helps (some) people to remember. You can
use this with different repetition modes: random, linear -useful for a text-,
set a time between interrogation of each word, etc.

You can configure your application with a configuration file in your home:
  $HOME/.repeatit.yml`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if len(os.Args) < 2 {
			cmd.Help()
			os.Exit(1)
		}

		// Creer un objet fichier et tester si on peut le lire
		filename := os.Args[1]
		file, err := os.Open(filename)
		if err != nil {
			fmt.Printf("Open of the source file failed: %v\n", err)
			cmd.Help()
			os.Exit(1)
		}

		p := InterrogationParameters{
			interactive: isInteractive,
			wait:        time.Duration(waitTime),
			reversed:    isReversed,
			limit:       loopCount,
			qachan:      make(chan string),
			command:     make(chan string),
			publisher:   make(chan string),
		}

		switch {
		case isSummary:
			p.mode = summary
		case isLinear:
			p.mode = linear
		default:
			p.mode = random
		}

		p.in = os.Stdin
		p.out = os.Stdout
		p.subsections = topics

		tpp := TopicParsingParameters{
			TopicAnnounce: defaultTopicsAnnouncer,
			QaSep:         defaultQaSeparator,
		}
		topic := ParseTopic(file, tpp)
		file.Close()

		out := p.GetOutputStream()
		if p.IsSummaryMode() {
			list := topic.GetSubsectionsName()
			if len(list) == 0 {
				fmt.Fprintf(out, "No topic found in this file")
				return
			}
			fmt.Fprintln(out, "List of topics:")
			fmt.Fprintln(out, "===============")
			for i := 0; i < len(list); i++ {
				fmt.Fprintf(out, "  * %s\n", list[i])
			}
			return
		}

		qa := topic.BuildQuestionsSet(p.GetListOfSubsections()[:]...)

		AskQuestions(qa, p)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	defaultTimeInSec := strconv.Itoa((defaultWaitTime / 1000))

	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.repeatit.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolVarP(&isInteractive, "interactive", "i", false, "If set, you will have to press Return to get the answer. This allows you to be in a learning way or enforcing your knowledge. If this flag is not set, you will not have to press the Return key and you simply have to wait for a given time. See -t for details about time.")
	rootCmd.Flags().IntVarP(&waitTime, "time", "t", defaultWaitTime, "the time to wait between 2 questions. Default is "+defaultTimeInSec+" seconds. The time you set is in milliseconds.")
	rootCmd.Flags().BoolVarP(&isSummary, "summary", "s", false, "ask to show the different topics of  the file, no more. Execution stops after this. Sections are supposed to start with "+defaultTopicsAnnouncer+".")
	rootCmd.Flags().StringArrayVarP(&topics, "list", "l", []string{}, "ask to be questionned only on the topics that are listed here. The topics must be separated with a comma.")
	rootCmd.Flags().BoolVarP(&isReversed, "reversed", "r", false, "reverts the questioning. This is like a Jeopardy in fact. The right column becomes the questions while the left column becomes the answer.")
	rootCmd.Flags().BoolVarP(&isLinear, "linear", "", false, "requires the questions to be made as they appear in the source")
	rootCmd.Flags().IntVarP(&loopCount, "loop", "", defaultLoopCount, "set the number of loops. Default is "+strconv.Itoa(defaultLoopCount))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".repeatit" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".repeatit")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
