package set

import (
	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"
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
		ns:  []byte("xxx"),
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

func TestNamespacing(t *testing.T) {
	t.Run("Add", func(t *testing.T) {
		assert := assert.New(t)

		dir, err := ioutil.TempDir("", "test")
		assert.Nil(err)

		db, err := leveldb.OpenFile(dir, nil)
		assert.Nil(err)

		a := Set{
			ns:  []byte("xxx"),
			ldb: db,
		}
		b := Set{
			ns:  []byte("yyy"),
			ldb: db,
		}

		err = a.Add([]byte("foo"))
		assert.Nil(err)

		contains, err := a.Contains([]byte("foo"))
		assert.Nil(err)
		assert.True(contains)

		contains, err = b.Contains([]byte("foo"))
		assert.Nil(err)
		assert.False(contains)
	})

	t.Run("Remove", func(t *testing.T) {
		assert := assert.New(t)

		dir, err := ioutil.TempDir("", "test")
		assert.Nil(err)

		db, err := leveldb.OpenFile(dir, nil)
		assert.Nil(err)

		a := Set{
			ns:  []byte("xxx"),
			ldb: db,
		}
		b := Set{
			ns:  []byte("yyy"),
			ldb: db,
		}

		err = a.Add([]byte("foo"))
		assert.Nil(err)

		err = b.Remove([]byte("foo"))
		assert.Nil(err)

		contains, err := a.Contains([]byte("foo"))
		assert.Nil(err)
		assert.True(contains)
	})
}
