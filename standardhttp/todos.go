package main

import (
	"net/http"
	"time"
)

type Todo struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Name      string    `json:"name"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"createdAt"`
}

func (s *Server) handleTodos(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleTodosGet(w, r)
	case http.MethodPost:
		s.handleTodosPost(w, r)
	case http.MethodPut:
		s.handleTodosPut(w, r)
	case http.MethodDelete:
		s.handleTodosDelete(w, r)
	case http.MethodOptions:
		w.Header().Add("Access-Control-Allow-Methods", "DELETE")
		respond(w, r, http.StatusOK, nil)
	default:
		respondHTTPErr(w, r, http.StatusNotFound)
	}
}

func (s *Server) handleTodosGet(w http.ResponseWriter, r *http.Request) {
	p := NewPath(r.URL.Path)
	if p.HasID() {
		response := Todo{}
		if err := s.database.First(&response, p.ID).Error; err != nil {
			respondErr(w, r, http.StatusNotFound, "failed to retrieve todo", err)
			return
		}
		respond(w, r, http.StatusOK, response)
	} else {
		response := []Todo{}
		s.database.Find(&response)
		respond(w, r, http.StatusOK, response)
	}
}

func (s *Server) handleTodosPost(w http.ResponseWriter, r *http.Request) {
	todo := Todo{}
	if err := decodeBody(r, &todo); err != nil {
		respondErr(w, r, http.StatusBadRequest, "failed to read todo from request", err)
		return
	}

	if err := s.database.Create(&todo).Error; err != nil {
		respondErr(w, r, http.StatusNotFound, "failed to create todo", err)
		return
	}

	respond(w, r, http.StatusCreated, todo)
}

func (s *Server) handleTodosPut(w http.ResponseWriter, r *http.Request) {
	p := NewPath(r.URL.Path)
	if !p.HasID() {
		respondErr(w, r, http.StatusMethodNotAllowed, "cannot update all todos")
		return
	}

	todo := Todo{}
	if err := s.database.First(&todo, p.ID).Error; err != nil {
		respondErr(w, r, http.StatusNotFound, "failed to retrieve todo", err)
		return
	}

	if err := decodeBody(r, &todo); err != nil {
		respondErr(w, r, http.StatusBadRequest, "failed to read todo from request", err)
		return
	}

	if err := s.database.Save(&todo).Error; err != nil {
		respondErr(w, r, http.StatusInternalServerError, "failed to update todo", err)
		return
	}

	respond(w, r, http.StatusOK, todo)
}

func (s *Server) handleTodosDelete(w http.ResponseWriter, r *http.Request) {
	p := NewPath(r.URL.Path)
	if !p.HasID() {
		respondErr(w, r, http.StatusMethodNotAllowed, "cannot delete all todos")
		return
	}

	todo := Todo{}
	if err := s.database.First(&todo, p.ID).Error; err != nil {
		respondErr(w, r, http.StatusNotFound, "failed to retrieve todo", err)
		return
	}

	if err := s.database.Delete(&todo).Error; err != nil {
		respondErr(w, r, http.StatusInternalServerError, "failed to delete todo", err)
		return
	}

	respond(w, r, http.StatusOK, nil)
}
