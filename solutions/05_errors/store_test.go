package store

import (
	"errors"
	"strings"
	"testing"
)

func TestCreateAndGet(t *testing.T) {
	s := NewStore()
	u, err := s.Create("Tah", "tah@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if u.ID != 1 || u.Name != "Tah" {
		t.Fatalf("Create = %+v; want ID=1 Name=Tah", u)
	}
	u2, _ := s.Create("Bee", "bee@example.com")
	if u2.ID != 2 {
		t.Errorf("second user ID = %d; want 2", u2.ID)
	}

	got, err := s.Get(1)
	if err != nil || got != u {
		t.Fatalf("Get(1) = %+v, %v; want %+v, nil", got, err, u)
	}
}

func TestGetNotFound(t *testing.T) {
	s := NewStore()
	_, err := s.Get(42)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("Get(42) error = %v; want wraps ErrNotFound", err)
	}
	if !strings.Contains(err.Error(), "42") {
		t.Errorf("error %q should mention the id", err)
	}
}

func TestValidation(t *testing.T) {
	s := NewStore()
	tests := []struct {
		name, email, field string
	}{
		{"  ", "a@b.c", "name"},
		{"Tah", "no-at-sign", "email"},
	}
	for _, tt := range tests {
		_, err := s.Create(tt.name, tt.email)
		var ve *ValidationError
		if !errors.As(err, &ve) {
			t.Errorf("Create(%q,%q) error = %v; want *ValidationError", tt.name, tt.email, err)
			continue
		}
		if ve.Field != tt.field {
			t.Errorf("Field = %q; want %q", ve.Field, tt.field)
		}
		if !strings.Contains(ve.Error(), "validation failed on "+tt.field) {
			t.Errorf("Error() = %q", ve.Error())
		}
	}
}

func TestDuplicateEmail(t *testing.T) {
	s := NewStore()
	if _, err := s.Create("A", "same@x.com"); err != nil {
		t.Fatal(err)
	}
	_, err := s.Create("B", "same@x.com")
	if !errors.Is(err, ErrDuplicateEmail) {
		t.Fatalf("error = %v; want wraps ErrDuplicateEmail", err)
	}
}
