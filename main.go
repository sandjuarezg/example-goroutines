package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sync"
)

func main() {
	var wg = sync.WaitGroup{}
	var mutex = sync.Mutex{}
	var people []string

	var file, err = os.Open("./files/people.csv")
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

		wg.Add(1)
		go func(m *sync.Mutex, w *sync.WaitGroup, people *[]string) {
			m.Lock()
			for i := range data {
				*people = append(*people, data[i])
			}
			m.Unlock()
			w.Done()
		}(&mutex, &wg, &people)

	}

	wg.Wait()

	fmt.Println(people)
}
