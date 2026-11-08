package main

import (
	"encoding/json"
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
	tmpl   *template.Template
)

func init() {
	tmpl = template.Must(template.New("").Funcs(template.FuncMap{
		"slice": func(s ...string) []string { return s },
		"json": func(v interface{}) string {
			b, _ := json.Marshal(v)
			return string(b)
		},
		"dict": func(values ...interface{}) (map[string]interface{}, error) {
			dict := make(map[string]interface{})
			for i := 0; i < len(values); i += 2 {
				dict[values[i].(string)] = values[i+1]
			}
			return dict, nil
		},
	}).ParseFiles("templates/index.html", "templates/board.html"))
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/tasks", tasksHandler)
	http.HandleFunc("/tasks/move", moveTaskHandler)
	http.HandleFunc("/tasks/add", addTaskHandler)
	http.HandleFunc("/tasks/delete", deleteTaskHandler)

	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "index.html", nil)
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	data := struct {
		Tasks []Task
	}{Tasks: append([]Task(nil), tasks...)}
	mu.Unlock()
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
