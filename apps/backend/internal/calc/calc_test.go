package calc

import (
	"errors"
	"testing"
)

func TestAdd(t *testing.T) {
	got := Add(2, 3)
	if got != 5 {
		t.Fatalf("Add(2, 3) = %v, want 5", got)
	}
}

func TestSubtract(t *testing.T) {
	got := Subtract(10, 4)
	if got != 6 {
		t.Fatalf("Subtract(10, 4) = %v, want 6", got)
	}
}

func TestMultiply(t *testing.T) {
	got := Multiply(3, 4)
	if got != 12 {
		t.Fatalf("Multiply(3, 4) = %v, want 12", got)
	}
}

func TestDivide(t *testing.T) {
	got, err := Divide(10, 4)
	if err != nil {
		t.Fatalf("Divide(10, 4) unexpected error: %v", err)
	}
	if got != 2.5 {
		t.Fatalf("Divide(10, 4) = %v, want 2.5", got)
	}
}

func TestDivideByZero(t *testing.T) {
	_, err := Divide(1, 0)
	if !errors.Is(err, ErrDivideByZero) {
		t.Fatalf("Divide(1, 0) error = %v, want ErrDivideByZero", err)
	}
}
