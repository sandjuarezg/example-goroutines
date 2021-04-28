package migration

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"
)

func SqlMigration() {
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

	return

}
