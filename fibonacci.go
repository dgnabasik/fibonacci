package main

// fibonacci.go  Did not break out components into separate packages.
import (
	"fmt"
	"math"
	"os"
	"strconv"
	"time"
)

type FibonacciDB struct {
	ID       int     `json:"id"`
	FibValue float64 `json:"fibvalue"`
}

// fibonacci closure functor returns float64.
func fibonacci() func() float64 {
	var x float64 = 0
	var y float64 = 1
	return func() float64 {
		x, y = y, x+y
		return x
	}
}

func main() {
	/*FibonacciDBslice := Performance()
	err := WriteFibonacciToCsvFile(FibonacciDBslice, "/tmp/fibonacci.csv")
	if err != nil {
		fmt.Println(err)
	}*/

	iterations := 500
	if len(os.Args) > 1 {
		xiterations, err := strconv.Atoi(os.Args[1])
		if err != nil || xiterations < 0 { // math.MaxFloat64 = 1.7977+308 // 2**1023 * (2**53 - 1) / 2**52
			iterations = 500
			fmt.Printf("%s %d\n", "resetting iterations: ", iterations)
		} else {
			iterations = xiterations
		}
	}
	fmt.Printf("%s %d\n", "fibonacci iterations: ", iterations)

	bigmap := make(map[int]float64)
	startTime := time.Now()
	f := fibonacci()
	var result float64 // must be declared float64 for sake of f(closure).
	for iter := 0; iter < iterations; iter++ {
		result = f()
		bigmap[iter] = result
		if result > math.MaxFloat64/2 {
			break
		}
		fmt.Print(result)
		fmt.Print("  ")
	}
	elapsed := time.Since(startTime)
	fmt.Println("\nelapsed time: " + elapsed.String())

	err := SetMemoizedResults(bigmap)
	if err != nil {
		fmt.Println(err)
	}
}
