// main.go
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"

	_ "github.com/microsoft/go-mssqldb"
)

// Response structure for /hello endpoint
type HelloResponse struct {
	Message string `json:"message"`
}

type URLShortener struct {
	urls map[string]string
}

var db *sql.DB
var server string
var port int
var user string
var password string
var database string
var selfUrl string
var urlRedirect string
var schema string
var tableName string

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func LoadEnv() {
	// try to load .env file if it exists

	err := godotenv.Load(".env")
	if err != nil {
		// log.Fatal("Error loading .env file")
		fmt.Println("Error loading .env file")
	}

	// Load environment variables
	server = os.Getenv("SQL_URL")
	port, err = strconv.Atoi(os.Getenv("SQL_PORT"))
	if err != nil {
		log.Fatal("Error converting port to int: ", err.Error())
	}
	user = os.Getenv("SQL_USER")
	password = os.Getenv("SQL_USER_PASSWORD")
	database = os.Getenv("SQL_DATABASE")
	selfUrl = os.Getenv("SELF_URL")
	urlRedirect = os.Getenv("URL_REDIRECT")
	schema = os.Getenv("SQL_SCHEMA")
	tableName = os.Getenv("SQL_TABLE")
}

func init() {
	file, err := os.OpenFile("logs.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// defer file.Close()

	if err != nil {
		log.Fatal(err)
	}

	// Create MultiWriter to write to both file and stdout.
	multi := io.MultiWriter(file, os.Stdout)

	_, err = multi.Write([]byte("Direct write to multiwriter\n"))
	if err != nil {
		log.Fatalf("Failed to write to multiwriter: %v", err)
	}

	WarningLogger = log.New(multi,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	ErrorLogger = log.New(multi,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	InfoLogger = log.New(multi,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	LoadEnv()
}

func main() {
	// Build connection string

	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;",
		server, user, password, port, database)

	var err error
	// Create connection pool
	db, err = sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal("Error creating connection pool: ", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Verify the schema and table or create it if it doesn't exist
	InitDB()

	// Create connection string
	shortener := &URLShortener{
		urls: make(map[string]string),
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/create", shortener.HandleAddEmailEntry)
	mux.HandleFunc("/key/", logRequest(shortener.HandleRedirect))

	InfoLogger.Println("Server listening on port 8080")

	http.ListenAndServe(":8080", noCacheMiddleware(mux))
}

func noCacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
		w.Header().Set("Pragma", "no-cache")                                   // HTTP 1.0.
		w.Header().Set("Expires", "0")                                         // Proxies.

		next.ServeHTTP(w, r)
	})
}

func (us *URLShortener) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	shortKey := strings.TrimSuffix(r.URL.Path[len("/key/"):], "/")

	// Log the utilisation to the database
	res, err := UsedLink(shortKey)
	if err != nil {
		ErrorLogger.Println("Error updating database: ", err.Error())
	}

	if res == 0 {
		ErrorLogger.Println("No rows updated in database, key not found: ", shortKey)
	}

	// Redirect the user to the original URL
	http.Redirect(w, r, urlRedirect, http.StatusMovedPermanently)
}

func logRequest(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Received request:", r.Method, r.URL.String())
		next(w, r)
	}
}

func (us *URLShortener) HandleAddEmailEntry(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var newEmail NewEmail
	err := json.NewDecoder(r.Body).Decode(&newEmail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var res int64
	res, err = AddEmailInfo(newEmail)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if res == 0 {
		http.Error(w, "No rows affected", http.StatusInternalServerError)
		return
	} else {
		fmt.Println("Rows affected: ", res)
	}

	// Construct the full shortened URL
	shortenedURL := fmt.Sprintf("%s/key/%s", selfUrl, newEmail.Link)

	println("shortenedURL: ", shortenedURL)

	// Render the HTML response with the shortened URL
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, shortenedURL)

}
