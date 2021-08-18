package database

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"sync"

	"github.com/brianvoe/gofakeit"
)

var mtx sync.Mutex
var toolData = "./db/data/tools.json"
var userData = "./db/data/users.json"
var rentalData = "./db/data/rentals.json"

// init generates seed data for DB.
func init() {
	if !fileExists(toolData) || !fileExists(userData) || !fileExists(rentalData) {
		err := createFiles(toolData, userData, rentalData)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println("Populating Database")

		gofakeit.Seed(0)
		generateTool := func(i int) *Tool {
			return &Tool{
				ID:       i + 1,
				Name:     gofakeit.BuzzWord() + " " + gofakeit.HackerNoun(),
				Desc:     gofakeit.HipsterSentence(5),
				Price:    gofakeit.Price(1, 2000),
				Quantity: gofakeit.Number(1, 30),
			}
		}

		generateUser := func(i int) *User {
			return &User{
				ID:    i + 1,
				Name:  gofakeit.Name(),
				Email: gofakeit.Email(),
			}
		}

		generateRental := func(i int) *Rental {
			rand.Seed(rand.Int63n(1000))
			return &Rental{
				ID:     i + 1,
				Active: true,
				ToolID: rand.Intn(10) + 1,
				UserID: rand.Intn(10) + 1,
			}
		}

		tools := []Tool{}
		users := []User{}
		rentals := []Rental{}

		for i := 0; i <= 10; i++ {
			tools = append(tools, *generateTool(i))
			users = append(users, *generateUser(i))
			rentals = append(rentals, *generateRental(i))
		}

		Save(toolData, tools)
		Save(userData, users)
		Save(rentalData, rentals)
	}
}

// fileExists determines if the specified data file exists.
func fileExists(file string) bool {
	info, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func createFiles(files ...string) error {
	for i := 0; i < len(files); i++ {
		_, err := os.Create(files[i])
		if err != nil {
			return errors.New("We've encountered an error, please ensure that the ./db/data directory exists")
		}
	}
	return nil
}

// toJSON takes in structured data and formats appropriately.
// Upon success an a new io.Reader is returned
func toJSON(v interface{}) (io.Reader, error) {
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

// Save handles persisting data to in-memory database.
// An error is returned if data cannot be added or copied.
func Save(path string, x interface{}) error {
	mtx.Lock()
	defer mtx.Unlock()

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := toJSON(x)
	if err != nil {
		return err
	}

	_, err = io.Copy(file, data)
	return err
}

// Load reads data found in the in-memory database.
// Data is returned as a slice of bytes.
// An error is returned ReadFile fails.
func Load(path string) ([]byte, error) {
	mtx.Lock()
	defer mtx.Unlock()

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}
