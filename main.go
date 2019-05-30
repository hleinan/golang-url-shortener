package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
	"log"
	"net/http"
)

type Link struct {
	gorm.Model
	Slug string
	Url  string
}

type Result struct {
	Slug   string `json:"slug"`
	Url    string `json:"url"`
	ApiKey string `json:"api_key"`
}

var (
	gormDB, _  = gorm.Open("sqlite3", "database.db")
	apiKey     = "your-api-key"
	defaultUrl = "https://example.com"
	c          = cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"application/json"},
		Debug:            false,
	})
)

func redirectRootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, defaultUrl, 301)
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	link := Link{}
	if gormDB.Where(&Link{Slug: vars["slug"]}).First(&link).RecordNotFound() {
		http.Redirect(w, r, defaultUrl, 301)
	} else {
		http.Redirect(w, r, link.Url, 301)
	}
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		queryVars := r.URL.Query()
		if queryVars.Get("api_key") != apiKey {
			respondJSON(w, http.StatusForbidden, http.StatusForbidden)
			break
		}
		var links []Link
		gormDB.Find(&links)
		respondJSON(w, http.StatusOK, links)

	case "POST":
		result := GetResult(r)
		if result.ApiKey != apiKey {
			respondJSON(w, http.StatusForbidden, http.StatusForbidden)
			break
		}

		link := Link{}
		if gormDB.Where(&Link{Slug: result.Slug}).First(&link).RecordNotFound() {
			newLink := Link{Slug: result.Slug, Url: result.Url}
			reply := gormDB.Create(&newLink)
			respondJSON(w, http.StatusOK, reply)
		} else {
			respondJSON(w, http.StatusAlreadyReported, link) // error
		}

	case "PUT":
		result := GetResult(r)
		if result.ApiKey != apiKey {
			respondJSON(w, http.StatusForbidden, http.StatusForbidden)
			break
		}

		link := Link{}
		if gormDB.Where(&Link{Slug: result.Slug}).First(&link).RecordNotFound() {
			respondJSON(w, http.StatusNotFound, http.StatusNotFound)
		} else {
			gormDB.Model(&link).Update("url", result.Url)
			respondJSON(w, http.StatusOK, link)
		}

	case "DELETE":
		result := GetResult(r)
		if result.ApiKey != apiKey {
			respondJSON(w, http.StatusForbidden, http.StatusForbidden)
			break
		}

		link := Link{}
		if gormDB.Where(&Link{Slug: result.Slug}).First(&link).RecordNotFound() {
			respondJSON(w, http.StatusNotFound, http.StatusNotFound)
		} else {
			gormDB.Delete(&link)
			respondJSON(w, http.StatusOK, "Deleted")
		}
	}
}
func apiSlugHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	link := Link{}
	if gormDB.Where(&Link{Slug: vars["slug"]}).First(&link).RecordNotFound() {
		respondJSON(w, http.StatusNotFound, http.StatusNotFound)
	} else {
		respondJSON(w, http.StatusOK, link)
	}
}

func main() {
	gormDB.AutoMigrate(&Link{})
	r := mux.NewRouter()

	r.HandleFunc("/api/", apiHandler).Methods("GET", "POST", "PUT", "DELETE")
	r.HandleFunc("/api/{slug}", apiSlugHandler).Methods("GET")
	r.HandleFunc("/{slug}", redirectHandler)
	r.HandleFunc("/", redirectRootHandler)

	handler := c.Handler(r)

	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

func GetResult(r *http.Request) Result {
	var result Result
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&result)
	if err != nil {
		panic(err)
	}
	return result
}
