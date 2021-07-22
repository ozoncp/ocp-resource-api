package utils

import (
	"errors"
)

func sliceContainsIn(source []int, value int) (bool, error) {
	for i := 0; i < len(source); i++ {
		if source[i] == value {
			return true, nil
		}
	}
	return false, nil
}

func FilterOutInt(source []int, filterOutValues []int) ([]int, error) {
	if source == nil {
		return nil, errors.New("source should not be nil")
	}
	if filterOutValues == nil {
		return nil, errors.New("filterOutValues should not be nil")
	}
	result := make([]int, 0, len(source))
	for i := 0; i < len(source); i++ {
		value := source[i]
		if contains, _ := sliceContainsIn(filterOutValues, value); !contains {
			result = append(result, value)
		}
	}
	return result, nil
}

func SplitInt(source []int, batchSize int) ([][]int, error) {
	if batchSize <= 0 {
		return nil, errors.New("batch size should be greater that 0")
	}
	if source == nil {
		return nil, errors.New("source should not be nil")
	}
	batchCount := calcChunkSize(source, batchSize)
	result := make([][]int, 0, batchCount)
	for i := 0; i < batchCount; i++ {
		start, end := batchBounds(len(source), batchSize, i)
		result = append(result, source[start:end])
	}
	return result, nil
}

func batchBounds(length int, batchSize int, i int) (int, int) {
	start := i * batchSize
	end := (i + 1) * batchSize
	if end > length {
		end = length
	}
	return start, end
}

func calcChunkSize(source []int, batchSize int) int {
	batchCount := len(source) / batchSize
	if len(source)%batchSize != 0 {
		batchCount += 1
	}
	return batchCount
}

func ReverseMapIntToInt(source map[int]int, forceRewrite bool) (map[int]int, error) {
	if source == nil {
		return nil, errors.New("source should not be nil")
	}
	result := make(map[int]int, len(source))
	for key, value := range source {
		_, ok := result[value]
		if ok && !forceRewrite {
			return nil, errors.New("key should not be duplicated in result map")
		}
		result[value] = key
	}
	return result, nil
}
