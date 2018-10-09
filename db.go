package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

type Gist struct {
  Id int        `json:"id"`
	Title string  `json:"title"`
	Url string   `json:"url"`
}

const dbPath string = "./storage/db.sqlite3"

// remoteGists returns fake data for seeding into a database.
func remoteGists() []Gist {
	fetched := []Gist{
		Gist{1, "purchase_orders.html", "https://gist.github.com/mikedll/8eaa6df25ac7a10ae3ded33e7f00b306"},
		Gist{2, "gist:8ea5f31a1269ed482f3ad0f7b274ee05", "https://gist.github.com/mikedll/8ea5f31a1269ed482f3ad0f7b274ee05"},
	}

	return fetched
}

func getGists() (collected []Gist) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, title, url FROM gists")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var cur Gist
		err = rows.Scan(&cur.Id, &cur.Title, &cur.Url)
		if err != nil {
			log.Fatal(err)
		}

		collected = append(collected, cur)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return
}

func makeGists() {
	os.Remove(dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
	CREATE TABLE gists (id INTEGER NOT NULL PRIMARY KEY, title TEXT, url TEXT);
	DELETE FROM gists;
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("INSERT INTO gists (id, title, url) VALUES(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	fetched := remoteGists()
	for _, f := range fetched {
		_, err = stmt.Exec(f.Id, f.Title, f.Url)
		if err != nil {
			log.Fatal(err)
		}
	}
	tx.Commit()
}
