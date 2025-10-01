package kvs_snapshot

import "errors"

var (
	ErrKeyNotFound = errors.New("key not found")
	ErrNoActiveTx  = errors.New("no active transaction")
)

// Store holds committed data and a stack of transaction snapshots.
type Store struct {
	stack []map[string]string
}

// New initializes a new KV store.
func New() *Store {
	return &Store{
		stack: []map[string]string{make(map[string]string)}, // committed store at index 0
	}
}

// current returns the top-most map (active transaction or committed store).
func (s *Store) current() map[string]string {
	return s.stack[len(s.stack)-1]
}

// Put sets key -> value in the current store.
func (s *Store) Put(key, value string) {
	s.current()[key] = value
}

// Get retrieves a value or returns ErrKeyNotFound.
func (s *Store) Get(key string) (string, error) {
	if v, ok := s.current()[key]; ok {
		return v, nil
	}
	return "", ErrKeyNotFound
}

// Delete removes a key. If the key does not exist, it returns ErrKeyNotFound.
func (s *Store) Delete(key string) error {
	if _, ok := s.current()[key]; !ok {
		return ErrKeyNotFound
	}
	delete(s.current(), key)
	return nil
}

// Begin starts a new transaction by copying the current store.
func (s *Store) Begin() {
	newSnap := make(map[string]string, len(s.current()))
	for k, v := range s.current() {
		newSnap[k] = v
	}
	s.stack = append(s.stack, newSnap)
}

// Rollback discards the current transaction snapshot.
func (s *Store) Rollback() error {
	if len(s.stack) == 1 {
		return ErrNoActiveTx
	}
	s.stack = s.stack[:len(s.stack)-1]
	return nil
}

// Commit replaces the parent snapshot with the current one.
func (s *Store) Commit() error {
	if len(s.stack) == 1 {
		return ErrNoActiveTx
	}
	snap := s.current()
	s.stack = s.stack[:len(s.stack)-1]
	// overwrite parent with committed snapshot
	s.stack[len(s.stack)-1] = snap
	return nil
}

// NumActiveTransactions returns the number of open transactions.
func (s *Store) NumActiveTransactions() int {
	return len(s.stack) - 1
}
