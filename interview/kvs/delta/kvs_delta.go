package kvs_delta

import "errors"

// Errors returned by the package.
var (
	ErrKeyNotFound     = errors.New("key not found")
	ErrNoActiveTx      = errors.New("no active transaction")
	ErrKeyDoesNotExist = errors.New("key does not exist") // used for DELETE outside tx when key missing
)

// entry represents a change in a transaction layer.
type entry struct {
	value   string
	exists  bool // whether this entry is present in the transaction layer
	deleted bool // whether this entry is an explicit delete (tombstone)
}

// Store is the in-memory key-value store with support for nested transactions.
// It is NOT goroutine-safe.
type Store struct {
	committed map[string]string  // main committed data
	txStack   []map[string]entry // stack of transaction overlays (0 = bottom-most active tx)
}

// New returns an initialized Store.
func New() *Store {
	return &Store{
		committed: make(map[string]string),
		txStack:   nil,
	}
}

// Put stores key -> value in the active transaction if present, otherwise in committed store.
func (s *Store) Put(key, value string) {
	if s.inTx() {
		top := s.topTx()
		top[key] = entry{value: value, exists: true, deleted: false}
		return
	}
	s.committed[key] = value
}

// Get returns the value for key or ErrKeyNotFound if not present (or deleted).
func (s *Store) Get(key string) (string, error) {
	// Search from top transaction down to committed map.
	for i := len(s.txStack) - 1; i >= 0; i-- {
		if e, ok := s.txStack[i][key]; ok && e.exists {
			if e.deleted {
				return "", ErrKeyNotFound
			}
			return e.value, nil
		}
	}
	// Fallback to committed store
	if v, ok := s.committed[key]; ok {
		return v, nil
	}
	return "", ErrKeyNotFound
}

// Delete marks the key as deleted in active transaction, or removes from committed store when no tx.
// Returns ErrKeyDoesNotExist if trying to delete a key not present in committed store and no transaction.
func (s *Store) Delete(key string) error {
	if s.inTx() {
		top := s.topTx()
		// Even if key is not present in any underlying layer, we still create a tombstone in the transaction.
		top[key] = entry{exists: true, deleted: true}
		return nil
	}

	// No transaction: delete directly from committed store if present
	if _, ok := s.committed[key]; ok {
		delete(s.committed, key)
		return nil
	}
	return ErrKeyDoesNotExist
}

// Begin starts a new nested transaction.
func (s *Store) Begin() {
	// push a fresh overlay map
	s.txStack = append(s.txStack, make(map[string]entry))
}

// Rollback discards the top-most active transaction.
// Returns ErrNoActiveTx if there is no transaction to rollback.
func (s *Store) Rollback() error {
	if !s.inTx() {
		return ErrNoActiveTx
	}
	// pop
	s.txStack = s.txStack[:len(s.txStack)-1]
	return nil
}

// Commit merges the top-most active transaction into the parent transaction or into the committed store.
// Returns ErrNoActiveTx if there is no active transaction.
func (s *Store) Commit() error {
	if !s.inTx() {
		return ErrNoActiveTx
	}
	top := s.topTx()
	s.txStack = s.txStack[:len(s.txStack)-1] // pop

	if s.inTx() {
		// merge into new top (parent) transaction
		parent := s.topTx()
		for k, e := range top {
			parent[k] = e
		}
		return nil
	}

	// no parent transaction: apply to committed store
	for k, e := range top {
		if e.deleted {
			delete(s.committed, k)
		} else if e.exists {
			s.committed[k] = e.value
		}
	}
	return nil
}

// inTx reports whether there is at least one active transaction.
func (s *Store) inTx() bool {
	return len(s.txStack) > 0
}

// topTx returns the topmost transaction overlay map. Caller must ensure inTx() == true.
func (s *Store) topTx() map[string]entry {
	return s.txStack[len(s.txStack)-1]
}

// NumActiveTransactions returns how many nested transactions are currently active.
func (s *Store) NumActiveTransactions() int {
	return len(s.txStack)
}
