// Package store — เฉลยข้อ 5: error handling, sentinel error, custom error type, wrapping
package store

import (
	"errors"
	"fmt"
	"strings"
)

// sentinel errors — ให้ caller เช็คด้วย errors.Is
var (
	ErrNotFound       = errors.New("user not found")
	ErrDuplicateEmail = errors.New("duplicate email")
)

// ValidationError — custom error type ที่พกข้อมูลเพิ่ม (field ไหนผิด)
// ให้ caller ดึงออกมาด้วย errors.As
type ValidationError struct {
	Field string
	Msg   string
}

// Error ทำให้ *ValidationError implement interface `error`
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed on %s: %s", e.Field, e.Msg)
}

type User struct {
	ID    int
	Name  string
	Email string
}

type Store struct {
	users   map[int]User
	byEmail map[string]int
	nextID  int
}

func NewStore() *Store {
	return &Store{
		users:   make(map[int]User),
		byEmail: make(map[string]int),
		nextID:  1,
	}
}

// Create สร้าง user ใหม่ (ID เริ่มที่ 1 และเพิ่มทีละ 1)
//   - name ว่าง (หลัง trim)         → *ValidationError{Field: "name"}
//   - email ไม่มี "@"              → *ValidationError{Field: "email"}
//   - email ซ้ำ                    → error ที่ wrap ErrDuplicateEmail
func (s *Store) Create(name, email string) (User, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return User{}, &ValidationError{Field: "name", Msg: "must not be empty"}
	}
	if !strings.Contains(email, "@") {
		return User{}, &ValidationError{Field: "email", Msg: "must contain @"}
	}
	if _, exists := s.byEmail[email]; exists {
		return User{}, fmt.Errorf("create user %q: %w", email, ErrDuplicateEmail)
	}
	u := User{ID: s.nextID, Name: name, Email: email}
	s.users[u.ID] = u
	s.byEmail[email] = u.ID
	s.nextID++
	return u, nil
}

// Get หา user ตาม id — ไม่เจอให้ wrap ErrNotFound พร้อม context
func (s *Store) Get(id int) (User, error) {
	u, ok := s.users[id]
	if !ok {
		return User{}, fmt.Errorf("get user %d: %w", id, ErrNotFound)
	}
	return u, nil
}
