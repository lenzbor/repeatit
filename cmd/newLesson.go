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
	"github.com/boris-lenzinger/repeatit/tools"
	"github.com/spf13/cobra"
)

// newLessonCmd represents the newLesson command
var newLessonCmd = &cobra.Command{
	Use:   "lesson",
	Short: "Creates a new lesson in your lang book",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		checkFileOrFail()
		f, err := os.Open(pathToLessonsFile)
		if err != nil {
			tools.Error(err, fmt.Sprintf("error while trying to open %q", pathToLessonsFile))
			os.Exit(1)
		}
		lessons, err := datamodel.LoadLessons(f)
		if err != nil {
			tools.Error(err, fmt.Sprintf("failed to load the lessons file %q", pathToLessonsFile))
			os.Exit(1)
		}
		f.Close()
		tools.Info(fmt.Sprintf("Loaded %d lessons", lessons.GetLessonsCount()))

		tools.QuestionWithPrompt("Title in the learning language")
		// Read the input and reject empty value
		titleInLearningLanguage := tools.ReadFromStdin(false)
		tools.QuestionWithPrompt("Title in your native language")
		titleInNativeLanguage := tools.ReadFromStdin(false)
		tools.Info(fmt.Sprintf("User has entered %q for the original title and %q for the title in its native language", titleInLearningLanguage, titleInNativeLanguage))
		res := datamodel.Resource{Learning: titleInLearningLanguage, Native: titleInNativeLanguage}
		fmt.Printf("Loaded res as %+v", res)

	},
}

func init() {
	newCmd.AddCommand(newLessonCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newLessonCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newLessonCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
