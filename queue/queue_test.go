package queue

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"
)

func TestQueueProperties(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	parameters.MinSize = 1 // ensures minimum one element generated in random slices

	properties := gopter.NewProperties(parameters)

	properties.Property("first appended element is always the result of pop", prop.ForAll(
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
				if err := q.Enqueue([]byte(s)); err != nil {
					return false
				}
			}

			front, err := q.Dequeue()
			if err != nil {
				return false
			}

			if !bytes.Equal(front, []byte(ss[0])) {
				return false
			}

			return true
		},
		gen.SliceOf(gen.Identifier()),
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

		err = q.Enqueue([]byte("foo"))
		assert.Nil(err)

		err = q.Enqueue([]byte("bar"))
		assert.Nil(err)

		got, err := q.Dequeue()
		assert.Nil(err)
		assert.Equal([]byte("foo"), got)
	})

	t.Run("push duplicates", func(t *testing.T) {
		dir, err := ioutil.TempDir("", "test")
		assert.Nil(err)

		db, err := leveldb.OpenFile(dir, nil)
		assert.Nil(err)

		q := NewQueue([]byte("test"), db)

		err = q.Enqueue([]byte("foo"))
		assert.Nil(err)

		err = q.Enqueue([]byte("foo"))
		assert.Nil(err)

		got, err := q.Dequeue()
		assert.Nil(err)
		assert.Equal([]byte("foo"), got)

		got, err = q.Dequeue()
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

	err = a.Enqueue([]byte("foo"))
	assert.Nil(err)

	err = b.Enqueue([]byte("bar"))
	assert.Nil(err)

	front, err := b.Dequeue()
	assert.Equal([]byte("bar"), front)
	assert.Nil(err)

	front, err = a.Dequeue()
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

		a.Enqueue([]byte("cz9qanCc"))
		a.Enqueue([]byte("wiekc00p"))
		a.Dequeue()
		a.Enqueue([]byte("t"))
		a.Dequeue()
		a.Enqueue([]byte("t"))
		a.Enqueue([]byte("h1lvfxhb"))
		a.Dequeue()

		front, err := a.Dequeue()
		assert.Nil(err)
		assert.Equal([]byte("t"), front)

	})
}
