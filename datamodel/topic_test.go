package datamodel

import "testing"

func TestComputeRangesOnArrayOfInts(t *testing.T) {
	tests := []struct {
		input    []int
		expected string
	}{
		{
			input:    []int{1, 2, 3, 4, 5, 6},
			expected: "1:6",
		},
		{
			input:    []int{1, 2, 4, 5, 6},
			expected: "1:2,4:6",
		},
		{
			input:    []int{1, 3, 4, 5, 6},
			expected: "1,3:6",
		},
	}
	for i := 0; i < len(tests); i++ {
		test := tests[i]
		computed := computeRangesOnArrayOfInts(test.input)
		if computed != test.expected {
			t.Errorf("for input %v, was expecting %s but received %s", test.input, test.expected, computed)
		}
	}

}
