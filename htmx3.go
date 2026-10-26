package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"time"
)

const eventsTemplate = `<!DOCTYPE html>
<html>
<head>
    <title>HTMX Events Demo (Go)</title>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
</head>
<body>
    <h1>Dynamic Content with HTMX (Go Version)</h1>

    <h2>Live Search</h2>
    <input type="text" name="query" placeholder="Search fruits..." 
           hx-get="/search" 
           hx-trigger="keyup changed delay:500ms" 
           hx-target="#search-results" 
           hx-swap="innerHTML">
    
    <div id="search-results">
        <p>Start typing to see results...</p>
    </div>

    <hr>

    <h2>Live Time (Polling)</h2>
    <div hx-get="/time" hx-trigger="every 2s" hx-swap="innerHTML">
        <p>Loading time...</p>
    </div>
</body>
</html>`

func main() {
	mux := http.NewServeMux()

	// 1. Startseite ausliefern
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.New("index").Parse(eventsTemplate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})

	// 2. Such-Endpunkt (Live Search)
	mux.HandleFunc("GET /search", func(w http.ResponseWriter, r *http.Request) {
		// Query-Parameter aus der URL auslesen (?query=...)
		query := strings.ToLower(r.URL.Query().Get("query"))
		data := []string{"Apple", "Banana", "Orange", "Grape", "Watermelon"}

		var results []string
		for _, item := range data {
			if strings.Contains(strings.ToLower(item), query) {
				results = append(results, item)
			}
		}

		// Wenn keine Treffer da sind
		if len(results) == 0 {
			fmt.Fprint(w, "<p>No results found</p>")
			return
		}

		// Treffer als HTML zurückgeben
		for _, item := range results {
			fmt.Fprintf(w, "<p>%s</p>", item)
		}
	})

	// 3. Zeit-Endpunkt (Polling alle 2 Sekunden)
	mux.HandleFunc("GET /time", func(w http.ResponseWriter, r *http.Request) {
		// Aktuelle Uhrzeit im Format HH:MM:SS formatieren
		now := time.Now().Format("15:04:05")
		fmt.Fprintf(w, "<p>Current time: %s</p>", now)
	})

	fmt.Println("Server is running on http://localhost:8080 ...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
