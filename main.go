package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sandjuarezg/example-goroutines/migration"
)

func main() {
	var wg = sync.WaitGroup{}
	var mutex = sync.Mutex{}
	var people []string
	migration.SqlMigration()

	db, err := sql.Open("sqlite3", "./people.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	file, err := os.Open("./files/people.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var r = csv.NewReader(file)

	for {
		data, err := r.Read()
		if err != nil {
			break
		}

		for i := range data {
			people = append(people, data[i])
		}
	}

	for i := 0; i < len(people); i = i + 5 {
		wg.Add(1)
		insertData(i, db, people, &mutex, &wg)
	}

	/*go func() {
		for _ = range time.Tick(1 * time.Millisecond) {
			fmt.Println("Mira porque el mundo se destruye:", runtime.NumGoroutine())
		}
	}()*/

	wg.Wait()

	fmt.Println("Finish")
}

func insertData(i int, db *sql.DB, people []string, mutex *sync.Mutex, wg *sync.WaitGroup) {
	mutex.Lock()
	smt, err := db.Prepare("INSERT INTO people (nombre, apellidoP, apellidoM, genero, edad) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer smt.Close()

	_, err = smt.Exec(people[i], people[i+1], people[i+2], people[i+3], people[i+4])
	if err != nil {
		log.Fatal(err)
	}
	mutex.Unlock()
	wg.Done()
}
