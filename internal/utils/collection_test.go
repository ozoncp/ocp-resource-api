package utils

import (
	"errors"
	"testing"

	"github.com/ozoncp/ocp-resource-api/internal/models"
)

func equalsNilValues(a interface{}, b interface{}) bool {
	return !(a == nil && b != nil) && !(a != nil && b == nil)
}

func sliceIntIntEquals(a [][]int, b [][]int) bool {
	if !equalsNilValues(a, b) || len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if !sliceIntEquals(a[i], b[i]) {
			return false
		}
	}
	return true
}

func sliceIntEquals(a []int, b []int) bool {
	if !equalsNilValues(a, b) || len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func sliceResourceEquals(a []models.Resource, b []models.Resource) bool {
	if !equalsNilValues(a, b) || len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i].Id != b[i].Id ||
			a[i].UserId != b[i].UserId ||
			a[i].Status != b[i].Status ||
			a[i].Type != b[i].Type {
			return false
		}
	}
	return true
}

func sliceResourceMatrixEquals(a [][]models.Resource, b [][]models.Resource) bool {
	if !equalsNilValues(a, b) || len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if !sliceResourceEquals(a[i], b[i]) {
			return false
		}
	}
	return true
}

func mapIntIntEquals(a map[int]int, b map[int]int) bool {
	if !equalsNilValues(a, b) || len(a) != len(b) {
		return false
	}
	for key, aValue := range a {
		bValue := b[key]
		if aValue != bValue {
			return false
		}
	}
	return true
}

func mapUInt64ResourceEquals(a map[uint64]models.Resource, b map[uint64]models.Resource) bool {
	if !equalsNilValues(a, b) || len(a) != len(b) {
		return false
	}
	for key, aValue := range a {
		bValue := b[key]
		if aValue.Id != bValue.Id ||
			aValue.UserId != bValue.UserId ||
			aValue.Status != bValue.Status ||
			aValue.Type != bValue.Type {
			return false
		}
	}
	return true
}

func errorEquals(a error, b error) bool {
	if a == nil && b == nil {
		return true
	}
	if (a != nil && b == nil) || (a == nil && b != nil) {
		return false
	}
	return a.Error() == b.Error()
}

func assertSplitToBulksInt(t *testing.T, source []int, batchSize uint, expectedSlice [][]int, expectedErr error) {
	actualSlice, actualErr := SplitToBulksInt(source, batchSize)
	if !errorEquals(actualErr, expectedErr) {
		t.Fatalf("Error assertion failed. Actual '%v'. Expected '%v'", actualErr, expectedErr)
	}
	if !sliceIntIntEquals(actualSlice, expectedSlice) {
		t.Fatalf("Assertion failed. Actual %v. Excepted %v", actualSlice, expectedSlice)
	}
}

func assertFilterOut(t *testing.T, source []int, filterOutValues []int, expectedSlice []int, expectedErr error) {
	actualSlice, actualErr := FilterOutInt(source, filterOutValues)
	if !errorEquals(actualErr, expectedErr) {
		t.Fatalf("Error assertion failed. Actual '%v'. Expected '%v'", actualErr, expectedErr)
	}
	if !sliceIntEquals(actualSlice, expectedSlice) {
		t.Fatalf("Assertion failed. Actual %v. Excepted %v", actualSlice, expectedSlice)
	}
}

func assertMapReverseIntInt(t *testing.T, source map[int]int, rewrite bool, expectedMap map[int]int, expectedErr error) {
	actualMap, actualErr := ReverseMapIntToInt(source, rewrite)
	if !errorEquals(actualErr, expectedErr) {
		t.Fatalf("Error assertion failed. Actual '%v'. Expected '%v'", actualErr, expectedErr)
	}
	if !mapIntIntEquals(actualMap, expectedMap) {
		t.Fatalf("Assertion failed. Actual %v. Excepted %v", actualMap, expectedMap)
	}
}

// TODO rewrite old tests to a new format
func TestSplitToBulksInt(t *testing.T) {
	assertSplitToBulksInt(t, []int{1, 2, 3, 4, 5, 6}, 5, [][]int{{1, 2, 3, 4, 5}, {6}}, nil)
}

func TestSplitToBulksIntNil(t *testing.T) {
	assertSplitToBulksInt(t, nil, 5, nil, ErrSourceIsNil)
}

func TestSplitToBulksIntBatchSizeZero(t *testing.T) {
	assertSplitToBulksInt(t, []int{1}, 0, nil, ErrBatchSizeIsNull)
}

func TestSplitToBulksIntBatchWithSameLength(t *testing.T) {
	assertSplitToBulksInt(t, []int{1, 2, 3, 4, 5}, 5, [][]int{{1, 2, 3, 4, 5}}, nil)
}

func TestSplitToBulksIntBatchGreaterThanLength(t *testing.T) {
	assertSplitToBulksInt(t, []int{1}, 5, [][]int{{1}}, nil)
}

