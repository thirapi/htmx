package main

import (
	"html/template"
	"net/http"
	"strconv"
	"sync"
)

type Task struct {
	ID     int
	Title  string
	Status string
}

var (
	tasks  = []Task{{ID: 1, Title: "Learn HTMX", Status: "todo"}}
	nextID = 2
	mu     sync.Mutex
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/tasks", tasksHandler)
	http.HandleFunc("/tasks/move", moveTaskHandler)
	http.HandleFunc("/tasks/add", addTaskHandler)
	http.HandleFunc("/tasks/delete", deleteTaskHandler)

	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html", "templates/board.html"))
	tmpl.ExecuteTemplate(w, "index.html", nil)
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("board").Funcs(template.FuncMap{
		"slice": func(s ...string) []string { return s },
	}).ParseFiles("templates/board.html"))
	data := struct {
		Tasks []Task
	}{Tasks: tasks}
	tmpl.ExecuteTemplate(w, "board", data)
}

func moveTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("id"))
	status := r.FormValue("status")
	mu.Lock()
	for i, t := range tasks {
		if t.ID == id {
			tasks[i].Status = status
		}
	}
	mu.Unlock()
	tasksHandler(w, r)
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	mu.Lock()
	tasks = append(tasks, Task{ID: nextID, Title: title, Status: "todo"})
	nextID++
	mu.Unlock()
	tasksHandler(w, r)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.FormValue("id"))
	mu.Lock()
	newTasks := []Task{}
	for _, t := range tasks {
		if t.ID != id {
			newTasks = append(newTasks, t)
		}
	}
	tasks = newTasks
	mu.Unlock()
	tasksHandler(w, r)
}
