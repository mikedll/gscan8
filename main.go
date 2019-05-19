package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"io/ioutil"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/qor/render"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"github.com/gorilla/sessions"
	"github.com/gorilla/securecookie"
)

type UserApiResponse struct {
	Login       string  `json:"login"`
}

type FileApiResponse struct {
	Filename   string `json:"filename"`
	Language   string `json:"language"`
	RawUrl     string `json:"raw_url"`
}

type GistApiResponse struct {
	Id     string                     `json:"id"`
	Title  string                     `json:"description"`
	Files  map[string]FileApiResponse `json:"files"`
}

var isProduction bool

func stateStr() (string, error) {
	c := 30
	b := make([]byte, c)
	_, err := rand.Read(b)
	if err != nil {
		return "", errors.New("couldn't make random str: " + err.Error())
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func main() {
	isProduction = false

	flag.Parse()

	var addr string = ":8081"
	port := os.Getenv("PORT")
	appEnv := os.Getenv("APP_ENV")

	// Going to use this to determine production environment...LOL!
	if appEnv == "production" || port != "" {
		isProduction = true
		addr = fmt.Sprintf(":%s", port)
	}

	openDbForProject(isProduction)
	defer closeDbForProject()

	if flag.NArg() > 0 {
		if flag.Arg(0) == "schema" {
			err := makeSchema()
			if err != nil {
				log.Println("failed to create schema.")
				return
			}

			log.Println("created schema.")
		} else if flag.Arg(0) == "empty" {
			err := emptyDb()
			if err != nil {
				log.Println("failed to empty db", err)
				return
			}

			log.Println("emptied database.")
		} else if flag.Arg(0) == "keygen" {
			key := securecookie.GenerateRandomKey(32)
			if key == nil {
				log.Println("failed to generate random key")
				return
			}

			log.Println("key: " + base64.StdEncoding.EncodeToString(key))
		}
		return
	}

	Render := render.New(&render.Config{
		ViewPaths:     []string{},
		DefaultLayout: "",
		FuncMapMaker:  nil,
	})

	keyBytes, decodeErr := base64.StdEncoding.DecodeString(os.Getenv("SESSION_KEY"))
	if decodeErr != nil {
		log.Println("Unable to decode session key.")
		return
	}
	var sessionStore = sessions.NewCookieStore(keyBytes)
	sessionName := "gscan8session"

	logout := func(w http.ResponseWriter, req *http.Request) {
 		session, err := sessionStore.Get(req, sessionName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		
		delete(session.Values, "userId")
		err = session.Save(req, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return			
		}
		
		http.Redirect(w, req, "/", http.StatusFound)
	}
	
	root := func(w http.ResponseWriter, req *http.Request) {

 		session, err := sessionStore.Get(req, sessionName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user := User{}
		dbConn.Where("id = ?", session.Values["userId"]).First(&user)
		if dbConn.Error != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := make(map[string]interface{})
		if !dbConn.NewRecord(user) {
			ctx["username"] = template.HTML(user.Username)
			ctx["loggedIn"] = true
		} else {
			ctx["username"] = template.HTML("(not logged in)")
			ctx["loggedIn"] = false
		}
		
		Render.Execute("index", ctx, req, w)
	}

	search := func(w http.ResponseWriter, req *http.Request) {
		// search db for json
		snippets := searchGistFiles(req.URL.Query().Get("q"))
		snippetsJson, err := json.Marshal(snippets)
		if err != nil {
			log.Println("error while marshalling snippets: ", err)
			snippetsJson = []byte{}
			// StatusInternalServerError
			// write "Error while marshalling snippets.
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(snippetsJson)
		return
	}

	oauth2Conf := &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		Scopes:       []string{},
		Endpoint: github.Endpoint,
	}
	StateCookieName := "OAuth2-Github-State"
	
	oauth2Github := func(w http.ResponseWriter, req *http.Request) {
		stateStr, err := stateStr()
		if err != nil {
			log.Println("Unable to generate state string for oauth redirect.")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		stateCookie := http.Cookie{Name: StateCookieName, Value: stateStr, MaxAge: 60 * 15}
		http.SetCookie(w, &stateCookie)

		w.Header().Add("Content-Type", "text/html")

		url := oauth2Conf.AuthCodeURL(stateStr, oauth2.AccessTypeOffline)
		log.Println("Redirecting to ", url)
		http.Redirect(w, req, url, http.StatusFound)
	}

  oauth2GithubCallback := func(w http.ResponseWriter, req *http.Request) {
		writeError := func(msg string, errorNum int) {
			http.Error(w, msg, errorNum)
		}

		writeInteralServerError := func(msg string) {
			writeError(msg, http.StatusInternalServerError)
		}
		
		cookies := req.Cookies()

		var stateCookieVal string
		for _, cookie := range cookies {
			if cookie.Name == StateCookieName {
				stateCookieVal = cookie.Value
			}
		}		
		code := req.URL.Query().Get("code")
		stateInUrl := req.URL.Query().Get("state")

		if(stateInUrl == stateCookieVal) {
			token, err := oauth2Conf.Exchange(oauth2.NoContext, code)
			if err != nil {
				writeInteralServerError("Could not get access token from code=" + code + ", error: " + err.Error())
				return
			}

			// Access token obtained.
			client := oauth2Conf.Client(oauth2.NoContext, token)

			userUrl := "https://api.github.com/user"
			response, err := client.Get(userUrl)
			if err != nil {
				writeInteralServerError("Error fetching user information at " + userUrl)
				return
			}

			// UserApiResponse
			var responseBody []byte
			responseBody, err = ioutil.ReadAll(response.Body)
			
			if err != nil {
				writeInteralServerError(err.Error())
				return
			}

			userApiResponse := UserApiResponse{}
			json.Unmarshal(responseBody, &userApiResponse)
			
			userQueried := User{}
			err = findUserByLogin(userApiResponse.Login, &userQueried)
			if err != nil {
				writeInteralServerError("Unable to query for user: " + err.Error())
				return
			}

			userQueried.AccessToken = token.AccessToken
			userQueried.TokenExpiry = token.Expiry
			
			if(dbConn.NewRecord(userQueried)) {
				userQueried.Username = userApiResponse.Login

				if userQueried.Username == "" {
					writeError("Got empty username", http.StatusBadRequest)
					return
				}
				
				dbConn.Create(&userQueried)
				if err = dbConn.Error; err != nil {
					writeInteralServerError("Unable to create user: " + err.Error())
					return
				}
			} else {
				dbConn.Save(&userQueried)
				if err = dbConn.Error; err != nil {
					writeInteralServerError("Unable to save user: " + err.Error())
				}
			}

			session, err := sessionStore.Get(req, sessionName)
			if err != nil {
				writeInteralServerError(err.Error())
				return
			}
			
			session.Values["userId"] = userQueried.Id
			session.Save(req, w)
			if err != nil {
				writeInteralServerError(err.Error())
				return			
			}

			http.Redirect(w, req, "/", http.StatusFound)
		} else {
			writeError("OAuth2 state variables did not match.", http.StatusBadRequest)
		}
		log.Println("Received callback.")
	}

	fetchAllGists := func(w http.ResponseWriter, req * http.Request) {
		session, err := sessionStore.Get(req, sessionName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, ok := session.Values["userId"]; !ok {
			http.Error(w, "", http.StatusForbidden)
			return
		}
		
		user := User{}
		dbConn.Where("id = ?", session.Values["userId"]).First(&user)
		if dbConn.Error != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		token := oauth2.Token{AccessToken: user.AccessToken}
		client := oauth2Conf.Client(oauth2.NoContext, &token)
		response, err := client.Get("https://api.github.com/users/" + user.Username + "/gists")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var responseBody []byte
		responseBody, err = ioutil.ReadAll(response.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		gistsResponse := []GistApiResponse{}
		json.Unmarshal(responseBody, &gistsResponse)

		// For every file in every gist
		for _, gist := range gistsResponse {
			for _, fileInfo := range gist.Files {
				
				newGist := GistFile{
					UserId: user.Id,
					VendorId: gist.Id,
					Title: gist.Title,
					Filename: fileInfo.Filename,
					Language: fileInfo.Language,
				}

				response, err = http.Get(fileInfo.RawUrl)
				if err != nil {
					log.Println("Unable to fetch url " + fileInfo.RawUrl)
					continue
				}

				var bodyBytes []byte
				bodyBytes, err = ioutil.ReadAll(response.Body)
				if err != nil {
					log.Println("Unable to read response from url: " + fileInfo.RawUrl)
					continue
				}

				existingGistFile := GistFile{}
				
				gistQuery := GistFile{
					UserId: newGist.UserId,
					VendorId: newGist.VendorId,
					Filename: newGist.Filename,
				}

				dbConn.Where(&gistQuery).First(&existingGistFile)

				if dbConn.NewRecord(existingGistFile) {
					newGist.Body = string(bodyBytes)
					dbConn.Create(&newGist)
					if err = dbConn.Error; err != nil {
						log.Println("Unable to create gist file: " + newGist.VendorId + "/" + newGist.Filename)
					}
				} else {
					existingGistFile.Title = newGist.Title
					existingGistFile.Body = string(bodyBytes)
					existingGistFile.Language = newGist.Language
					dbConn.Save(&existingGistFile)
					if err = dbConn.Error; err != nil {
						log.Println("Unable to save gist file: " + newGist.VendorId + "/" + newGist.Filename)
					}					
				}
				
			}
		}
		w.Header().Add("Content-Type", "text/html")
		w.Write([]byte("I think I found a response"))
		// w.Write([]byte(responseBody))
	}

	getGists := func(w http.ResponseWriter, req *http.Request) {
		session, err := sessionStore.Get(req, sessionName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if _, ok := session.Values["userId"]; !ok {
			http.Error(w, "", http.StatusForbidden)
			return
		}
		
		gistFiles := []GistFile{}
		dbConn.Where(GistFile{UserId: session.Values["userId"].(int64)}).Find(&gistFiles)
		if dbConn.Error != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var gistFilesJson []byte
		gistFilesJson, err = json.Marshal(gistFiles)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(gistFilesJson)
	}
	
	defaultHandler := func(w http.ResponseWriter, req *http.Request) {
		file := "index.jsx"
		err := sendFile(file, w)
		if err != nil {
			log.Println("Unable to serve static file: ", file)
		}
	}

	log.Println("Starting server...")
	http.Handle("/", http.HandlerFunc(root))
	http.Handle("/logout", http.HandlerFunc(logout))

	http.Handle("/oauth/github", http.HandlerFunc(oauth2Github))
	http.Handle("/oauth/github/callback", http.HandlerFunc(oauth2GithubCallback))
	http.Handle("/api/gists/fetchAll", http.HandlerFunc(fetchAllGists))	
	http.Handle("/api/gists/search", http.HandlerFunc(search))
	http.Handle("/api/gists", http.HandlerFunc(getGists))	
  http.Handle("/index.jsx", http.HandlerFunc(defaultHandler))

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
