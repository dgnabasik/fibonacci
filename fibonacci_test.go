package main

// All not-test files init() functions are executed first, then all test files init() functions are executed (hopefully in lexical order).
import (
	"fmt"
	"testing"
)

// Test_database func
func Test_database(t *testing.T) {
	if IsProduction() && GetHost() != "http://server" {
		t.Error("a1a:GetHost incorrect ")
	}
	if !IsProduction() && GetHost() != "http://localhost" {
		t.Error("a1b:GetHost incorrect ")
	}

	if IsProduction() && GetPort() != "8080" {
		t.Error("a2a:GetPort incorrect ")
	}
	if !IsProduction() && GetPort() != "5000" {
		t.Error("a2b:GetPort incorrect ")
	}

	dbConn := GetDatabaseConnectionString()
	if dbConn == "" {
		t.Error("a3:GetDatabaseConnectionString produced nothing.")
	}

	db, err := GetDatabaseReference()
	CheckErr(err)
	defer db.Close()
	if err != nil {
		t.Error("a4:GetDatabaseReference produced error ", err)
	}

	// NoRowsReturned(err error) bool {

	err = ClearDataStore()
	if err != nil {
		t.Error("a5:ClearDataStore produced error ", err)
	}

	var bigmap = map[int]float64{0: 1, 1: 1, 2: 2, 3: 3, 4: 5, 5: 8, 6: 13, 7: 21, 8: 34, 9: 55, 10: 89, 11: 144}
	err = SetMemoizedResults(bigmap)
	if err != nil {
		t.Error("a6:SetMemoizedResults produced error ", err)
	}

	var fibLimit float64 = 120 // expect 11 results
	fibList, err := GetMemoizedResults(fibLimit)
	if err != nil {
		t.Error("a7a:GetMemoizedResults produced error ", err)
	}
	if fibList[len(fibList)-1].ID < 10 {
		t.Error("a7b:GetMemoizedResults does not have enough data for fibLimit of ", fibLimit)
	}
	if len(fibList) != 11 {
		t.Error("a7c:GetMemoizedResults expected 11 but produced ", len(fibList))
	}
	fmt.Printf("%s %f %f", "Actual value < ", fibLimit, fibList[len(fibList)-1].FibValue)

	outputFile := "/tmp/fibonacci.csv"
	lines := []string{"line1", "line2"}
	err = WriteTextFile(lines, outputFile)
	if err != nil {
		t.Error("a8:WriteTextFile produced error ", err)
	}

	data := make([]FibonacciDB, 0)
	data = append(data, FibonacciDB{ID: 0, FibValue: 4})
	err = WriteFibonacciToCsvFile(data, outputFile)
	if err != nil {
		t.Error("a9:WriteFibonacciToCsvFile produced error ", err)
	}

	fmt.Println()
}

// Test_fibonacci func
func Test_fibonacci(t *testing.T) {
	iterations := 12
	bigmap := make(map[int]float64)
	f := fibonacci()
	var result float64
	for iter := 0; iter < iterations; iter++ {
		result = f()
		bigmap[iter] = result
	}
	if len(bigmap) < iterations {
		t.Error("b1:fibonacci expected iterations of ", iterations, " but got ", len(bigmap))
	}
	fmt.Printf("%s%d%s", "Fib(", iterations, ")=")
	fmt.Println(result)

	FibonacciDBslice := Performance()
	err := WriteFibonacciToCsvFile(FibonacciDBslice, "/tmp/fibonacci.csv")
	if err != nil {
		t.Error("b2:Failed Performance(): ", err)
	}

}
