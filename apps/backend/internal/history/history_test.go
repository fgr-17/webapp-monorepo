package history

import (
	"testing"
	"time"
)

func TestAddAndListNewestFirst(t *testing.T) {
	s := New(10)
	t1 := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	t2 := t1.Add(time.Second)

	s.Add("add", 1, 2, 3, t1)
	s.Add("multiply", 3, 4, 12, t2)

	got := s.List()
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
	if got[0].Op != "multiply" || got[0].Result != 12 {
		t.Fatalf("newest = %+v, want multiply/12", got[0])
	}
	if got[1].Op != "add" || got[1].Result != 3 {
		t.Fatalf("oldest = %+v, want add/3", got[1])
	}
	if !got[0].Timestamp.Equal(t2) || !got[1].Timestamp.Equal(t1) {
		t.Fatalf("timestamps = %v, %v", got[0].Timestamp, got[1].Timestamp)
	}
}

func TestCapDropsOldest(t *testing.T) {
	s := New(3)
	base := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := 0; i < 5; i++ {
		s.Add("add", float64(i), 1, float64(i+1), base.Add(time.Duration(i)*time.Second))
	}

	if s.Len() != 3 {
		t.Fatalf("len = %d, want 3", s.Len())
	}
	got := s.List()
	// Newest first: adds with a=4,3,2 (0 and 1 dropped).
	wantA := []float64{4, 3, 2}
	for i, want := range wantA {
		if got[i].A != want {
			t.Fatalf("got[%d].A = %v, want %v (full=%+v)", i, got[i].A, want, got)
		}
	}
}

func TestListReturnsCopy(t *testing.T) {
	s := New(10)
	s.Add("add", 1, 1, 2, time.Unix(1, 0).UTC())
	got := s.List()
	got[0].Op = "mutated"
	again := s.List()
	if again[0].Op != "add" {
		t.Fatalf("store mutated via List copy: %q", again[0].Op)
	}
}

func TestNewDefaultCap(t *testing.T) {
	s := New(0)
	for i := 0; i < DefaultCap+3; i++ {
		s.Add("add", float64(i), 0, float64(i), time.Time{})
	}
	if s.Len() != DefaultCap {
		t.Fatalf("len = %d, want DefaultCap %d", s.Len(), DefaultCap)
	}
}
