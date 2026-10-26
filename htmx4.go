package main

import (
	"fmt"
	"html/template"
	"net/http"
)

// Das HTML-Layout mit hx-boost="true" auf dem Body
const fullTemplate = `<!DOCTYPE html>
<html>
<head>
    <title>HTMX Navigation Demo (Go Boosted)</title>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <style>
        body { font-family: sans-serif; margin: 40px; }
        nav a { margin-right: 10px; text-decoration: none; color: #0066cc; }
        nav a:hover { text-decoration: underline; }
        #content { margin-top: 20px; padding: 20px; border: 1px solid #ddd; border-radius: 5px; }
    </style>
</head>
<body hx-boost="true">
    <h1>My Website (Go + hx-boost)</h1>
    <nav>
        <a href="/">Home</a> | 
        <a href="/about">About</a> | 
        <a href="/contact">Contact</a>
    </nav>
    <hr>
    <div id="content">
        {{.Content}}
    </div>
</body>
</html>`

// Hilfstruktur, um Inhalte an das Template zu übergeben
type PageData struct {
	Content template.HTML
}

func renderPage(w http.ResponseWriter, htmlContent string) {
	tmpl, err := template.New("layout").Parse(fullTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	data := PageData{
		Content: template.HTML(htmlContent),
	}
	tmpl.Execute(w, data)
}

func main() {
	mux := http.NewServeMux()

	// 1. Home Route
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		content := `<h2>Home Page</h2><p>Welcome to the home page (loaded via hx-boost).</p>`
		renderPage(w, content)
	})

	// 2. About Route
	mux.HandleFunc("GET /about", func(w http.ResponseWriter, r *http.Request) {
		content := `<h2>About Page</h2><p>This is the about page. No extra logic needed in Go!</p>`
		renderPage(w, content)
	})

	// 3. Contact Route
	mux.HandleFunc("GET /contact", func(w http.ResponseWriter, r *http.Request) {
		content := `<h2>Contact Page</h2><p>Get in touch with us at contact@example.com.</p>`
		renderPage(w, content)
	})

	fmt.Println("Boosted Server running on http://localhost:8080 ...")
	http.ListenAndServe(":8080", mux)
}
