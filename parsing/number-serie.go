package parsing

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ParseNumberSerie takes a string describing a number sequence
// and create the real sequence.
// The expected syntax of the sequence is:
//   * n:m will be expanded as n, n+1, n+2, ..., m
//   * n,m remains as n, m
//   * syntax can be combined as n:m,p,q:r
// If the sequence has an invalid order (n is greater than m for
// instance), an error is returned.
// If an element of the string is not an integer nor a comma nor a
// colon, an error is returned.
func ParseNumberSerie(serie string) ([]int, error) {
	if !containsOnlyAllowedCharacters(serie) {
		return []int{}, fmt.Errorf("the serie entered contained unauthorized characters")
	}
	out := []int{}
	var n int
	splitted := strings.Split(serie, ",")
	for i := 0; i < len(splitted); i++ {
		contiguous := strings.Split(splitted[i], ":")
		if len(contiguous) == 1 {
			n, _ = strconv.Atoi(contiguous[0])
			out = append(out, n)
		} else {
			for k := 0; k < len(contiguous)-1; k++ {
				// the first element can be larger than the following: in this case
				// we are making a decreasing questionning !
				min, _ := strconv.Atoi(contiguous[k])
				max, _ := strconv.Atoi(contiguous[k+1])
				if min == max {
					goto nextRound
				}
				// Warning: 1:3:1 would return 1:2:3:3:2:1 if we leave the code this
				// way. The fact that we are in a multiple : sequence changes the context
				if k >= 1 {
					if min > max {
						// 1:3:1 -> when reached 3, we want to start from 2, not 3...
						min = min - 1
					} else {
						// 5:3:5 must not be 5, 4, 3, 3, 4, 5
						min = min + 1
					}
				}
				if min > max {
					for j := min; j >= max; j-- {
						out = append(out, j)
					}
				} else {
					for j := min; j <= max; j++ {
						out = append(out, j)
					}
				}
			}
		nextRound:
		}
	}
	return out, nil
}

// Checks that a string starts with a number and is a sequence of
// numbers separated by comma or colon. Sequence of 1 number is allowed.
func containsOnlyAllowedCharacters(s string) bool {
	r, _ := regexp.Compile("^([0-9]+([,:][0-9]+)*)+$")
	return r.MatchString(s)
}
