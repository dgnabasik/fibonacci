package main

// database.go:: database driver: go get -u github.com/jackc/pgx

import (
	"bufio"
	"context" // pgx driver uses context: see https://golang.org/pkg/context/
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	pgx "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool" // https://pkg.go.dev/github.com/jackc/pgx/v4/pgxpool
)

const (
	DBHOST   = "database"
	DBPORT   = 5432
	DBSCHEMA = "" // public.
)

func IsProduction() bool {
	mode := os.Getenv("NODE_ENV")
	return strings.ToLower(mode) == "production"
}

// GetHost func returns full hostname from fib.env (do NOT include :port)
func GetHost() string {
	if IsProduction() {
		return "http://localhost"
	}

	host := os.Getenv("FIB_API_DOMAIN")
	if host == "" {
		host = "http://localhost"
	}
	return host
}

// GetPort func
func GetPort() string {
	if IsProduction() {
		return "8080"
	}
	return "5000"
}

// GetDatabaseConnectionString func uses environment var FIB_DATABASE_URL
func GetDatabaseConnectionString() string {
	if IsProduction() {
		username := os.Getenv("POSTGRES_USER")
		database := os.Getenv("POSTGRES_DB")
		password := os.Getenv("POSTGRES_PASSWORD")
		if username == "" || database == "" || password == "" {
			fmt.Println("POSTGRES variables not found in environment variables...exiting")
			os.Exit(1)
		}
		return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", DBHOST, DBPORT, username, password, database)
	}

	connStr := os.Getenv("FIB_DATABASE_URL")
	if connStr == "" {
		fmt.Println("FIB_DATABASE_URL not found in environment variables...exiting")
		os.Exit(1)
	}
	return connStr
}

// CheckErr database error handler.
func CheckErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

// GetDatabaseReference opens a database specified by its database driver name and a driver-specific data source name: db,err := GetDatabaseReference()
// defer db.Close() must follow a call to this function in the calling function. sslmode is set to 'required' by default.
func GetDatabaseReference() (*pgxpool.Pool, error) {
	dbConn := GetDatabaseConnectionString()
	fmt.Println("Connecting to: " + dbConn)
	db, err := pgxpool.Connect(context.Background(), dbConn)
	CheckErr(err)
	err = db.Ping(context.Background())
	if err != nil {
		fmt.Println(err)
	}
	return db, err
}

/**************************************************************************************/

// GetMemoizedResults func fetches all the memoized results less than a given value (e.g. there are 11 intermediate results less than 120),
func GetMemoizedResults(fibLimit float64) ([]FibonacciDB, error) {
	db, err := GetDatabaseReference()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	// Primary Key always returns in key order.
	query := "SELECT id, fibvalue FROM " + DBSCHEMA + "Fibonacci WHERE fibvalue < $1"
	rows, err := db.Query(context.Background(), query, fibLimit)
	CheckErr(err)
	defer rows.Close()

	fibList := make([]FibonacciDB, 0)
	var item FibonacciDB
	for rows.Next() {
		err := rows.Scan(&item.ID, &item.FibValue)
		CheckErr(err)
		fibList = append(fibList, item)
	}
	err = rows.Err()
	CheckErr(err)
	return fibList, err
}

// ClearDataStore func clears the data store but does not RESTART IDENTITY since no sequence defined.
func ClearDataStore() error {
	db, err := GetDatabaseReference()
	CheckErr(err)
	defer db.Close()
	_, err = db.Exec(context.Background(), "TRUNCATE TABLE "+DBSCHEMA+"Fibonacci")
	CheckErr(err)
	return err
}

// convertMapToFibonacciSlice in key order.
func convertMapToFibonacciSlice(bigmap map[int]float64) []FibonacciDB {
	// Convert map to slice of keys.
	keys := []int{}
	for key := range bigmap {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	// Convert map to slice of values.
	values := []float64{}
	for _, value := range bigmap {
		values = append(values, value)
	}
	sort.Float64s(values)

	bigFibSlice := make([]FibonacciDB, len(keys))
	for ndx := range keys {
		bigFibSlice[ndx] = FibonacciDB{ID: keys[ndx], FibValue: values[ndx]}
	}

	return bigFibSlice
}

// SetMemoizedResults func performs Bulk insert/append.
func SetMemoizedResults(bigmap map[int]float64) error {
	/*err := ClearDataStore()
	if err != nil {
		return err
	}*/

	bigFibSlice := convertMapToFibonacciSlice(bigmap)

	db, err := GetDatabaseReference()
	if err != nil {
		return err
	}
	defer db.Close()

	txn, err := db.Begin(context.Background())
	CheckErr(err)

	// Bulk insert must use lowercase column names!
	copyCount, err := db.CopyFrom(
		context.Background(),
		pgx.Identifier{DBSCHEMA + "fibonacci"}, // tablename
		[]string{"id", "fibvalue"},
		pgx.CopyFromSlice(len(bigFibSlice), func(i int) ([]interface{}, error) {
			return []interface{}{bigFibSlice[i].ID, bigFibSlice[i].FibValue}, nil
		}),
	)

	CheckErr(err)
	if copyCount == 0 {
		fmt.Println("SetMemoizedResults: no rows inserted")
	}
	err = txn.Commit(context.Background())
	CheckErr(err)

	return err
}

/**************************************************************************************/

// WriteTextFile func
func WriteTextFile(lines []string, filePath string) error {
	_ = os.Remove(filePath)
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	datawriter := bufio.NewWriter(file)
	for _, line := range lines {
		_, _ = datawriter.WriteString(line + "\n")
	}

	datawriter.Flush()
	file.Close()
	return nil
}

// WriteFibonacciToCsvFile func
func WriteFibonacciToCsvFile(data []FibonacciDB, outputFile string) error {
	fmt.Println("\nWriting performance data to " + outputFile)
	lines := make([]string, len(data)+1)
	lines[0] = "Iterations,Seconds"
	ndx := 1
	for _, fib := range data {
		lines[ndx] = fmt.Sprintf("%s,%g", strconv.Itoa(fib.ID), fib.FibValue)
		ndx++
	}
	return WriteTextFile(lines, outputFile)
}

// Performance func
func Performance() []FibonacciDB {
	f := fibonacci()
	var result float64
	bigFibSlice := make([]FibonacciDB, 0)

	for iterations := 32; iterations <= 1024; iterations *= 2 {
		startTime := time.Now()
		for iter := 0; iter < iterations; iter++ {
			result = f()
			if result > math.MaxFloat64-1 {
				fmt.Println("Numeric limit exceeded!")
				break
			}
		}
		elapsed := time.Since(startTime)
		bigFibSlice = append(bigFibSlice, FibonacciDB{ID: iterations, FibValue: elapsed.Seconds()})
	}

	return bigFibSlice
}
