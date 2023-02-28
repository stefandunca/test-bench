package learning

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/tour/tree"
)

func TestWalkTree(t *testing.T) {
	testTree := tree.New(1)
	ch := make(chan int)
	go Walk(testTree, ch)
	for i := 1; i <= 10; i++ {
		require.Equal(t, i, <-ch)
	}
}

func TestSameTrees(t *testing.T) {
	testTree1 := tree.New(1)
	testTree2 := tree.New(1)
	require.True(t, Same(testTree1, testTree2))

	testTree3 := tree.New(2)
	require.False(t, Same(testTree1, testTree3))
}
