package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"sync"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sandjuarezg/exampleGoroutines/migration"
)

func main() {
	var wg = sync.WaitGroup{}
	var mutex = sync.Mutex{}
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

	for err != io.EOF {
		wg.Add(1)
		go func(m *sync.Mutex, w *sync.WaitGroup) {
			m.Lock()
			err = insertData(r, db)
			m.Unlock()
			w.Done()
		}(&mutex, &wg)
		fmt.Println(runtime.NumGoroutine())
	}

	wg.Wait()

	fmt.Println("Data added successfully")
}

func insertData(r *csv.Reader, db *sql.DB) (err error) {
	smt, err := db.Prepare("INSERT INTO people (nombre, apellidoP, apellidoM, genero, edad) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer smt.Close()

	data, err := r.Read()
	if err != nil {
		return
	}

	var people []string
	for i := range data {
		people = append(people, data[i])
	}
	_, err = smt.Exec(people[0], people[1], people[2], people[3], people[4])
	if err != nil {
		log.Fatal(err)
	}

	return
}
