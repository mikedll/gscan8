package main

import (
	"errors"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
)

type GistFile struct {
	Id       int64  `json:"id"        gorm:"PRIMARY_KEY;AUTO_INCREMENT"`
	VendorId string `json:"vendor_id" gorm:"not null"`
	Title    string `json:"title"     gorm:"default '';not null"`
	Filename string `json:"filename"  gorm:"not null"`
	Body     string `json:"body"      gorm:"type:character varying;default:''"`
	Language string `json:"language"  gorm:"not null"`
}

type Snippet struct {
	Id       int64  `json:"id"`
	Title    int    `json:"title"`
	Body     string `json:"body"`
	Language string `json:"language"`
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
func mockGistFiles() ([]GistFile, error) {
	fetched := []GistFile{
		GistFile{1, "8eaa6df25ac7a10ae3ded33e7f00b306", "purchase_orders.html", "DbWrapper.cs", "", "CSharp"},
		GistFile{2, "7051b354ac1f5705587386a4cf07efe8", "pos.sql", "pos.sql", "", "SQL"},
	}

	var err error
	var bytes []byte
	for i, gistFile := range fetched {
		filePath := "mockdata/" + strconv.FormatInt(gistFile.Id, 10) + ".txt"
		bytes, err = ioutil.ReadFile(filePath)
		if err != nil {
			return []GistFile{}, err
		}

		fetched[i].Body = string(bytes)
	}

	return fetched, nil
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
func makeGistFiles() error {
	fetched, err := mockGistFiles()
	if err != nil {
		return errors.New("error while retrieving mock gist files")
	}
	for _, gistFile := range fetched {
		dbConn.Create(&gistFile)
		if err := dbConn.Error; err != nil {
			return errors.New("insert failed")
		}
	}

	return nil
}

//
// Deletes all data in db.
//
func emptyDb() error {
	dbConn.Delete(GistFile{})
	err := dbConn.Error
	if err != nil {
		return errors.New("Error while deleting db data.")
	}

	return err
}

func getGistFiles() (results []GistFile) {
	if err := dbConn.Find(&results).Error; err != nil {
		log.Fatal(err)
	}

	return
}

func makeSchema() error {
	dbConn.AutoMigrate(&GistFile{})

	return nil
}
