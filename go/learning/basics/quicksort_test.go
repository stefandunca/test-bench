package learning

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func quicksort(array []int, pivotSrc func(int) int) []int {
	if len(array) < 2 {
		return array
	}

	left, right := 0, len(array)-1

	pivot := pivotSrc(len(array))
	array[pivot], array[right] = array[right], array[pivot]

	for i := range array {
		if array[i] < array[right] {
			array[left], array[i] = array[i], array[left]
			left++
		}
	}

	array[left], array[right] = array[right], array[left]

	quicksort(array[:left], pivotSrc)
	quicksort(array[left+1:], pivotSrc)

	return array
}

// Benchmark for quicksort vs. quicksort1
func TestQuicksort(t *testing.T) {
	testCases := []struct {
		name     string
		in       []int
		pivotSrc func(int) int
		expected []int
	}{
		{
			name:     "unsorted",
			in:       []int{3, 4, 5, 2, 1},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "reversed",
			in:       []int{5, 4, 3, 2, 1},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "sorted",
			in:       []int{1, 2, 3, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
		{
			name:     "one item",
			in:       []int{6},
			expected: []int{6},
		},
		{
			name:     "empty",
			in:       []int{},
			expected: []int{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			in := make([]int, len(tc.in))
			copy(in, tc.in)
			out := quicksort(in, func(len int) int { return rand.Intn(len) })
			require.Equal(t, tc.expected, out, "quicksort failed for %v", tc.name)
		})
	}
}

func BenchmarkQuicksort(b *testing.B) {
	sortedDataSet := make([]int, 10000)
	reversedDataSet := make([]int, len(sortedDataSet))
	randomDataSet := make([]int, len(sortedDataSet))
	for i := 0; i < len(sortedDataSet); i++ {
		sortedDataSet[i] = i
		reversedDataSet[i] = len(reversedDataSet) - i
		randomDataSet[i] = rand.Int()
	}

	testCases := []struct {
		name string
		in   []int
	}{
		{
			name: "sorted",
			in:   sortedDataSet,
		},
		{
			name: "reversed",
			in:   reversedDataSet,
		},
		{
			name: "random",
			in:   randomDataSet,
		},
	}

	type pivotInfo struct {
		name     string
		pivotSrc func(int) int
	}

	pivots := []pivotInfo{
		{
			name:     "HalfPivot",
			pivotSrc: func(len int) int { return len / 2 },
		},
		{
			name:     "RandomPivot",
			pivotSrc: func(len int) int { return rand.Intn(len) },
		},
		{
			name:     "LastPivot",
			pivotSrc: func(len int) int { return len - 1 },
		},
		{
			name:     "FirstPivot",
			pivotSrc: func(len int) int { return 0 },
		},
	}

	for pI := 0; pI < len(pivots); pI++ {
		for _, tc := range testCases {
			b.Run(tc.name+pivots[pI].name, func(b *testing.B) {
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					b.StopTimer()
					in := make([]int, len(tc.in))
					copy(in, tc.in)
					b.StartTimer()
					quicksort(in, pivots[pI].pivotSrc)
				}
			})
		}
	}
}
