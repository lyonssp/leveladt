package queue

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/arbitrary"
	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"
)

func TestQueueProperties(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSize = 1

	properties := gopter.NewProperties(parameters)

	arbitraries := arbitrary.DefaultArbitraries()
	properties.Property("first appended element is always the result of pop", arbitraries.ForAll(
		func(ss []string) bool {
			dir, err := ioutil.TempDir("", "test")
			if err != nil {
				return false
			}
			db, err := leveldb.OpenFile(dir, nil)
			if err != nil {
				return false
			}

			q := NewQueue([]byte("test"), db)

			for _, s := range ss {
				if err := q.Push([]byte(s)); err != nil {
					return false
				}
			}

			front, err := q.Pop()
			if err != nil {
				return false
			}

			if !bytes.Equal(front, []byte(ss[0])) {
				return false
			}

			return true
		},
	))

	properties.TestingRun(t)
}

func TestQueue(t *testing.T) {
	assert := assert.New(t)

	t.Run("push then pop", func(t *testing.T) {
		dir, err := ioutil.TempDir("", "test")
		assert.Nil(err)

		db, err := leveldb.OpenFile(dir, nil)
		assert.Nil(err)

		q := NewQueue([]byte("test"), db)

		err = q.Push([]byte("foo"))
		assert.Nil(err)

		err = q.Push([]byte("bar"))
		assert.Nil(err)

		got, err := q.Pop()
		assert.Nil(err)
		assert.Equal([]byte("foo"), got)
	})

	t.Run("push duplicates", func(t *testing.T) {
		dir, err := ioutil.TempDir("", "test")
		assert.Nil(err)

		db, err := leveldb.OpenFile(dir, nil)
		assert.Nil(err)

		q := NewQueue([]byte("test"), db)

		err = q.Push([]byte("foo"))
		assert.Nil(err)

		err = q.Push([]byte("foo"))
		assert.Nil(err)

		got, err := q.Pop()
		assert.Nil(err)
		assert.Equal([]byte("foo"), got)

		got, err = q.Pop()
		assert.Nil(err)
		assert.Equal([]byte("foo"), got)
	})
}

func TestNamespacing(t *testing.T) {
	assert := assert.New(t)

	dir, err := ioutil.TempDir("", "test")
	assert.Nil(err)

	db, err := leveldb.OpenFile(dir, nil)
	assert.Nil(err)

	a := NewQueue([]byte("xxx"), db)
	b := NewQueue([]byte("yyy"), db)

	err = a.Push([]byte("foo"))
	assert.Nil(err)

	err = b.Push([]byte("bar"))
	assert.Nil(err)

	front, err := b.Pop()
	assert.Equal([]byte("bar"), front)
	assert.Nil(err)

	front, err = a.Pop()
	assert.Equal([]byte("foo"), front)
	assert.Nil(err)
}

// Capture failed model test sequences
func TestRegressions(t *testing.T) {
	assert := assert.New(t)

	t.Run("regression 0", func(t *testing.T) {
		dir, err := ioutil.TempDir("", "test")
		assert.Nil(err)

		db, err := leveldb.OpenFile(dir, nil)
		assert.Nil(err)

		a := NewQueue([]byte("xxx"), db)

		a.Push([]byte("cz9qanCc"))
		a.Push([]byte("wiekc00p"))
		a.Pop()
		a.Push([]byte("t"))
		a.Pop()
		a.Push([]byte("t"))
		a.Push([]byte("h1lvfxhb"))
		a.Pop()

		front, err := a.Pop()
		assert.Nil(err)
		assert.Equal([]byte("t"), front)

	})
}
