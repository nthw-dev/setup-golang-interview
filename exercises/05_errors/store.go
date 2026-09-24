// Package store — ข้อ 5: error handling, sentinel error, custom error type, wrapping
package store

import "errors"

var (
	ErrNotFound       = errors.New("user not found")
	ErrDuplicateEmail = errors.New("duplicate email")
)

// ValidationError บอกว่า field ไหนไม่ถูกต้อง
type ValidationError struct {
	Field string
	Msg   string
}

// Error ต้องคืนรูปแบบ "validation failed on <Field>: <Msg>"
func (e *ValidationError) Error() string {
	// TODO
	return ""
}

type User struct {
	ID    int
	Name  string
	Email string
}

type Store struct {
	// TODO: เก็บข้อมูลอะไรบ้าง?
}

func NewStore() *Store {
	// TODO
	return &Store{}
}

// Create สร้าง user ใหม่ (ID เริ่มที่ 1 และเพิ่มทีละ 1)
//   - name ว่าง (หลัง trim)  → *ValidationError{Field: "name"}
//   - email ไม่มี "@"       → *ValidationError{Field: "email"}
//   - email ซ้ำ             → error ที่ wrap ErrDuplicateEmail
func (s *Store) Create(name, email string) (User, error) {
	// TODO
	return User{}, nil
}

// Get หา user ตาม id — ไม่เจอให้คืน error ที่ wrap ErrNotFound และมี id อยู่ในข้อความ
func (s *Store) Get(id int) (User, error) {
	// TODO
	return User{}, nil
}
