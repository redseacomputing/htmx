package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

const formTemplate = `<!DOCTYPE html>
<html>
<head>
    <title>HTMX Forms Demo (Go)</title>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
</head>
<body>
    <h1>Contact Form</h1>
    <form hx-post="/submit" hx-target="#response" hx-swap="innerHTML">
        <label>Name:</label><br>
        <input type="text" name="name" required><br><br>
        <label>Email:</label><br>
        <input type="email" name="email" required><br><br>
        <label>Message:</label><br>
        <textarea name="message" required></textarea><br><br>
        <button type="submit">Send</button>
    </form>
    <div id="response"></div>
</body>
</html>`

func main() {
	mux := http.NewServeMux()

	// 1. Startseite mit dem Formular liefern
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.New("form").Parse(formTemplate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	})

	// 2. POST-Endpunkt für die Formularverarbeitung
	mux.HandleFunc("POST /submit", func(w http.ResponseWriter, r *http.Request) {
		// Wichtig in Go: Formulardaten erst einlesen/parsen
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		name := strings.TrimSpace(r.FormValue("name"))
		email := strings.TrimSpace(r.FormValue("email"))
		message := strings.TrimSpace(r.FormValue("message"))

		// 1. Validierung: Sind alle Felder ausgefüllt?
		if name == "" || email == "" || message == "" {
			fmt.Fprint(w, `<p style="color:red;">All fields are required.</p>`)
			return
		}

		// 2. Validierung: Enthält die E-Mail ein '@'?
		if !strings.Contains(email, "@") {
			fmt.Fprint(w, `<p style="color:red;">Invalid email address.</p>`)
			return
		}

		// Erfolgreiche Antwort zurückschicken
		fmt.Fprintf(w, `<p style="color:green;">Thank you, %s! Your message has been received.</p>`, name)
	})

	fmt.Println("Server is running on http://localhost:8080 ...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
