// Package todoapi — เฉลยข้อ 14: REST API ด้วย net/http + JSON + ทดสอบด้วย httptest
package todoapi

import (
	"encoding/json"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type Server struct {
	mu     sync.RWMutex // handler ถูกเรียกพร้อมกันหลาย goroutine (1 request = 1 goroutine)
	todos  map[int]Todo
	nextID int
	mux    *http.ServeMux
}

func NewServer() *Server {
	s := &Server{todos: make(map[int]Todo), nextID: 1, mux: http.NewServeMux()}
	// Go 1.22+ รองรับ method และ path parameter ใน pattern
	s.mux.HandleFunc("GET /todos", s.listTodos)
	s.mux.HandleFunc("POST /todos", s.createTodo)
	s.mux.HandleFunc("GET /todos/{id}", s.getTodo)
	s.mux.HandleFunc("DELETE /todos/{id}", s.deleteTodo)
	return s
}

// ServeHTTP ทำให้ *Server เป็น http.Handler
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// GET /todos → 200 + array เรียงตาม id (ว่างต้องเป็น [] ไม่ใช่ null)
func (s *Server) listTodos(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	list := make([]Todo, 0, len(s.todos)) // make → encode เป็น [] (nil slice จะได้ null)
	for _, t := range s.todos {
		list = append(list, t)
	}
	s.mu.RUnlock()
	sort.Slice(list, func(i, j int) bool { return list[i].ID < list[j].ID })
	writeJSON(w, http.StatusOK, list)
}

// POST /todos {"title": "..."} → 201 + todo ที่สร้าง
// JSON พัง หรือ title ว่าง → 400 {"error": "..."}
func (s *Server) createTodo(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}
	in.Title = strings.TrimSpace(in.Title)
	if in.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}
	s.mu.Lock()
	t := Todo{ID: s.nextID, Title: in.Title}
	s.todos[t.ID] = t
	s.nextID++
	s.mu.Unlock()
	writeJSON(w, http.StatusCreated, t)
}

// parseID อ่าน {id} จาก path — ไม่ใช่ตัวเลข → 400
func parseID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return 0, false
	}
	return id, true
}

// GET /todos/{id} → 200 / 404
func (s *Server) getTodo(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	s.mu.RLock()
	t, found := s.todos[id]
	s.mu.RUnlock()
	if !found {
		writeError(w, http.StatusNotFound, "todo not found")
		return
	}
	writeJSON(w, http.StatusOK, t)
}

// DELETE /todos/{id} → 204 (ไม่มี body) / 404
func (s *Server) deleteTodo(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r)
	if !ok {
		return
	}
	s.mu.Lock()
	_, found := s.todos[id]
	delete(s.todos, id)
	s.mu.Unlock()
	if !found {
		writeError(w, http.StatusNotFound, "todo not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
