package learning

import (
	"fmt"
	"image"
	"image/color"
	"strings"
)

func WordCount(s string) map[string]int {
	var counted = make(map[string]int)
	for _, w := range strings.Split(s, " ") {
		counted[w]++
	}
	return counted
}

func doWordCount() {
	fmt.Println(WordCount("First Second Second Again Third Forth Forth Again"))
}

// fibonacci is a function that returns
// a function that returns an int.
func fibonacci() func() int {
	prev, current := 0, 1
	return func() int {
		prevPrev := prev
		prev = current
		current = prevPrev + prev
		return prevPrev
	}
}

func doFibonacci() {
	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}

type IPAddr [4]byte

func (t IPAddr) String() string {
	var formatted string
	const countWithoutLast int = len(t) - 1
	for i := 0; i < countWithoutLast; i++ {

		formatted += fmt.Sprintf("%d.", t[i])
	}
	formatted += fmt.Sprintf("%d", t[countWithoutLast])
	return formatted
}

type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string {
	return fmt.Sprintf("Cannot Sqrt negative number: %v", float64(e))
}

func Sqrt(x float64) (float64, error) {
	if x < 0 {
		return 0, ErrNegativeSqrt(x)
	}
	res := x
	for i := 0; i < 100; i++ {
		prevRes := res
		res -= (res*res - x) / (2 * res)
		if (prevRes - res) < 0.0001 {
			break
		}
	}
	return res, nil
}

func doErrors() {
	fmt.Println(Sqrt(20008978900))
	fmt.Println(Sqrt(9))
	fmt.Println(Sqrt(-2))
}

func doStringify() {
	hosts := map[string]IPAddr{
		"loopback":  {127, 0, 0, 1},
		"googleDNS": {8, 8, 8, 8},
	}
	for name, ip := range hosts {
		fmt.Printf("%v: %v\n", name, ip)
	}
}

type MyReader struct{}

func (t MyReader) Read(buff []byte) (int, error) {
	for i := 0; i < len(buff); i++ {
		buff[i] = 'A'
	}
	return len(buff), nil
}

func doStreams() {
	for {
		buff := make([]byte, 8)
		_, err := MyReader{}.Read(buff)
		if err != nil {
			break
		}
		fmt.Print(string(buff))
	}
}

type Image struct {
	w int
	h int
}

func (t Image) ColorModel() color.Model {
	return color.RGBAModel
}

func (t Image) Bounds() image.Rectangle {
	return image.Rect(0, 0, t.w, t.h)
}

func (t Image) At(x, y int) color.Color {
	return color.RGBA{uint8(x * y), uint8(x * y), 255, 255}
}

func doImage() {
	m := Image{50, 50}
	fmt.Println(m.At(10, 10))
}

func exercise() {
	//doWordCount()
	//doFibonacci()
	//doStringify()
	//doErrors()
	//doStreams()
	//doStreams()
	doImage()
}
