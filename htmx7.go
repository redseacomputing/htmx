package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

// Task-Struktur
type Task struct {
	ID    int
	Title string
}

// In-Memory "Datenbank"
var tasks = []Task{
	{ID: 1, Title: "Learn HTMX"},
	{ID: 2, Title: "Build a CRUD app"},
}
var nextID = 3

const crdTemplate = `
{{define "task_item"}}
<li id="task-{{.ID}}">
    {{.Title}}
    <button hx-get="/edit/{{.ID}}" hx-target="#task-{{.ID}}" hx-swap="outerHTML">Edit</button>
    <button hx-delete="/delete/{{.ID}}" hx-target="#task-list" hx-swap="innerHTML">Delete</button>
</li>
{{end}}

{{define "task_edit"}}
<li id="task-{{.ID}}">
    <form hx-post="/update/{{.ID}}" hx-target="#task-{{.ID}}" hx-swap="outerHTML">
        <input type="text" name="title" value="{{.Title}}">
        <button type="submit">Save</button>
    </form>
</li>
{{end}}

<!DOCTYPE html>
<html>
<head>
    <title>HTMX CRUD App (Go)</title>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
		<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.min.css">
</head>
<body>
    <h1>Task Manager</h1>
    
    <!-- Formular zum Hinzufügen -->
    <form hx-post="/add" hx-target="#task-list" hx-swap="innerHTML">
        <input type="text" name="title" placeholder="New task" required>
        <button type="submit">Add Task</button>
    </form>

    <!-- Aufgabenliste -->
    <ul id="task-list">
        {{range .Tasks}}
            {{template "task_item" .}}
        {{end}}
    </ul>
</body>
</html>
`

func main() {
	mux := http.NewServeMux()
	tmpl := template.Must(template.New("index").Parse(crdTemplate))

	// Helper-Funktion: Rendert die komplette Liste aller Tasks neu
	renderTaskList := func(w http.ResponseWriter) {
		for _, t := range tasks {
			tmpl.ExecuteTemplate(w, "task_item", t)
		}
	}

	// 1. Home / Read
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Tasks []Task
		}{Tasks: tasks}
		tmpl.Execute(w, data)
	})

	// 2. Create (Add)
	mux.HandleFunc("POST /add", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		title := r.FormValue("title")
		if title != "" {
			tasks = append(tasks, Task{ID: nextID, Title: title})
			nextID++
		}
		renderTaskList(w)
	})

	// 3. Delete
	mux.HandleFunc("DELETE /delete/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.PathValue("id"))
		
		// Task aus der Liste filtern
		newTasks := []Task{}
		for _, t := range tasks {
			if t.ID != id {
				newTasks = append(newTasks, t)
			}
		}
		tasks = newTasks
		renderTaskList(w)
	})

	// 4. Edit (Formular anzeigen)
	mux.HandleFunc("GET /edit/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.PathValue("id"))
		for _, t := range tasks {
			if t.ID == id {
				tmpl.ExecuteTemplate(w, "task_edit", t)
				return
			}
		}
		w.Write([]byte("<li>Task not found</li>"))
	})

	// 5. Update (Änderung speichern)
	mux.HandleFunc("POST /update/{id}", func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(r.PathValue("id"))
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}
		newTitle := r.FormValue("title")

		var updatedTask *Task
		for i := range tasks {
			if tasks[i].ID == id {
				tasks[i].Title = newTitle
				updatedTask = &tasks[i]
				break
			}
		}

		if updatedTask != nil {
			// Nur das einzelne Listenelement zurückschicken, das aktualisiert wurde
			tmpl.ExecuteTemplate(w, "task_item", *updatedTask)
		} else {
			w.Write([]byte("<li>Task not found</li>"))
		}
	})

	fmt.Println("CRUD Server running on http://localhost:8080 ...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
