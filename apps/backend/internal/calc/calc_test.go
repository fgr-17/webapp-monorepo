package calc

import (
	"errors"
	"testing"
)

func TestOperations(t *testing.T) {
	tests := []struct {
		name    string
		op      string
		a, b    float64
		want    float64
		wantErr error
	}{
		{name: "add_positive", op: "add", a: 2, b: 3, want: 5},
		{name: "add_negative", op: "add", a: -2, b: 5, want: 3},
		{name: "subtract", op: "subtract", a: 10, b: 4, want: 6},
		{name: "multiply", op: "multiply", a: 3, b: 4, want: 12},
		{name: "multiply_zero", op: "multiply", a: 3, b: 0, want: 0},
		{name: "divide", op: "divide", a: 10, b: 4, want: 2.5},
		{name: "divide_by_zero", op: "divide", a: 1, b: 0, wantErr: ErrDivideByZero},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				got float64
				err error
			)
			switch tt.op {
			case "add":
				got = Add(tt.a, tt.b)
			case "subtract":
				got = Subtract(tt.a, tt.b)
			case "multiply":
				got = Multiply(tt.a, tt.b)
			case "divide":
				got, err = Divide(tt.a, tt.b)
			default:
				t.Fatalf("unknown op %q", tt.op)
			}

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}
