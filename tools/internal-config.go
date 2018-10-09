package tools

import "github.com/spf13/viper"

// IsDebugActivated tells if the debug configuration key has been
// set to debug or not.
func IsDebugActivated() bool {
	return viper.GetBool("debug")
}

// IsChavaModeActicated tells if the mode for people who have problem to
// read text in red is activated or not. It is named after a friend who
// has this problem and told me he has issues in reading the error message
// reported in red.
func IsChavaModeActivated() bool {
	return viper.GetBool("chavaMode")
}
