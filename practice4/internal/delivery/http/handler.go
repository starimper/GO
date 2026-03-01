package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"practice_4/internal/usecase"
	"practice_4/pkg/modules"

	"github.com/gorilla/mux"
)

type Handler struct {
	usecase *usecase.UserUsecase
}

func NewHandler(u *usecase.UserUsecase) *Handler {
	return &Handler{usecase: u}
}

func (h *Handler) Register(r *mux.Router) {
	r.HandleFunc("/health", h.Health).Methods("GET")
	r.HandleFunc("/users", h.GetUsers).Methods("GET")
	r.HandleFunc("/users/{id}", h.GetUserByID).Methods("GET")
	r.HandleFunc("/users", h.CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", h.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", h.DeleteUser).Methods("DELETE")
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users, err := h.usecase.GetUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(users)
}

func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	user, err := h.usecase.GetUserByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user modules.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	id, err := h.usecase.CreateUser(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var user modules.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user.ID = id

	if err := h.usecase.UpdateUser(&user); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "updated"})
}

func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	rows, err := h.usecase.DeleteUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rows == 0 {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "deleted"})
}