package main

import (
	"fmt"
	"html/template"
	"net/http"
)

// Das HTML-Template als String (wie render_template_string in Flask)
const htmlTemplate = `<!DOCTYPE html>
<html>
<head>
    <title>HTMX Essentials Demo (Go)</title>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
</head>
<body>
    <h1>HTMX Essentials (Golang Version)</h1>

    <h2>Load Message (hx-get)</h2>
    <button hx-get="/message" hx-target="#message-box" hx-swap="innerHTML">
        Load Message
    </button>
    <div id="message-box">
        <p>No message loaded yet.</p>
    </div>

    <hr>

    <h2>Submit Form (hx-post)</h2>
    <form hx-post="/submit" hx-target="#form-result" hx-swap="innerHTML">
        <input type="text" name="username" placeholder="Enter your name" required>
        <button type="submit">Submit</button>
    </form>
    <div id="form-result"></div>
</body>
</html>`

func main() {
	mux := http.NewServeMux()

	// 1. Home page (GET /)
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.New("index").Parse(htmlTemplate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})

	// 2. GET endpoint (/message)
	mux.HandleFunc("GET /message", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "<p>Hello! This content was loaded using hx-get (from Go).</p>")
	})

	// 3. POST endpoint (/submit)
	mux.HandleFunc("POST /submit", func(w http.ResponseWriter, r *http.Request) {
		// Formulardaten parsen
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		
		username := r.FormValue("username")
		if username == "" {
			username = "Guest"
		}

		fmt.Fprintf(w, "<p>Welcome, %s! Your form was submitted successfully via Go.</p>", username)
	})

	fmt.Println("Server is running on http://localhost:8080 ...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
