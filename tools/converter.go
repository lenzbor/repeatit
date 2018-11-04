package tools

import "fmt"

// ConvertToStringsArray takes an array of int and convert it to
// an array of strings. The specificity of this transformation is that
// we will prefix the number that do not have the good number of digits
// with front zero(s).
func ConvertToStringsArray(v []int, numberOfDigits int) []string {
	output := make([]string, len(v))
	for i := 0; i < len(v); i++ {
		output[i] = fmt.Sprintf("%0*d", numberOfDigits, v[i])
	}
	return output
}
