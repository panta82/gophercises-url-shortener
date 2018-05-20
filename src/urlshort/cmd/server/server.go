package main

import (
	"fmt"
	"net/http"

	"urlshort"
	"urlshort/database"
)

func main() {
	mux := defaultMux()

	// Initialize db at default path
	db, err := database.NewDatabase(database.GetDefaultDatabasePath())
	if err != nil {
		panic(err)
	}

	// Build db handler. This is the lowest handler in the list, because it is going to HDD
	dbHandler, err := urlshort.DatabaseHandler(db, mux)
	if err != nil {
		panic(err)
	}

	// Build the MapHandler using the db as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, dbHandler)

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	yamlHandler, err := urlshort.YAMLHandler([]byte(yaml), mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}

