package kvs_delta

import "testing"

func TestBasicPutGetDeleteNoTx(t *testing.T) {
	s := New()
	_, err := s.Get("x")
	if err != ErrKeyNotFound {
		t.Fatalf("expected ErrKeyNotFound for missing key, got %v", err)
	}

	s.Put("x", "10")
	v, err := s.Get("x")
	if err != nil || v != "10" {
		t.Fatalf("expected 10, got %v %v", v, err)
	}

	// delete
	if err := s.Delete("x"); err != nil {
		t.Fatalf("unexpected delete error: %v", err)
	}
	if _, err := s.Get("x"); err != ErrKeyNotFound {
		t.Fatalf("expected ErrKeyNotFound after delete, got %v", err)
	}
}

func TestTransactionCommitRollbackSimple(t *testing.T) {
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

	// commit case
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

func TestDeleteInsideTransactionAndRollback(t *testing.T) {
	s := New()
	s.Put("k", "v")
	s.Begin()
	if err := s.Delete("k"); err != nil {
		t.Fatalf("unexpected delete error in tx: %v", err)
	}
	if _, err := s.Get("k"); err != ErrKeyNotFound {
		t.Fatalf("expected key to be deleted in tx, got err=%v", err)
	}
	// rollback should restore
	if err := s.Rollback(); err != nil {
		t.Fatalf("rollback failed: %v", err)
	}
	v, err := s.Get("k")
	if err != nil || v != "v" {
		t.Fatalf("expected original value after rollback, got %v %v", v, err)
	}
}

func TestNestedTransactionsCommitBehavior(t *testing.T) {
	s := New()
	s.Put("x", "0")

	s.Begin() // tx1
	s.Put("x", "1")

	s.Begin() // tx2 nested
	s.Put("x", "2")

	// commit tx2 -> merges into tx1
	if err := s.Commit(); err != nil {
		t.Fatalf("commit tx2 failed: %v", err)
	}
	val, _ := s.Get("x")
	if val != "2" {
		t.Fatalf("expected 2 after committing inner tx, got %s", val)
	}

	// rollback tx1 -> should revert to committed (0)
	if err := s.Rollback(); err != nil {
		t.Fatalf("rollback tx1 failed: %v", err)
	}
	if v, err := s.Get("x"); err != nil || v != "0" {
		t.Fatalf("expected 0 after rollback of outer tx, got %v %v", v, err)
	}
}

func TestCommitAllToCommitted(t *testing.T) {
	s := New()
	s.Begin() // tx1
	s.Put("a", "1")
	s.Begin() // tx2
	s.Put("b", "2")

	// commit tx2 -> into tx1
	if err := s.Commit(); err != nil {
		t.Fatalf("commit tx2 failed: %v", err)
	}
	// commit tx1 -> into committed
	if err := s.Commit(); err != nil {
		t.Fatalf("commit tx1 failed: %v", err)
	}

	if v, err := s.Get("a"); err != nil || v != "1" {
		t.Fatalf("expected a=1 in committed, got %v %v", v, err)
	}
	if v, err := s.Get("b"); err != nil || v != "2" {
		t.Fatalf("expected b=2 in committed, got %v %v", v, err)
	}
}

func TestErrorsOnNoActiveTransaction(t *testing.T) {
	s := New()
	if err := s.Commit(); err != ErrNoActiveTx {
		t.Fatalf("expected ErrNoActiveTx when committing with none, got %v", err)
	}
	if err := s.Rollback(); err != ErrNoActiveTx {
		t.Fatalf("expected ErrNoActiveTx when rolling back with none, got %v", err)
	}
}
