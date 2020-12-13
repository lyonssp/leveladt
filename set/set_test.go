package set

import (
    "github.com/syndtr/goleveldb/leveldb"
    "io/ioutil"
    "fmt"
    "testing"
)

func TestSet(t *testing.T) {
    dir, err := ioutil.TempDir("", "test")
    assert(t, err == nil)

    db, err := leveldb.OpenFile(dir, nil)
    assert(t, err == nil)

    s := Set{db}

    err = s.Add([]byte("foo"))
    assert(t, err == nil)

    contains, err := s.Contains([]byte("foo"))
    assert(t, err == nil)
    assert(t, contains)

    err = s.Remove([]byte("foo"))
    assert(t, err == nil)

    contains, err = s.Contains([]byte("foo"))
    assert(t, err == nil)
    assert(t, !contains)
}

func assert(t *testing.T, b bool) {
    if !b {
        t.Error(fmt.Errorf("assert failed"))
    }
}
