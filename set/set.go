package set

import (
    "fmt"
    "github.com/syndtr/goleveldb/leveldb"
)

type Set struct {
    ldb *leveldb.DB
}

func (s *Set) Add(x []byte) error {
    return s.ldb.Put(x, []byte{}, nil)
}

func (s *Set) Remove(x []byte) error {
    return s.ldb.Delete(x, nil)
}

func (s *Set) Contains(x []byte) (bool, error) {
    _, err := s.ldb.Get(x, nil)
    if err != nil {
        if err == leveldb.ErrNotFound {
            return false, nil
        }
        return false, fmt.Errorf("leveldb get: %v", err)
    }
    return true, nil
}
