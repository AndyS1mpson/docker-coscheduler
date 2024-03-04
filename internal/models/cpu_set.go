package models

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

// CPUSet представляет набор CPU
type CPUSet struct {
	From  int64 `json:"from"`  // Начальный индекс
	Count int64 `json:"count"` // Количество CPU в наборе
}

// NewCPUSet конструктор для CPUSet
func NewCPUSet(str string) (*CPUSet, error) {
	re := regexp.MustCompile(`(\d+)-(\d+)`)
	matches := re.FindStringSubmatch(str)

	if len(matches) != 3 {
		return &CPUSet{}, errors.New("can not match params")
	}

	from, err := strconv.ParseInt(matches[1], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("convert from: %w", err)
	}

	to, err := strconv.ParseInt(matches[2], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("convert count: %w", err)
	}

	count := to - from + 1
	if count <= 0 || from < 0 {
		return nil, errors.New("count and from must be positive")
	}

	return &CPUSet{From: from, Count: count}, nil
}

// AsString возвращает строковое представление набора CPU
func (c *CPUSet) AsString() string {
	return fmt.Sprintf("%d-%d", c.From, c.From+c.Count-1)
}
