package set

import (
    "github.com/syndtr/goleveldb/leveldb"
    "github.com/stretchr/testify/assert"
    "io/ioutil"
    "testing"
)

func TestSet(t *testing.T) {
    assert := assert.New(t)

    dir, err := ioutil.TempDir("", "test")
    assert.Nil(err)

    db, err := leveldb.OpenFile(dir, nil)
    assert.Nil(err)

    s := Set{
        ns: []byte("xxx"),
        ldb: db,
    }

    err = s.Add([]byte("foo"))
    assert.Nil(err)

    contains, err := s.Contains([]byte("foo"))
    assert.Nil(err)
    assert.True(contains)

    err = s.Remove([]byte("foo"))
    assert.Nil(err)

    contains, err = s.Contains([]byte("foo"))
    assert.Nil(err)
    assert.False(contains)
}
