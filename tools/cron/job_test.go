package cron

import "testing"

func TestJobId(t *testing.T) {
	expected := "test"

	j := Job{id: expected}

	if j.Id() != expected {
		t.Fatalf("Expected job with id %q, got %q", expected, j.Id())
	}
}

func TestJobExpr(t *testing.T) {
	expected := "1 2 3 4 5"

	s, err := NewSchedule(expected)
	if err != nil {
		t.Fatal(err)
	}

	j := Job{schedule: s}

	if j.Expr() != expected {
		t.Fatalf("Expected job with cron expression %q, got %q", expected, j.Expr())
	}
}

func TestJobRun(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Shouldn't panic: %v", r)
		}
	}()

	calls := ""

	j1 := Job{}
	j2 := Job{fn: func() { calls += "2" }}

	j1.Run()
	j2.Run()

	expected := "2"
	if calls != expected {
		t.Fatalf("Expected calls %q, got %q", expected, calls)
	}
}
