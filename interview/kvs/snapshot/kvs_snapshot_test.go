package kvs_snapshot

import "testing"

func TestBasicOperations(t *testing.T) {
	s := New()
	_, err := s.Get("x")
	if err != ErrKeyNotFound {
		t.Fatalf("expected ErrKeyNotFound, got %v", err)
	}

	s.Put("x", "10")
	v, _ := s.Get("x")
	if v != "10" {
		t.Fatalf("expected 10, got %s", v)
	}

	if err := s.Delete("x"); err != nil {
		t.Fatalf("unexpected delete error: %v", err)
	}
	if _, err := s.Get("x"); err != ErrKeyNotFound {
		t.Fatalf("expected ErrKeyNotFound after delete")
	}
}

func TestTransactionCommitRollback(t *testing.T) {
	s := New()
	s.Put("a", "1")

	s.Begin()
	s.Put("a", "2")
	v, _ := s.Get("a")
	if v != "2" {
		t.Fatalf("expected 2 inside txn, got %s", v)
	}
	if err := s.Rollback(); err != nil {
		t.Fatalf("rollback failed: %v", err)
	}
	v, _ = s.Get("a")
	if v != "1" {
		t.Fatalf("expected 1 after rollback, got %s", v)
	}

	// commit
	s.Begin()
	s.Put("a", "3")
	if err := s.Commit(); err != nil {
		t.Fatalf("commit failed: %v", err)
	}
	v, _ = s.Get("a")
	if v != "3" {
		t.Fatalf("expected 3 after commit, got %s", v)
	}
}

func TestNestedTransactions(t *testing.T) {
	s := New()
	s.Put("x", "0")

	s.Begin() // tx1
	s.Put("x", "1")

	s.Begin() // tx2
	s.Put("x", "2")
	if v, _ := s.Get("x"); v != "2" {
		t.Fatalf("expected 2 in tx2, got %s", v)
	}
	if err := s.Commit(); err != nil {
		t.Fatalf("commit tx2 failed: %v", err)
	}
	if v, _ := s.Get("x"); v != "2" {
		t.Fatalf("expected 2 after commit tx2, got %s", v)
	}
	if err := s.Rollback(); err != nil {
		t.Fatalf("rollback tx1 failed: %v", err)
	}
	if v, _ := s.Get("x"); v != "0" {
		t.Fatalf("expected 0 after rollback, got %s", v)
	}
}
