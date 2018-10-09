package parsing

import (
	"testing"
)

func TestContainsOnlyAllowedCharacters(t *testing.T) {
	allowedStrings := []string{
		"1", "1:4", "1,2:5,10",
	}
	for i := 0; i < len(allowedStrings); i++ {
		if !containsOnlyAllowedCharacters(allowedStrings[i]) {
			t.Errorf("Expected %q to be considered as allowed string", allowedStrings[i])
		}
	}

	unallowedStrings := []string{
		"1:", ":1", "1s", ":",
	}
	for i := 0; i < len(unallowedStrings); i++ {
		if containsOnlyAllowedCharacters(unallowedStrings[i]) {
			t.Errorf("Expected %q to be considered as unallowed string", unallowedStrings[i])
		}
	}
}

func TestParseNumberSerie(t *testing.T) {
	series := []string{
		"1",
		"1,3",
		"1:4",
		"1:3,5:7,10",
		"1:3:1:3",
		"1:3:2",
	}
	expectedResults := [][]int{
		{1},
		{1, 3},
		{1, 2, 3, 4},
		{1, 2, 3, 5, 6, 7, 10},
		{1, 2, 3, 2, 1, 2, 3},
		{1, 2, 3, 2},
	}

	for i := 0; i < len(series); i++ {
		computed, err := ParseNumberSerie(series[i])
		if err != nil {
			t.Errorf("Input string: %s\n", series[i])
			t.Errorf("Computed: %v\n", computed)
			t.Errorf("Expected: %v\n", expectedResults[i])
			t.Errorf("parsing a correct string %s should not result in error %v", series[i], err)
		}

		if len(computed) != len(expectedResults[i]) {
			t.Errorf("Input string: %s\n", series[i])
			t.Errorf("Computed: %v\n", computed)
			t.Errorf("Expected: %v\n", expectedResults[i])
			t.Fatalf("For string %s, expected %d elements, received an array with %d", series[i], len(expectedResults[i]), len(computed))
		}
		for j := 0; j < len(computed); j++ {
			if computed[j] != expectedResults[i][j] {
				t.Errorf("Input string: %s\n", series[i])
				t.Errorf("Computed: %v\n", computed)
				t.Errorf("Expected: %v\n", expectedResults[i])
				t.Errorf("Serie computed for %s is invalid. Expected %d at position %d but got %d", series[i], expectedResults[i][j], j, computed[j])
			}
		}
	}
}
