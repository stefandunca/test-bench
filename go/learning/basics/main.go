package learning

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

func doPointers() {
	varInt := 5
	varIntPointer := &varInt
	fmt.Println("varInt:", varInt, "varIntPointer:", varIntPointer)

	stringPtr := new(string)
	*stringPtr = "Hello, pointers!"
	fmt.Println("stringPtr value:", *stringPtr, "stringPtr address:", stringPtr)
}

func allocateMemory() *int {
	return new(int)
}

func localMemory() *int {
	var localInt int = 12
	return &localInt
}

func doMemory() {
	allocatedMemory := allocateMemory()
	fmt.Println("allocatedMemory:", *allocatedMemory, "localMemory", *localMemory())
}

func doRunes() {
	emoji := rune('😇')
	fmt.Printf("emoji: %d; type: %T; reflect %s \n", emoji, emoji, reflect.TypeOf(emoji))
}

type Weekday int

const (
	Sunday Weekday = iota
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

func doIota() {
	today := Wednesday
	fmt.Printf("Sunday: %#v; today: %#v (%T)\n", Sunday, today, today)
}

func doSwitches() {
	today := Wednesday
	switch today {
	case Monday:
		fmt.Println("Monday")
	case Tuesday:
		fmt.Println("Tuesday")
	case Wednesday:
		fmt.Println("Wednesday")
	default:
		fmt.Println("Other day")
	}

	var i int = 4
	for i < 7 {
		switch {
		case i < 5:
			fmt.Println("i < 5")
		case i > 5:
			fmt.Println("i > 5")
		default:
			fmt.Println("i must be 5")
		}
		i++
	}
}

func doCasting() {
	var int16Var int16 = 5
	var int32Var int32 = int32(int16Var)
	fmt.Println("int16Var:", int16Var, "int32Var:", int32Var)
	int32Var = 10
	int16Var = int16(int32Var)
	fmt.Println("int16Var:", int16Var, "int32Var:", int32Var)
}

func doContainers() {
	intArray := []int{1, 2, 3, 4, 5}
	fmt.Println("int array:", intArray, "len:", len(intArray), "cap:", cap(intArray), "type:", reflect.TypeOf(intArray))
	intSlice := make([]int, 5)
	fmt.Println("int slice:", intSlice, "len:", len(intSlice), "cap:", cap(intSlice), "type:", reflect.TypeOf(intSlice))
	var testMap = map[string]int{"one": 1, "two": 2}
	fmt.Println("testMap:", testMap, "len:", len(testMap), "type:", reflect.TypeOf(testMap))
	_, exists := testMap["three"]
	fmt.Println("\"three\" exists in map:", exists)
	val, exists := testMap["two"]
	if exists {
		fmt.Print("\"", val, "\" aka \"two\" exists in map \n")
	} else {
		fmt.Println("Impossible")
	}
}

type TestStruct struct {
	strVar    string         `json:"str_var"`
	intStrMap map[int]string `json:"int_str_map", omitempty`
}

func doStructs() {
	testStruct := TestStruct{
		strVar:    "String member of struct",
		intStrMap: map[int]string{1: "one", 2: "two"},
	}
	testStruct.strVar += " " + reflect.TypeOf(testStruct).String()
	fmt.Println("testStruct:", testStruct)
	jsonBytes, err := json.Marshal(testStruct)
	if err != nil {
		panic(err)
	}
	fmt.Print("testStruct json \"", string(jsonBytes), "\"\n")

	type Vertex struct {
		Lat, Long float64
	}

	var m map[string]*Vertex

	verts := []Vertex{{
		40.68433, -74.39967,
	}}
	m = make(map[string]*Vertex)
	m["Bell Labs"] = &verts[0]
	fmt.Println("*m[\"Bell Labs\"]", *m["Bell Labs"], "TypeOf m[X]", reflect.TypeOf(m["Bell Labs"]))
}

func Pic(dx, dy int) [][]uint8 {
	const draw = 1 // 0 - 2
	res := make([][]uint8, dy)
	for y := 0; y < dy; y++ {
		res[y] = make([]uint8, dx)
		for x := 0; x < dx; x++ {
			var pixel uint8
			switch {
			case draw == 0:
				pixel = uint8((x + y) / 2)
			case draw == 1:
				pixel = uint8(x * y)
			case draw == 2:
				pixel = uint8(x ^ y)
			default:
				panic("Bad choice for draw")
			}
			res[y][x] = uint8(pixel)/uint8(256/('~'-'!')) + uint8('!')
		}
	}
	return res
}

func doPlot() {
	for _, row := range Pic(100, 50) {
		fmt.Println(string(row))
	}
}

func returnClosureBoundToLocalVar() func(string) string {
	localVar := "Initial value;"
	return func(addString string) string {
		localVar += addString
		return localVar
	}
}

func doClosures() {
	closure := returnClosureBoundToLocalVar()
	firstCallResult := closure("First call;")
	secondCallResult := closure("Second call;")
	fmt.Println("firstCallResult:", firstCallResult, "secondCallResult:", secondCallResult)
}

func doBatchIteration() {
	items := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22}
	maxItemsPerBatch := 5
	for i := 0; i < len(items); i += maxItemsPerBatch {
		last := i + maxItemsPerBatch
		if last > len(items) {
			last = len(items)
		}
		fmt.Println("[", i, "] ", items[i:last])
	}
}

type Learned func()

func main() {
	funcs := [...]Learned{doPointers, doMemory, doCasting, doRunes, doIota,
		doSwitches, doContainers, doStructs, doPlot, doClosures, doBatchIteration, concurrentEmit, walkTree}
	exit := false
	for !exit {
		choice := 0
		if len(os.Args) > 1 {
			choice, _ = strconv.Atoi(os.Args[1])
			exit = true
		} else {
			for i, f := range funcs {
				allName := strings.Split(runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(), ".")
				fmt.Println(i, allName[len(allName)-1])
			}
			fmt.Print("Choose function to run or ", len(funcs), " to exit: ")

			_, err := fmt.Scan(&choice)
			if err != nil {
				fmt.Println("\nInvalid input; error", err)
				break
			}
		}

		switch {
		case choice < len(funcs):
			fmt.Println("")
			funcs[choice]()
			fmt.Println("")
		case choice == len(funcs):
			exit = true
		default:
			fmt.Println("\nInvalid input", choice)
		}
	}
}
