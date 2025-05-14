package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/vborisov12/go_final/pkg/storage"
)

type TaskService struct {
	storage storage.TaskStore
}

func NewTaskService(storage storage.TaskStore) *TaskService {
	return &TaskService{storage: storage}
}

type taskResponse struct {
}

func writeJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// NextDateHendler обработчик запроса /api/nextdate
// опирается на параметры now, date, repeat
// для обработки исопльзует функцию NewDate из scheduler.go
func (s TaskService) NextDateHendler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse("20060102", nowStr)
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid 'now' date format"), http.StatusBadRequest)
			log.Printf("Invalid 'now' date format: %v", err)
			return
		}
	}

	if dateStr == "" {
		http.Error(w, fmt.Sprintf("Missing 'date' parameter"), http.StatusBadRequest)
		log.Print("Missing 'date' parameter")
		return
	}

	if repeat == "" {
		http.Error(w, fmt.Sprintf("Missing 'repeat' parameter"), http.StatusBadRequest)
		log.Print("Missing 'repeat' parameter")
		return
	}

	nexDate, err := NewDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	log.Printf("Next date response: %s", nexDate)
	fmt.Fprint(w, nexDate)
}

// AddTask обработчик запроса /api/task
// POST запрос с телом в формате JSON
func (s *TaskService) AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var newTask storage.Task
	var respError struct {
		Error string `json:"error"`
	}
	var respID struct {
		ID int64 `json:"id"`
	}

	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		log.Printf("Error while decoding JSON: %v", err)
		respError.Error = fmt.Sprintf("Error while decoding JSON: %v", err)
		writeJSON(w, respError, http.StatusBadRequest)
		return
	}

	if newTask.Title == "" {
		log.Print("Missing 'title' parameter")

		respError.Error = fmt.Sprintf("Missing 'title' parameter")
		writeJSON(w, respError, http.StatusBadRequest)
		return
	}

	err = checkDate(&newTask)
	if err != nil {
		log.Printf("Error: %v", err)
		respError.Error = err.Error()
		writeJSON(w, respError, http.StatusBadRequest)
		return
	}

	id, err := s.storage.AddTask(newTask)
	if err != nil {
		log.Printf("Error: %v", err)
		respError.Error = err.Error()
		writeJSON(w, respError, http.StatusInternalServerError)
		return
	}

	respID.ID = id
	log.Printf("Handler response AddTask ID: %d", id)
	writeJSON(w, respID, http.StatusCreated)
}

// GetTasksHandler обработчик запроса /api/tasks
// GET запрос
// Параметр search опциональный, необходим для поиска задач
func (s *TaskService) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	var respError struct {
		Error string `json:"error"`
	}
	var respTasks struct {
		Tasks []*storage.Task `json:"tasks"`
	}

	// Проверка на наличие параметра search
	if r.URL.Query().Has("search") {
		// Проверка на дату в параметре
		if t, err := time.Parse("02.01.2006", r.URL.Query().Get("search")); err == nil {
			respTasks.Tasks, _ = s.storage.SearchTasksByDate(t.Format("20060102"), 50)
			log.Printf("Search by date: %s -> result: %d", t.Format("20060102"), len(respTasks.Tasks))
			writeJSON(w, respTasks, http.StatusOK)
			return
		}
		respTasks.Tasks, _ = s.storage.SearchTasks(r.URL.Query().Get("search"), 50)
		log.Printf("Search: %s -> result: %d", r.URL.Query().Get("search"), len(respTasks.Tasks))
		writeJSON(w, respTasks, http.StatusOK)
		return
	}

	tasks, err := s.storage.GetTasks(50)

	if err != nil {
		log.Printf("Error: %v", err)
		respError.Error = err.Error()
		writeJSON(w, respError, http.StatusInternalServerError)
		return
	}

	respTasks.Tasks = tasks
	log.Printf("Executed GetTasksHandler: %d tasks", len(respTasks.Tasks))
	writeJSON(w, respTasks, http.StatusOK)
}

func (s *TaskService) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	var respError struct {
		Error string `json:"error"`
	}

	if r.URL.Query().Has("id") {
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			log.Print("Missing 'id' parameter")
			respError.Error = fmt.Sprintf("Missing 'id' parameter")
			writeJSON(w, respError, http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Printf("Error: %v", err)
			respError.Error = err.Error()
			writeJSON(w, respError, http.StatusBadRequest)
			return
		}

		task, err := s.storage.GetTaskByID(id)
		if err != nil {
			log.Printf("Error: %v", err)
			respError.Error = err.Error()
			writeJSON(w, respError, http.StatusInternalServerError)
			return
		}
		log.Printf("Executed GetTaskHandler, search by id: %d", id)
		writeJSON(w, task, http.StatusOK)
		return
	}

	log.Print("Missing 'id' parameter")
	respError.Error = fmt.Sprintf("Missing 'id' parameter")
	writeJSON(w, respError, http.StatusBadRequest)
}

