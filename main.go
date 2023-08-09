package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

var tpl = template.Must(template.ParseFiles("index.html"))

// Function that writes HTTP response that is sent to client and is recieved from client
func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.URL.String())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	params := u.Query()
	searchQuery := params.Get("q")
	page := params.Get("page")
	if page == "" {
		page = "1"
	}

	fmt.Println("Search Query is: ", searchQuery)
	fmt.Println("Page is: ", page)
}

// Main entry point of program
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	// Calls to check if an environment variable for port is set
	port := os.Getenv("PORT")
	// If port is not set, then it is set to 3000
	if port == "" {
		port = "3000"
	}

	fs := http.FileServer(http.Dir("assets"))
	// Creates HTTP request router (multiplexer); used to route incoming requests to correct handler functions
	mux := http.NewServeMux()
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	// Registers the indexHandler to handle requests for the / route
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/search", searchHandler)
	// Starts HTTP server telling it to listen on the port.
	// Mux is the second parameter so the servers uses it to route all incoming requests
	http.ListenAndServe(":"+port, mux)
}
