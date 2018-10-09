package tools

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/fatih/color"
)

// WriteInBoldWhite writes a text in bold white color. This set of functions
// (WriteInXXX) makes possible to produced colored text on the same line.
// So don't forget to add the return carriage in the text when needed.
func WriteInBoldWhite(msg string) {
	c := color.New(color.FgWhite).Add(color.Bold)
	c.Printf(msg)
}

// WriteInGreen writes a text in green color. This set of functions
// (WriteInXXX) makes possible to produced colored text on the same line.
// So don't forget to add the return carriage in the text when needed.
func WriteInGreen(msg string) {
	c := color.New(color.FgGreen)
	c.Printf(msg)
}

// WriteInRed writes a text in red color. This set of functions
// (WriteInXXX) makes possible to produced colored text on the same line.
// So don't forget to add the return carriage in the text when needed.
func WriteInRed(msg string) {
	c := color.New(color.FgRed)
	if IsChavaModeActivated() {
		c = color.New(color.FgCyan)
	}
	c.Printf(msg)
}

// WriteInCyan writes a text in cyan color. This set of functions
// (WriteInXXX) makes possible to produced colored text on the same line.
// So don't forget to add the return carriage in the text when needed.
func WriteInCyan(msg string) {
	c := color.New(color.FgCyan)
	c.Printf(msg)
}

// WriteInMagenta writes a text in magenta color. This set of functions
// (WriteInXXX) makes possible to produced colored text on the same line.
// So don't forget to add the return carriage in the text when needed.
func WriteInMagenta(msg string) {
	c := color.New(color.FgMagenta)
	c.Printf(msg)
}

// WriteInYellow writes a text in yellow color. This set of functions
// (WriteInXXX) makes possible to produced colored text on the same line.
// So don't forget to add the return carriage in the text when needed.
func WriteInYellow(msg string) {
	c := color.New(color.FgYellow)
	c.Printf(msg)
}

// Error gives a convenient and uniform way to report error to the user console.
// This function uses ASCII code to display color so it can make log files look
// like a little bit funny. Use less -R to view log files with colors.
func Error(err error, msg string) {
	c := color.New(color.FgRed)
	if IsChavaModeActivated() {
		c = color.New(color.FgCyan)
	}
	switch err {
	case nil:
		c.Printf("[ERROR] %s\n", msg)
	default:
		errMsg := err.Error()
		errMsg = strings.Replace(errMsg, `\n`, "\n", -1)
		errMsg = strings.Replace(errMsg, `\t`, "\t", -1)
		c.Printf("[ERROR] %s : %s\n", msg, errMsg)
	}
}

// Info gives a convenient and uniform way to show info trace to the user console.
// This function uses ASCII code to display color so it can make log files look
// like a little bit funny. Use less -R to view log files with colors.
// Carriage return is added at the end of the text.
func Info(msg string) {
	whiteBold := color.New(color.FgWhite).Add(color.Bold)
	whiteBold.Printf("%s\n", msg)
}

// TitleNoNewLine displays a bold text in white and a ':' so it can
// be used to make some enum for instance.
// No carriage return is added at the end of the text.
func TitleNoNewLine(msg string) {
	whiteBold := color.New(color.FgWhite).Add(color.Bold)
	whiteBold.Printf("%s: ", msg)
}

// Section give a delimitation to annonce something with a colored text
// No carriage return is added at the end of the text.
func Section(title string) {
	c := color.New(color.FgGreen).Add(color.Bold)
	c.Println(title)
	c.Println(strings.Repeat("-", utf8.RuneCountInString(title)))
}

// Warning gives a convenient and uniform way to show warning trace to the user console.
// This function uses ASCII code to display color so it can make log files look
// like a little bit funny. Use less -R to view log files with colors.
// Carriage return is added to the end of the text.
func Warning(msg string) {
	magenta := color.New(color.FgMagenta)
	magenta.Printf("[WARN] %s\n", msg)
}

// Debug gives a convenient and uniform way to show debug trace to the user console.
// This function uses ASCII code to display color so it can make log files look
// like a little bit funny. Use less -R to view log files with colors.
// Carriage return is added to the end of the text.
func Debug(msg string) {
	if IsDebugActivated() {
		yellow := color.New(color.FgYellow)
		yellow.Printf("[DEBUG] %s\n", msg)
	}
}

