package main

import (
	"errors"
	"os"
	"log"
	"time"
	"regexp"
	"strings"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"
)

func min(a, b int) int {
	if(a < b) {
		return a
	}
	return b
}

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
		connString = "host=localhost user=postgres dbname=gscan8_development"
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
	if query == "" {
		return
	}
	
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

		matchIndex := 0
		matchLines := []int{}
		blocks := [][]int{}
		lines := []string{}
		curLineIndex := 0
		curLineMin := 0

		// store all lines, and store line indices surrounding matches
		for i, c := range gistFile.Body {
			if c == '\n' {
				lines = append(lines, gistFile.Body[curLineMin:i])
				
				for matchIndex < len(matches) && curLineMin <= matches[matchIndex][0] && matches[matchIndex][1] <= i {
					matchLines = append(matchLines, curLineIndex)
					blockIndices := []int{}
					for j := -2; j < 3; j++ {
						blockIndices = append(blockIndices, curLineIndex + j)
					}
					blocks = append(blocks, blockIndices)

					matchIndex += 1
				}

				curLineMin = i + 1
				curLineIndex += 1
			}
		}

		// len(blocks) = len(matches), assuming \n is not in the pattern. fix above.

		for i, _ := range matches {

			blocksAsStrings := []string{}
			for j := 0; j < min(5, len(blocks[i])); j++ {
				lineIndex := blocks[i][j]
				if lineIndex >= 0 && lineIndex < len(lines) {
					blocksAsStrings = append(blocksAsStrings, lines[blocks[i][j]])
				}
			}

			snippet := Snippet{
				GistFileId: gistFile.Id,
				VendorId: gistFile.VendorId,
				Filename: gistFile.Filename,
				LineNumber: blocks[i][2] + 1,
				Title: gistFile.Title,
				Body: strings.Join(blocksAsStrings, "\n"),
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
