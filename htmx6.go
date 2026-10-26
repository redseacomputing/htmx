package main

import (
	"html/template"
	"net/http"
	"strconv"
)

type User struct {
	ID    int
	Name  string
	Email string
}

var users = []User{
	{ID: 1, Name: "Alice", Email: "alice@example.com"},
	{ID: 2, Name: "Bob", Email: "bob@example.com"},
	{ID: 3, Name: "Charlie", Email: "charlie@example.com"},
}

// Das Template enthält jetzt sowohl das Hauptlayout als auch die wiederverwendbare Komponente ("user_card")
const tplSource = `
{{define "user_card"}}
<div class="user-card" style="border:1px solid #ccc; padding:10px; margin:5px;">
    <h3>{{.Name}}</h3>
    <p>Email: {{.Email}}</p>
    <button hx-get="/user/{{.ID}}" hx-target="#details" hx-swap="innerHTML">
        View Details
    </button>
</div>
{{end}}

<!DOCTYPE html>
<html>
<head>
    <title>Go Templates Demo</title>
    <script src="https://unpkg.com/hx.org@1.9.10"></script>
</head>
<body>
    <h1>User List</h1>
    <div id="user-list">
        {{range .Users}}
            {{template "user_card" .}}
        {{end}}
    </div>
    <hr>
    <h2>User Details</h2>
    <div id="details">
        <p>Select a user to view details.</p>
    </div>
</body>
</html>
`

func main() {
	mux := http.NewServeMux()

	// Template einmal beim Start parsen
	tmpl := template.Must(template.New("index").Parse(tplSource))

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		// Daten an das Template übergeben
		data := struct {
			Users []User
		}{
			Users: users,
		}
		tmpl.Execute(w, data)
	})

	mux.HandleFunc("GET /user/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.PathValue("id"))
		for _, u := range users {
			if u.ID == id {
				// Hier nutzen wir direkt denselben "user_card"-Baustein oder ein eigenes Detail-Template!
				tmpl.ExecuteTemplate(w, "user_card", u)
				return
			}
		}
		w.Write([]byte("<p>User not found</p>"))
	})

	http.ListenAndServe(":8080", mux)
}
