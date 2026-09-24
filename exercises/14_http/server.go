// Package todoapi — ข้อ 14: REST API ด้วย net/http + JSON
//
//	GET    /todos       → 200 [ ...todos เรียงตาม id ] (ว่าง = [])
//	POST   /todos       → 201 todo ที่สร้าง   | JSON พัง/title ว่าง → 400
//	GET    /todos/{id}  → 200 todo            | ไม่เจอ → 404, id ไม่ใช่ตัวเลข → 400
//	DELETE /todos/{id}  → 204 (ไม่มี body)    | ไม่เจอ → 404, id ไม่ใช่ตัวเลข → 400
//
// ทุก error ตอบเป็น {"error": "..."} และทุก JSON response ต้องมี Content-Type: application/json
package todoapi

import (
	"encoding/json"
	"net/http"
)

type Todo struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type Server struct {
	mux *http.ServeMux
	// TODO: storage (อย่าลืมว่า handler ถูกเรียกพร้อมกันหลาย goroutine)
}

func NewServer() *Server {
	s := &Server{mux: http.NewServeMux()}
	// TODO: ลงทะเบียน route (Go 1.22+ รองรับ "GET /todos/{id}")
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// helper ให้ใช้ได้เลย
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
