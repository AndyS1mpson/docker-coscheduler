package slices_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/AndyS1mpson/docker-coscheduler/internal/utils/slices"
)

type elem struct {
	id        int64
	anotherID int
}

func TestIndexByKey_byElemID(t *testing.T) {
	var (
		idExtractor = func(elem elem) int64 {
			return elem.id
		}
		emptyMap = map[int64]elem{}
	)

	testCases := []struct {
		name      string
		in        []elem
		extractor func(elem elem) int64
		expected  map[int64]elem
	}{
		// Пустой массив на вход
		{
			name:      "empty input",
			in:        []elem{},
			extractor: idExtractor,
			expected:  emptyMap,
		},
		// Нулевой массив на вход
		{
			name:      "input is nil",
			in:        nil,
			extractor: idExtractor,
			expected:  emptyMap,
		},
		// Nil в качестве функции экстрактора
		{
			name: "keyExtractor is nil",
			in: []elem{
				{id: 1, anotherID: 3},
				{id: 2, anotherID: 4},
			},
			extractor: nil,
			expected:  emptyMap,
		},
		// Индексирования по ID
		{
			name: "index by id",
			in: []elem{
				{id: 1, anotherID: 3},
				{id: 2, anotherID: 4},
			},
			extractor: idExtractor,
			expected: map[int64]elem{
				1: {id: 1, anotherID: 3},
				2: {id: 2, anotherID: 4},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := IndexByKey(tc.in, tc.extractor)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestIndexByKey_byElemAnotherID(t *testing.T) {
	anotherIDExtractor := func(elem elem) int {
		return elem.anotherID
	}

	testCases := []struct {
		name      string
		in        []elem
		extractor func(elem elem) int
		expected  map[int]elem
	}{
		// Индексирование по AnotherID без пересечений
		{
			name: "index by another id",
			in: []elem{
				{id: 1, anotherID: 3},
				{id: 2, anotherID: 4},
			},
			extractor: anotherIDExtractor,
			expected: map[int]elem{
				3: {id: 1, anotherID: 3},
				4: {id: 2, anotherID: 4},
			},
		},
		// Индексирование по AnotherID с пересечениями
		{
			name: "index by another id with removing duplicates",
			in: []elem{
				{id: 1, anotherID: 4},
				{id: 2, anotherID: 4},
				{id: 3, anotherID: 4},
			},
			extractor: anotherIDExtractor,
			expected: map[int]elem{
				4: {id: 3, anotherID: 4},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := IndexByKey(tc.in, tc.extractor)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
