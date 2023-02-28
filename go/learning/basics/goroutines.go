package learning

import (
	"fmt"
	"time"

	"golang.org/x/tour/tree"
)

// Duration
const oneTickMs = 10

func logTime(c chan rune) {
	for i := 0; i < 150; i++ {
		time.Sleep(oneTickMs * time.Millisecond)
		c <- '-'
	}
	close(c)
}

func emitCharEvery(char rune, durationMs time.Duration, c chan rune) {
	for {
		time.Sleep(durationMs * time.Millisecond)
		c <- char
	}
}

func concurrentEmit() {
	c := make(chan rune)
	go emitCharEvery('5', 50, c)
	go emitCharEvery('7', 70, c)
	go emitCharEvery('9', 90, c)
	go logTime(c)

	resultString := ""
	for char := range c {
		resultString += string(char)
	}
	fmt.Println(">>>", resultString)
}

//=============================================================================

// type Tree struct {
//     Left  *Tree
//     Value int
//     Right *Tree
// }

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func doWalk(t *tree.Tree, ch chan int) {
	if t == nil {
		return
	}
	doWalk(t.Left, ch)
	ch <- t.Value
	doWalk(t.Right, ch)
}

func Walk(t *tree.Tree, ch chan int) {
	doWalk(t, ch)
	close(ch)
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go Walk(t1, ch1)
	go Walk(t2, ch2)

	v1, ok1 := <-ch1
	v2, ok2 := <-ch2
	if ok1 && ok2 {
		if v1 != v2 {
			return false
		}
	} else {
		return false
	}
	return true
}

func walkTree() {

}
