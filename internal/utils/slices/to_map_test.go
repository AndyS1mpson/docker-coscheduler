package slices_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/AndyS1mpson/docker-coscheduler/internal/utils/slices"
)

func TestToMap(t *testing.T) {
	testCases := []struct {
		name     string
		in       []int64
		expected map[int64]struct{}
	}{
		// tagID 16677 на вход перед пустой список
		// tmsTestType unit
		{
			name:     "empty in",
			in:       []int64{},
			expected: map[int64]struct{}{},
		},
		// tagID 16677 на вход перед нулевой список
		// tmsTestType unit
		{
			name:     "nil in",
			in:       nil,
			expected: map[int64]struct{}{},
		},
		// tagID 16677 на вход передан список с повторяющимися значениями
		// tmsTestType unit
		{
			name: "index non unique list",
			in:   []int64{1, 1, 2, 1, 3, 4, 3, 4},
			expected: map[int64]struct{}{
				1: {},
				2: {},
				3: {},
				4: {},
			},
		},
		// tagID 16677 на вход передан список уникальных значений
		// tmsTestType unit
		{
			name: "index unique list",
			in:   []int64{1, 2, 3, 4},
			expected: map[int64]struct{}{
				1: {},
				2: {},
				3: {},
				4: {},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := ToMap(tc.in)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
