package main

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
//	"os"
	// "io/ioutil"
	"errors"
	// "fmt"
)

type GistFile struct {
  Id int64          `json:"id"        gorm:"PRIMARY_KEY;AUTO_INCREMENT"`
	VendorId string   `json:"vendor_id" gorm:"not null"`
	Title string      `json:"title"     gorm:"default '';not null"`
	Filename string   `json:"filename"  gorm:"not null"`
	Body string       `json:"body"      gorm:"type:character varying;default:''"`
	Language string   `json:"language"  gorm:"not null"`
}

type Snippet struct {
	Id int64          `json:"id"`
	Title int         `json:"title"`
	Body string       `json:"body"`
	Language string   `json:"language"`
}

const dbPath string = "./storage/db.sqlite3"

var dbConn *gorm.DB

func openDbForProject(isProduction bool) {
	var err error
	connString := "host=localhost user=gscan8dev dbname=gscan8development password=thintent"
	if !isProduction {
		 connString = connString + " sslmode=disable"
	}
	dbConn, err = gorm.Open("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}
}

func closeDbForProject() {
	dbConn.Close()	
}

// remoteGists returns fake data for seeding into a database.
func remoteGistFiles() []GistFile {
	fetched := []GistFile{
		GistFile{1, "8eaa6df25ac7a10ae3ded33e7f00b306", "purchase_orders.html", "blah.py", "some code", "Ruby"},
		GistFile{2, "8ea5f31a1269ed482f3ad0f7b274ee05", "app_world.txt", "blah.txt", "some text", "Ruby"},
	}

	return fetched
}

func searchGistFiles(query string) (results []Snippet) {
	results = []Snippet{}
	// search db, get back bodies

	// search again, get indices.

	// search backward/forward to discover nearby lines.

	// assemble snippets with languages
	return
}

//
// Seeds db with some fake data.
// 
func makeGistFiles() (error) {
	fetched := remoteGistFiles()
	for _, gistFile := range fetched {
		dbConn.Create(&gistFile)
		if err := dbConn.Error; err != nil {
			return errors.New("insert failed")
		}
	}
	
	return nil
}

func getGistFiles() (results []GistFile) {
	if err := dbConn.Find(&results).Error; err != nil {
		log.Fatal(err)
	}

	return
}

func makeSchema() (error) {
	dbConn.AutoMigrate(&GistFile{})

 	return nil
}

