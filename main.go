package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var wg = sync.WaitGroup{}
	var mutex = sync.Mutex{}

	var _, err = os.Stat("./database.sql")
	if os.IsNotExist(err) {
		log.Fatal(err)
	}

	var file, _ = os.Open("./database.sql")
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = os.Stat("./people.db")
	if os.IsNotExist(err) {
		var file, err = os.Create("./people.db")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	}

	db, err := sql.Open("sqlite3", "./people.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Query("SELECT * from people")
	if err != nil {
		_, err = db.Exec(string(content))
		if err != nil {
			log.Fatal(err)
		}
	}

	file, err = os.Open("./files/people.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var r = csv.NewReader(file)

	for err == nil {
		wg.Add(1)
		go func() {
			err = insertData(r, db, &mutex, &wg)
		}()
	}

	wg.Wait()

	fmt.Println("Data added successfully")
}

func insertData(r *csv.Reader, db *sql.DB, mutex *sync.Mutex, wg *sync.WaitGroup) (err error) {
	defer mutex.Unlock()
	defer wg.Done()
	mutex.Lock()

	smt, err := db.Prepare("INSERT INTO people (nombre, apellidoP, apellidoM, genero, edad) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer smt.Close()

	data, err := r.Read()
	if err != nil {
		if err == io.EOF {
			return
		}
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