func (s *TaskService) UpdateTaskHendler(w http.ResponseWriter, r *http.Request) {
	var task storage.Task
	var respError struct {
		Error string `json:"error"`
	}
	resp := struct{}{}

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		log.Printf("Error while decoding JSON: %v", err)
		respError.Error = fmt.Sprintf("Error while decoding JSON: %v", err)
		writeJSON(w, respError, http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		log.Print("Missing 'id' parameter")
		respError.Error = fmt.Sprintf("Missing 'id' parameter")
		writeJSON(w, respError, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		log.Print("Missing 'title' parameter")
		respError.Error = fmt.Sprintf("Missing 'title' parameter")
		writeJSON(w, respError, http.StatusBadRequest)
		return
	}

	err = checkDate(&task)
	if err != nil {
		log.Printf("Error: %v", err)
		respError.Error = err.Error()
		writeJSON(w, respError, http.StatusBadRequest)
		return
	}

	err = s.storage.UpdateTaskByID(&task)
	if err != nil {
		log.Printf("Error: %v", err)
		respError.Error = err.Error()
		writeJSON(w, respError, http.StatusInternalServerError)
		return
	}

	log.Printf("Executed UpdateTaskHendler by id: %s", task.ID)
	writeJSON(w, resp, http.StatusOK)

}

// DeleteTaskHandler обработчик запроса /api/task
// DELETE запрос с параметром id
func (s *TaskService) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	var respError struct {
		Error string `json:"error"`
	}

	resp := struct{}{}

	if r.URL.Query().Has("id") {
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			log.Print("Missing 'id' parameter")
			respError.Error = fmt.Sprintf("Missing 'id' parameter")
			writeJSON(w, respError, http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Printf("Error: %v", err)
			respError.Error = err.Error()
			writeJSON(w, respError, http.StatusBadRequest)
			return
		}

		err = s.storage.DeleteTaskByID(id)
		if err != nil {
			log.Printf("Error: %v", err)
			respError.Error = err.Error()
			writeJSON(w, respError, http.StatusInternalServerError)
			return
		}
		log.Printf("Executed DeleteTaskHandler, search by id: %d", id)
		writeJSON(w, resp, http.StatusOK)
		return
	}

	log.Print("Missing 'id' parameter")
	respError.Error = fmt.Sprintf("Missing 'id' parameter")
	writeJSON(w, respError, http.StatusBadRequest)
}

func (s *TaskService) DoneTaskHandler(w http.ResponseWriter, r *http.Request) {
	var respError struct {
		Error string `json:"error"`
	}
	resp := struct{}{}

	if r.URL.Query().Has("id") {
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			log.Print("Missing 'id' parameter")
			respError.Error = fmt.Sprintf("Missing 'id' parameter")
			writeJSON(w, respError, http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Printf("Error: %v", err)
			respError.Error = err.Error()
			writeJSON(w, respError, http.StatusBadRequest)
			return
		}

		task, err := s.storage.GetTaskByID(id)
		if err != nil {
			log.Printf("Error: %v", err)
			respError.Error = err.Error()
			writeJSON(w, respError, http.StatusInternalServerError)
			return
		}

		log.Printf("Mark task with id: %d as done", id)
		if task.Repeat == "" {
			log.Printf("Task with id: %d is not repeatable, delete it", id)
			err = s.storage.DeleteTaskByID(id)
			if err != nil {
				log.Printf("Error: %v", err)
				respError.Error = err.Error()
				writeJSON(w, respError, http.StatusInternalServerError)
				return
			}
			writeJSON(w, resp, http.StatusOK)
			return
		}

		nextDate, err := NewDate(time.Now(), task.Date, task.Repeat)
		if err != nil {
			log.Printf("Error: %v", err)
			respError.Error = err.Error()
			writeJSON(w, respError, http.StatusInternalServerError)
			return
		}
		log.Printf("Next exec date %s for task with id: %d", nextDate, id)

		task.Date = nextDate
		err = s.storage.UpdateTaskByID(task)
		if err != nil {
			log.Printf("Error: %v", err)
			respError.Error = err.Error()
			writeJSON(w, respError, http.StatusInternalServerError)
			return
		}

		writeJSON(w, resp, http.StatusOK)
		return
	}

}
