package tests

import "regexp"

// EmptyLine is a regular expression to test if a line is empty or not.
var EmptyLine = regexp.MustCompile("^\\s*$")

// Separator is a regular expression to test if a line is a separator used when displaying
// vocabulary.
var Separator = regexp.MustCompile("^---------------------------$")

// Loop is the regular expression to test if a line is the announcement of a loop with its count
// in the process of interrogation of the user.
var Loop = regexp.MustCompile("^Loop \\([0-9]+/[0-9]+\\)$")

// NbOfQuestions is the regular expression to test if a line is the announcement of a the
// number of questions in the process of interrogation of the user.
var NbOfQuestions = regexp.MustCompile("^Nb of questions: [0-9]+$")

// LimitReached is a regular expression to test if a line of text is the annoncement of
// reaching the maximum number of loops.
var LimitReached = regexp.MustCompile("^Limit reached. Exiting. Number of loops set to: [0-9]+$")
