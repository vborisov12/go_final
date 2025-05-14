package api

import (
	"log"
	"net/http"
)

type Api struct {
	taskService *TaskService
}

func NewApi(taskService *TaskService) *Api {
	return &Api{
		taskService: taskService,
	}
}

func (a *Api) taskHendler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		log.Printf("Handled POST request to /api/task")
		a.taskService.AddTaskHandler(w, r)
	case http.MethodGet:
		log.Printf("Handled GET request to /api/task")
		a.taskService.GetTaskHandler(w, r)
	case http.MethodPut:
		log.Printf("Handled PUT request to /api/task")
		a.taskService.UpdateTaskHendler(w, r)
	case http.MethodDelete:
		log.Printf("Handled DELETE request to /api/task")
		a.taskService.DeleteTaskHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a *Api) tasksHendler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		log.Printf("Handled GET request to /api/tasks")
		a.taskService.GetTasksHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a *Api) taskDoneHendler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		log.Printf("Handled POST request to /api/task/done")
		a.taskService.DoneTaskHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a *Api) InitRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/nextdate", a.taskService.NextDateHendler)
	mux.HandleFunc("/api/task", a.authMiddleware(a.taskHendler))
	mux.HandleFunc("/api/tasks", a.authMiddleware(a.tasksHendler))
	mux.HandleFunc("/api/task/done", a.authMiddleware(a.taskDoneHendler))
	mux.HandleFunc("/api/signin", a.SignInHandler)

}
