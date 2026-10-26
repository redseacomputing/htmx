package main

import (
	"fmt"
	"html/template"
	"net/http"
)

// Das komplette Layout (für den direkten Browser-Aufruf oder F5)
const fullTemplate = `<!DOCTYPE html>
<html>
<head>
    <title>HTMX Targeted Navigation</title>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <style>
        body { font-family: sans-serif; margin: 40px; }
        nav a { margin-right: 10px; text-decoration: none; color: #0066cc; cursor: pointer; }
        nav a:hover { text-decoration: underline; }
        #content { margin-top: 20px; padding: 20px; border: 1px solid #ddd; border-radius: 5px; }
    </style>
</head>
<body>
    <h1>My Website</h1>
    <nav>
        <a hx-get="/" hx-target="#content" hx-push-url="true">Home</a> | 
        <a hx-get="/about" hx-target="#content" hx-push-url="true">About</a> | 
        <a hx-get="/contact" hx-target="#content" hx-push-url="true">Contact</a>
    </nav>
    <hr>
    <div id="content">
        {{.Content}}
    </div>
</body>
</html>`

type PageData struct {
	Content template.HTML
}

// Hilfsfunktion: Entweder komplette Seite oder nur das Inhalts-Snippet senden
func renderResponse(w http.ResponseWriter, r *http.Request, content string) {
	// Wenn es kein HTMX-Request ist, schicken wir das komplette Grundgerüst
	if r.Header.Get("HX-Request") == "" {
		tmpl, err := template.New("layout").Parse(fullTemplate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, PageData{Content: template.HTML(content)})
		return
	}
	
	// Wenn es ein HTMX-Request ist, schicken wir NUR das reine HTML-Snippet!
	fmt.Fprint(w, content)
}

func main() {
	mux := http.NewServeMux()

	// 1. Home
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		content := `<h2>Home Page</h2><p>Welcome to the home page (targeted via hx-get).</p>`
		renderResponse(w, r, content)
	})

	// 2. About
	mux.HandleFunc("GET /about", func(w http.ResponseWriter, r *http.Request) {
		content := `<h2>About Page</h2><p>This is the about page, loaded cleanly into #content.</p>`
		renderResponse(w, r, content)
	})

	// 3. Contact
	mux.HandleFunc("GET /contact", func(w http.ResponseWriter, r *http.Request) {
		content := `<h2>Contact Page</h2><p>Get in touch: contact@example.com</p>`
		renderResponse(w, r, content)
	})

	fmt.Println("Targeted Navigation Server running on http://localhost:8080 ...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