// Advise gives a convenient and uniform way to display advices to the end user.
// This function uses ASCII code to display color so it can make log files look
// like a little bit funny. Use less -R to view log files with colors.
// Carriage return is added to the end of the text.
func Advise(msg string) {
	c := color.New(color.FgYellow)
	c.Printf(fmt.Sprintf("[Advice] %s\n", msg))
}

// RunAction gives a convenient and uniform way to announce an action to the end user.
// This function uses ASCII code to display color so it can make log files look
// like a little bit funny. Use less -R to view log files with colors.
// Carriage return is added to the end of the text.
func RunAction(msg string) {
	cyan := color.New(color.FgCyan)
	cyan.Printf("%s\n", msg)
}

// PositiveStatus prints text in green so the user perceives this message as a
// positive message.
// Carriage return is added to the end of the text.
func PositiveStatus(s string) {
	green := color.New(color.FgGreen)
	green.Printf("%s\n", s)
}

// NegativeStatus prints text in red so the user perceives this message as a
// potential error and she/he should really read this message.
// Carriage return is added to the end of the text.
func NegativeStatus(s string) {
	red := color.New(color.FgRed)
	red.Printf("%s\n", s)
}

// Summary creates a block of text to summarise a list of things.
// Each string passed in parameter will be on a distinct line.
func Summary(msg ...string) {
	c := color.New(color.FgCyan).Add(color.Bold)
	longuest := 0
	for _, t := range msg {
		if len(t) > longuest {
			longuest = len(t)
		}
	}
	for _, t := range msg {
		additionalSpaces := getRepeatedPattern(" ", longuest, len(t))
		c.Printf("%s%s\n", t, additionalSpaces)
	}
}

// helper function to repeat a pattern N times. Minimum is one even if
// max is less than limit.
func getRepeatedPattern(pattern string, max int, limit int) string {
	// Be sure to have the pattern at least once
	s := pattern
	if max > limit {
		s = strings.Repeat(pattern, max-limit)
	}
	return s
}

// Success reports an action with a message in green. The display is done
// like for the Unix Services that start. It can be useful when writing
// a list of actions and summarises success/failure.
// For instance:
// Starting Network............................[OK]
// where OK is in green.
func Success(msg, successMsg string) {
	green := color.New(color.FgGreen)
	whiteBold := color.New(color.FgWhite).Add(color.Bold)
	dots := getRepeatedPattern(".", 80, utf8.RuneCountInString(msg))
	whiteBold.Printf("%s %s", msg, dots)
	green.Printf(fmt.Sprintf("[%s]\n", successMsg))
}

// Failure reports an action with a message in red. The display is done
// like for the Unix Services that start. It can be useful when writing
// a list of actions and summarises success/failure.
// For instance:
// Starting Network............................[failed because XXX]
// where 'failed because XXX' is in red.
func Failure(msg, failureMsg string) {
	red := color.New(color.FgRed)
	whiteBold := color.New(color.FgWhite).Add(color.Bold)
	dots := getRepeatedPattern(".", 80, utf8.RuneCountInString(msg))
	whiteBold.Printf("%s %s", msg, dots)
	red.Printf(fmt.Sprintf("[%s]\n", failureMsg))
}

// OK is the same as Success with a predefined success message: OK.
func OK(msg string) {
	Success(msg, "OK")
}

// NOK is the same as Failure with a predefined message: NOK.
func NOK(msg string) {
	Failure(msg, "NOK")
}

// Announce writes to stdout a colored text underlined with equals sign.
func Announce(msg string) {
	cyan := color.New(color.FgCyan).Add(color.Bold)
	cyan.Println(msg)
	cyan.Println(strings.Repeat("=", utf8.RuneCountInString(msg)))
}

// question writes a colored text on the stdout. A newline is added after
// the question and a prompt character is also added. Default is to use '>'.
func question(text string, withPrompt bool) {
	cyan := color.New(color.FgCyan)
	cyan.Println(text)
	if withPrompt {
		cyan.Printf(" > ")
	}
}

// QuestionWithPrompt see question function
func QuestionWithPrompt(text string) {
	question(text, true)
}

// QuestionWithNoPrompt see question function
func QuestionWithNoPrompt(text string) {
	question(text, false)
}
