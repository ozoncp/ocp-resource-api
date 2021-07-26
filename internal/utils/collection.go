package utils

import (
	"errors"
	"github.com/ozoncp/ocp-resource-api/internal/models"
)

var ErrBatchSizeIsNull = errors.New("batch size should be greater that 0")
var ErrSourceIsNil = errors.New("source should not be nil")
var ErrDuplicatedKey = errors.New("key should not be duplicated in result map")
var ErrFilterOutValuesIsNil = errors.New("filterOutValues should not be nil")

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
		return nil, ErrSourceIsNil
	}
	if filterOutValues == nil {
		return nil, ErrFilterOutValuesIsNil
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

func batchBounds(length uint, batchSize uint, i uint) (uint, uint) {
	start := i * batchSize
	end := (i + 1) * batchSize
	if end > length {
		end = length
	}
	return start, end
}

func calcChunkSize(sourceLength uint, batchSize uint) uint {
	batchCount := sourceLength / batchSize
	if sourceLength%batchSize != 0 {
		batchCount += 1
	}
	return batchCount
}

func ReverseMapIntToInt(source map[int]int, forceRewrite bool) (map[int]int, error) {
	if source == nil {
		return nil, ErrSourceIsNil
	}
	result := make(map[int]int, len(source))
	for key, value := range source {
		_, ok := result[value]
		if ok && !forceRewrite {
			return nil, ErrDuplicatedKey
		}
		result[value] = key
	}
	return result, nil
}

func SplitToBulksInt(source []int, batchSize uint) ([][]int, error) {
	if batchSize == 0 {
		return nil, ErrBatchSizeIsNull
	}
	if source == nil {
		return nil, ErrSourceIsNil
	}
	batchCount := calcChunkSize(uint(len(source)), batchSize)
	result := make([][]int, 0, batchCount)
	for i := uint(0); i < batchCount; i++ {
		start, end := batchBounds(uint(len(source)), batchSize, i)
		result = append(result, source[start:end])
	}
	return result, nil
}

func SplitToBulksResource(source []models.Resource, batchSize uint) ([][]models.Resource, error) {
	if batchSize == 0 {
		return nil, ErrBatchSizeIsNull
	}
	if source == nil {
		return nil, ErrSourceIsNil
	}
	sourceLength := uint(len(source))
	batchCount := calcChunkSize(sourceLength, batchSize)
	result := make([][]models.Resource, 0, batchCount)
	for i := uint(0); i < batchCount; i++ {
		start, end := batchBounds(sourceLength, batchSize, i)
		result = append(result, source[start:end])
	}
	return result, nil
}

func SliceToMapResource(source []models.Resource, forceRewrite bool) (map[uint64]models.Resource, error) {
	if source == nil {
		return nil, ErrSourceIsNil
	}
	result := make(map[uint64]models.Resource, len(source))
	for _, resource := range source {
		_, ok := result[resource.Id]
		if ok && !forceRewrite {
			return nil, ErrDuplicatedKey
		}
		result[resource.Id] = resource
	}
	return result, nil
}
