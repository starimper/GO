package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"practice_2/models"
)

var (
	tasks  = make(map[int]models.Task)
	lastID = 0
	mu     sync.Mutex
)

func TasksHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {

	case http.MethodGet:
		idParam := r.URL.Query().Get("id")
		doneParam := r.URL.Query().Get("done")

		mu.Lock()
		defer mu.Unlock()

		if idParam != "" {
			id, err := strconv.Atoi(idParam)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "invalid id"})
				return
			}

			task, ok := tasks[id]
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				json.NewEncoder(w).Encode(map[string]string{"error": "task not found"})
				return
			}

			json.NewEncoder(w).Encode(task)
			return
		}

		filtered := []models.Task{}
		if doneParam != "" {
			done, err := strconv.ParseBool(doneParam)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{"error": "invalid done filter"})
				return
			}

			for _, t := range tasks {
				if t.Done == done {
					filtered = append(filtered, t)
				}
			}

			json.NewEncoder(w).Encode(filtered)
			return
		}

		all := []models.Task{}
		for _, t := range tasks {
			all = append(all, t)
		}
		json.NewEncoder(w).Encode(all)

	case http.MethodPost:
		var body struct {
			Title string `json:"title"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Title == "" {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid title"})
			return
		}

		mu.Lock()
		defer mu.Unlock()

		lastID++
		task := models.Task{
			ID:    lastID,
			Title: body.Title,
			Done:  false,
		}
		tasks[lastID] = task

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(task)

	case http.MethodPatch:
		idParam := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid id"})
			return
		}

		var body struct {
			Done *bool `json:"done"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Done == nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid done value"})
			return
		}

		mu.Lock()
		defer mu.Unlock()

		task, ok := tasks[id]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "task not found"})
			return
		}

		task.Done = *body.Done
		tasks[id] = task

		json.NewEncoder(w).Encode(map[string]bool{"updated": true})

	case http.MethodDelete:
		idParam := r.URL.Query().Get("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid id"})
			return
		}

		mu.Lock()
		defer mu.Unlock()

		if _, ok := tasks[id]; !ok {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(map[string]string{"error": "task not found"})
			return
		}

		delete(tasks, id)
		json.NewEncoder(w).Encode(map[string]bool{"deleted": true})

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