func TestFilterOutInt(t *testing.T) {
	assertFilterOut(t, []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 5}, []int{4}, nil)
}

func TestFilterOutIntSourceNil(t *testing.T) {
	assertFilterOut(t, nil, []int{1, 2, 3, 5}, nil, ErrSourceIsNil)
}

func TestFilterOutIntFilterValuesNil(t *testing.T) {
	assertFilterOut(t, []int{1, 2, 3, 5}, nil,
		nil, ErrFilterOutValuesIsNil)
}

func TestReverseMapIntToInt(t *testing.T) {
	assertMapReverseIntInt(t, map[int]int{1: 2, 3: 4}, false, map[int]int{2: 1, 4: 3}, nil)
}

func TestReverseMapIntToIntNilMap(t *testing.T) {
	assertMapReverseIntInt(t, nil, false, nil, ErrSourceIsNil)
}

func TestReverseMapIntToIntDuplicate(t *testing.T) {
	assertMapReverseIntInt(t, map[int]int{1: 2, 3: 2}, false, nil, ErrDuplicatedKey)
}

func TestReverseMapIntToIntDuplicateRewrite(t *testing.T) {
	assertMapReverseIntInt(t, map[int]int{1: 2, 3: 2}, true, map[int]int{2: 3}, nil)
}

func TestSplitToBulksResource(t *testing.T) {
	type args struct {
		slice     []models.Resource
		batchSize uint
	}
	testSlice := []models.Resource{
		models.NewResource(1, 1, 1, 1),
		models.NewResource(2, 2, 2, 2),
		models.NewResource(3, 3, 3, 3),
		models.NewResource(4, 4, 4, 4),
	}

	tests := []struct {
		name          string
		args          args
		expectedSlice [][]models.Resource
		expectedError error
	}{
		{
			name: "base",
			args: args{slice: testSlice, batchSize: uint(3)},
			expectedSlice: [][]models.Resource{
				testSlice[0:3],
				testSlice[3:],
			},
			expectedError: nil,
		},
		{
			name:          "batch greater than slice len",
			args:          args{slice: testSlice, batchSize: uint(5)},
			expectedSlice: [][]models.Resource{testSlice},
			expectedError: nil,
		},
		{
			name:          "slice is nil",
			args:          args{slice: nil, batchSize: 10},
			expectedSlice: nil,
			expectedError: ErrSourceIsNil,
		},
		{
			name:          "batchSize is 0",
			args:          args{slice: testSlice, batchSize: 0},
			expectedSlice: nil,
			expectedError: ErrBatchSizeIsNull,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := SplitToBulksResource(test.args.slice, test.args.batchSize)
			if !equalsNilValues(err, test.expectedError) && !errors.Is(err, test.expectedError) {
				t.Errorf("Error from SplitToBulksResource() = %v, expected %v", actual, test.expectedSlice)
			}
			if !sliceResourceMatrixEquals(actual, test.expectedSlice) {
				t.Errorf("SplitToBulksResource() = %v, expected %v", actual, test.expectedSlice)
			}
		})
	}
}

func TestSliceToMapResource(t *testing.T) {
	testSlice := []models.Resource{
		models.NewResource(1, 1, 1, 1),
		models.NewResource(2, 2, 2, 2),
		models.NewResource(3, 3, 3, 3),
		models.NewResource(1, 4, 4, 4),
	}

	baseMap := map[uint64]models.Resource{
		1: models.NewResource(1, 4, 4, 4),
		2: models.NewResource(2, 2, 2, 2),
		3: models.NewResource(3, 3, 3, 3),
	}

	type args struct {
		source       []models.Resource
		forceRewrite bool
	}
	tests := []struct {
		name        string
		args        args
		expected    map[uint64]models.Resource
		expectedErr error
	}{
		{
			name:        "with force rewrite",
			args:        args{source: testSlice, forceRewrite: true},
			expected:    baseMap,
			expectedErr: nil,
		},
		{
			name:        "without force rewrite and duplicated keys",
			args:        args{source: testSlice, forceRewrite: false},
			expected:    nil,
			expectedErr: ErrDuplicatedKey,
		},
		{
			name:        "with nil source",
			args:        args{source: nil, forceRewrite: true},
			expected:    nil,
			expectedErr: ErrSourceIsNil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := SliceToMapResource(test.args.source, test.args.forceRewrite)
			if !equalsNilValues(err, test.expectedErr) && !errors.Is(err, test.expectedErr) {
				t.Errorf("SliceToMapResource() error = %v, expected %v", err, test.expectedErr)
				return
			}
			if !mapUInt64ResourceEquals(actual, test.expected) {
				t.Errorf("SliceToMapResource() actual = %v, expected %v", actual, test.expected)
			}
		})
	}
}
