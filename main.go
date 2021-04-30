//Aqui no es necesario el mutex
//Debido a que cada goroutine se crea en cada iteración
//Por lo tanto, cada uno tiene un valor distinto

package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sandjuarezg/example-goroutines/migration"
)

func main() {
	var wg = sync.WaitGroup{}
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
		go insertData(i, db, people, &wg)
	}

	go func() {
		for _ = range time.Tick(300 * time.Millisecond) {
			fmt.Println("Mira porque el mundo se destruye:", runtime.NumGoroutine())
		}
	}()

	wg.Wait()

	fmt.Println("Finish")
}

func insertData(i int, db *sql.DB, people []string, wg *sync.WaitGroup) {
	smt, err := db.Prepare("INSERT INTO people (nombre, apellidoP, apellidoM, genero, edad) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer smt.Close()

	_, err = smt.Exec(people[i], people[i+1], people[i+2], people[i+3], people[i+4])
	if err != nil {
		log.Fatal(err)
	}
	wg.Done()
}
