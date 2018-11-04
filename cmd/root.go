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

	"github.com/boris-lenzinger/repeatit/datamodel"
	"github.com/boris-lenzinger/repeatit/engine"
	"github.com/boris-lenzinger/repeatit/parsing"
	"github.com/boris-lenzinger/repeatit/tools"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	keySourceFile = "sourceFile"
)

// global string that stores the path to the configuration file in
// case you want to set a specific configuration file. Else the code
// will choose $HOME/.repeatit.yaml
var cfgFile string

// Debug activates the debug mode for the tool.
var Debug bool

// global parameter storing the path to the input that stores the
// sentences to learn
var pathToLessonsFile string

// params is a global variable of the command interpreter.
//
var params datamodel.InterrogationParameters

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "repeatit <path-to-lesson-file>",
	Short: "A tool to help you to learn things based on a simple thing: REPETITION !!.",
	Long:  ``,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			if viper.GetString(keySourceFile) == "" {
				tools.NegativeStatus(fmt.Sprintf("Please set a file on the command line or in your $HOME/.repeat.yaml with the key %q", keySourceFile))
				os.Exit(1)
			}
			pathToLessonsFile = viper.GetString(keySourceFile)
		}
		tools.Debug(fmt.Sprintf("[root] Arguments received: %+v\n", args))
		params = datamodel.NewInterrogationParameters()
		if interactive {
			params.SetInteractive()
		}
		exists, err := tools.FileExists(pathToLessonsFile)
		if err != nil {
			tools.Error(err, fmt.Sprintf("error while checking if lessons file %q exists", pathToLessonsFile))
			os.Exit(1)
		}
		if !exists {
			tools.NegativeStatus(fmt.Sprintf("File %q does not exist. Please set a file that exists.", pathToLessonsFile))
			os.Exit(1)
		}

		params.SetLessonsFile(pathToLessonsFile)
		if Debug {
			viper.Set("debug", true)
		}
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		parsingParameters := datamodel.TopicParsingParameters{
			LessonAnnounce:   viper.GetString("announcementForLessons"),
			SentenceAnnounce: viper.GetString("announcementForSentences"),
			QaSep:            viper.GetString("qaSep"),
		}
		t, err := parsing.ParseLanguageFile(pathToLessonsFile, parsingParameters)
		if err != nil {
			tools.NegativeStatus(fmt.Sprintf("failed to parse file %q due to %v", pathToLessonsFile, err))
			os.Exit(0)
		}
		engine.StartEngine(t)
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		tools.Debug("[root] Calling PersistentPostRun")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	tools.Debug("[root] Execute()")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	tools.Debug(fmt.Sprintf("[root] Is it interactive ? %t\n", interactive))
}

func init() {
	tools.Debug("[root] init()")
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.repeatit.yaml)")
	tools.Debug("Adding support for interactive flag.")
	rootCmd.PersistentFlags().BoolVarP(&interactive, "interactive", "", false, `If set, you will have to press Return to get the answer.
This allows you to be in a learning way or enforcing your knowledge. It lets you time to search in your
memory and answer when you feel ready.
If this flag is not set, you will not have to press the Return key and you
simply have to wait for a  given time. Questions and answers flow with a time
interval between them. See -t for details about time.`)
	rootCmd.PersistentFlags().StringVarP(&pathToLessonsFile, "lessons", "", "", "the path to the file containing the lessons.")
	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "", false, "Enables debug mode on the client.")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().StringP()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	tools.Debug("[root] initConfig()")
	if cfgFile != "" {
		tools.Debug(fmt.Sprintf("[root] Using configuration files passed in parameter %s \n", cfgFile))
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
		tools.Debug(fmt.Sprintf("[root] Using config file:  %s", viper.ConfigFileUsed()))
	} else {
		tools.Error(err, fmt.Sprintf("failed to use config file %q", viper.ConfigFileUsed()))
		os.Exit(1)
	}
}

func checkFileOrFail() {
	exists, err := tools.FileExists(pathToLessonsFile)
	if err != nil {
		tools.Error(err, fmt.Sprintf("error while checking if %q exists", pathToLessonsFile))
		os.Exit(1)
	}
	if !exists {
		tools.NegativeStatus(fmt.Sprintf("file to add the lesson (%s) does not exist", pathToLessonsFile))
		os.Exit(1)
	}
}
