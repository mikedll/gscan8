package main

import (
	"errors"
	"os"
	"log"
	"time"
	"regexp"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
)

type User struct {
	Id           int64      `gorm:"PRIMARY_KEY;AUTO_INCREMENT"`
	Username     string     `gorm:"not null;unique_index"`
	AccessToken  string     `gorm:"not null"`
	TokenExpiry  time.Time
}

type GistFile struct {
	Id       int64  `json:"id"        gorm:"PRIMARY_KEY;AUTO_INCREMENT"`
	UserId   int64  `                 gorm:"not null"`
	VendorId string `json:"vendor_id" gorm:"not null"`
	Title    string `json:"title"     gorm:"default '';not null"`
	Filename string `json:"filename"  gorm:"not null"`
	Body     string `json:"body"      gorm:"type:character varying;default:''"`
	Language string `json:"language"  gorm:"not null"`
}

type Snippet struct {
	GistFileId int64  `json:"id"`
	VendorId   string `json:"vendor_id"`
	Filename   string `json:"filename"`
	LineNumber int    `json:"line_number"`
	Title      string `json:"title"`
	Body       string `json:"body"`
	Language   string `json:"language"`
}

var dbConn *gorm.DB

func openDbForProject(isProduction bool) {
	var err error
	var connString string
	
	if isProduction {
		connString = os.Getenv("DATABASE_URL")
	} else {
		connString = "host=localhost user=gscan8dev dbname=gscan8development password=thintent"
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

func findUserByLogin(login string, user *User) error {
	dbConn.Where("username = ?", login).First(user)
	if err := dbConn.Error; err != nil {
		return errors.New("createUser failed")
	}

	return nil
}

func searchGistFiles(userId int64, query string) (results []Snippet, err error) {
	results = []Snippet{}
	err = nil

	gistFiles := []GistFile{}
	dbConn.Where("user_id = ? AND body like ?", userId, "%" + query + "%").Find(&gistFiles)

	if err = dbConn.Error; err != nil {
		return
	}

	var queryRegex *regexp.Regexp
	queryRegex, err = regexp.Compile(query) // TODO replace special characters
	if err != nil {
		err = errors.New("Unable to compile regular expression: " + err.Error())
		return
	}
	
	for _, gistFile := range gistFiles {
		matches := queryRegex.FindAllStringIndex(gistFile.Body, -1)
		for _, match := range matches {
			min := match[0] - 100
			max := match[1] + 100
			if min < 0 {
				min = 0
			}
			if max > len(gistFile.Body) {
				max = len(gistFile.Body)
			}

			snippet := Snippet{
				GistFileId: gistFile.Id,
				VendorId: gistFile.VendorId,
				Filename: gistFile.Filename,
				LineNumber: 1, // TODO calculate line number
				Title: gistFile.Title,
				Body: gistFile.Body[min:max],
				Language: gistFile.Language,
			}

			results = append(results, snippet)
		}
	}

	return
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

func makeSchema() error {
	dbConn.AutoMigrate(&User{})
	dbConn.AutoMigrate(&GistFile{})

	return nil
}
