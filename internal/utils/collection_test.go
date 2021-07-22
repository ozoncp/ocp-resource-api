package utils

import (
	"errors"
	"testing"
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

func errorEquals(a error, b error) bool {
	if a == nil && b == nil {
		return true
	}
	if (a != nil && b == nil) || (a == nil && b != nil) {
		return false
	}
	return a.Error() == b.Error()
}

func assertSplitInt(t *testing.T, source []int, batchSize int, expectedSlice [][]int, expectedErr error) {
	actualSlice, actualErr := SplitInt(source, batchSize)
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

func TestSplitInt(t *testing.T) {
	assertSplitInt(t, []int{1, 2, 3, 4, 5, 6}, 5, [][]int{{1, 2, 3, 4, 5}, {6}}, nil)
}

func TestSplitIntNil(t *testing.T) {
	assertSplitInt(t, nil, 5, nil, errors.New("source should not be nil"))
}

func TestSplitIntBatchSizeZero(t *testing.T) {
	assertSplitInt(t, []int{1}, 0, nil, errors.New("batch size should be greater that 0"))
}

func TestSplitIntBatchWithSameLength(t *testing.T) {
	assertSplitInt(t, []int{1, 2, 3, 4, 5}, 5, [][]int{{1, 2, 3, 4, 5}}, nil)
}

func TestSplitIntBatchGreaterThanLenght(t *testing.T) {
	assertSplitInt(t, []int{1}, 5, [][]int{{1}}, nil)
}

func TestFilterOutInt(t *testing.T) {
	assertFilterOut(t, []int{1, 2, 3, 4, 5}, []int{1, 2, 3, 5}, []int{4}, nil)
}

func TestFilterOutIntSourceNil(t *testing.T) {
	assertFilterOut(t, nil, []int{1, 2, 3, 5}, nil, errors.New("source should not be nil"))
}

func TestFilterOutIntFilterValuesNil(t *testing.T) {
	assertFilterOut(t, []int{1, 2, 3, 5}, nil,
		nil, errors.New("filterOutValues should not be nil"))
}

func TestReverseMapIntToInt(t *testing.T) {
	assertMapReverseIntInt(t, map[int]int{1: 2, 3: 4}, false, map[int]int{2: 1, 4: 3}, nil)
}

func TestReverseMapIntToIntNilMap(t *testing.T) {
	assertMapReverseIntInt(t, nil, false, nil, errors.New("source should not be nil"))
}

func TestReverseMapIntToIntDuplicate(t *testing.T) {
	assertMapReverseIntInt(t, map[int]int{1: 2, 3: 2}, false, nil, errors.New("key should not be duplicated in result map"))
}

func TestReverseMapIntToIntDuplicateRewrite(t *testing.T) {
	assertMapReverseIntInt(t, map[int]int{1: 2, 3: 2}, true, map[int]int{2: 3}, nil)
}
