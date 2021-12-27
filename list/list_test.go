package list

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"
)

func TestList(t *testing.T) {
	assert := assert.New(t)

	dir, err := ioutil.TempDir("", "test")
	assert.Nil(err)

	db, err := leveldb.OpenFile(dir, nil)
	assert.Nil(err)

	s := List{
		ns:  []byte("xxx"),
		ldb: db,
	}

	err = s.Append([]byte("foo"))
	assert.Nil(err)

	err = s.Append([]byte("bar"))
	assert.Nil(err)

	v, err := s.Get(0)
	assert.Nil(err)
	assert.Equal("foo", string(v))

	v, err = s.Get(1)
	assert.Nil(err)
	assert.Equal("bar", string(v))
}

func TestNamespacing(t *testing.T) {
	assert := assert.New(t)

	dir, err := ioutil.TempDir("", "test")
	assert.Nil(err)

	db, err := leveldb.OpenFile(dir, nil)
	assert.Nil(err)

	a := List{
		ns:  []byte("xxx"),
		ldb: db,
	}
	b := List{
		ns:  []byte("yyy"),
		ldb: db,
	}

	err = a.Append([]byte("foo"))
	assert.Nil(err)

	err = b.Append([]byte("bar"))
	assert.Nil(err)

	av, err := a.Get(0)
	assert.Nil(err)
	assert.Equal("foo", string(av))

	bv, err := b.Get(0)
	assert.Nil(err)
	assert.Equal("bar", string(bv))
}
