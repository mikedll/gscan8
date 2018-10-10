package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"io/ioutil"
	"errors"
	"fmt"
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

func schemaString(isProduction bool) (sql string, error error) {
	var pth string
	
	if isProduction {
		pth = "config/spostgres.sql"
	} else {
		pth = "config/ssqlite.sql"
	}

	if _, err := os.Stat(pth); err == nil {
		bytes, err := ioutil.ReadFile(pth)
		if err != nil {
			error = errors.New("unable to open schema file")
			return
		}
		sql = string(bytes)
	}

	return
}

func getDb(isProduction bool) (db *sql.DB, err error){
	if isProduction {
		log.Fatal("postgres not implemented yet")
	} else {
		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			log.Fatal(err)
		}
		return db, err
	}

	return nil, errors.New("didn't get database connection")
}

func makeSchema(isProduction bool) (error) {
	if (!isProduction) {
		os.Remove(dbPath)
	}

	db, err := getDb(isProduction)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var schemaStmt string
	schemaStmt, err = schemaString(isProduction)
	if err != nil {
		fmt.Println("unable to open schema file", err)
		return err
	}
	
	_, err = db.Exec(schemaStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, schemaStmt)
		return err
	}

	return nil
}

func makeGists(isProduction bool) (error) {
	db, err := getDb(isProduction)
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("INSERT INTO gists (id, title, url) VALUES(?, ?, ?)")
	if err != nil {
		return errors.New("failed to prepare insert statement")
	}
	defer stmt.Close()

	fetched := remoteGists()
	for _, f := range fetched {
		_, err = stmt.Exec(f.Id, f.Title, f.Url)
		if err != nil {
			return errors.New("insert failed")
		}
	}
	tx.Commit()
	
	return err
}
