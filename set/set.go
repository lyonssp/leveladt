package set

import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
)

type Set struct {
	ns  []byte
	ldb *leveldb.DB
}

// Add includes the value x to the set
func (s *Set) Add(x []byte) error {
	return s.ldb.Put(s.key(x), []byte{}, nil)
}

// Remove deletes the value x from the set
func (s *Set) Remove(x []byte) error {
	return s.ldb.Delete(s.key(x), nil)
}

// Contains returns true if x is in the set, and false otherwise
func (s *Set) Contains(x []byte) (bool, error) {
	_, err := s.ldb.Get(s.key(x), nil)
	if err != nil {
		if err == leveldb.ErrNotFound {
			return false, nil
		}
		return false, fmt.Errorf("leveldb get: %v", err)
	}
	return true, nil
}

func (s *Set) key(x []byte) []byte {
	namespaced := make([]byte, len(x)+len(s.ns))
	copy(namespaced[:len(s.ns)], s.ns)
	copy(namespaced[len(s.ns):], x)
	return namespaced
}
