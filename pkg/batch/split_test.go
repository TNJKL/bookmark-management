package batch_test

import (
	"testing"

	"github.com/TNJKL/bookmark-management/pkg/batch"
	"github.com/stretchr/testify/assert"
)

func TestSplitIntoBatches(t *testing.T) {
	t.Parallel()

	t.Run("split integers slice into batches", func(t *testing.T) {
		t.Parallel()

		testCases := []struct {
			name     string
			objects  []int
			size     int
			expected [][]int
		}{
			{
				name:     "slice smaller than batch size",
				objects:  []int{1, 2, 3},
				size:     5,
				expected: [][]int{{1, 2, 3}},
			},
			{
				name:     "slice length exact multiple of batch size",
				objects:  []int{1, 2, 3, 4},
				size:     2,
				expected: [][]int{{1, 2}, {3, 4}},
			},
			{
				name:     "slice length not exact multiple of batch size",
				objects:  []int{1, 2, 3, 4, 5},
				size:     2,
				expected: [][]int{{1, 2}, {3, 4}, {5}},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				result := batch.SplitIntoBatches(tc.objects, tc.size)
				assert.Equal(t, tc.expected, result)
			})
		}
	})

	t.Run("split strings slice into batches", func(t *testing.T) {
		t.Parallel()

		input := []string{"apple", "banana", "cherry", "date"}
		expected := [][]string{{"apple", "banana"}, {"cherry", "date"}}

		result := batch.SplitIntoBatches(input, 2)
		assert.Equal(t, expected, result)
	})
}
