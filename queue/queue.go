package queue

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/syndtr/goleveldb/leveldb"
)

// _root is a special key that always points to the front of the queue
const _root = "root"

// Queue is a FIFO queue backed by LevelDB
type Queue struct {
	ns []byte

	// prev is a key that always maps to the next pushed item
	//
	// this key should be referenced when queueing a new item
	// `next` such that a LevelDB mapping is created like so:
	//
	//	_prev -> next
	//
	// TODO: would prefer not to do this, needs to be initialized for any new queue
	_prev []byte

	ldb *leveldb.DB

	l sync.Mutex
}

// Push the value x to the back of the queue
func (ls *Queue) Push(v []byte) error {

	// encode value with namespace and deduplicating nonce
	encoded, err := encode(queueValue{ns: string(ls.ns), val: string(v), nonce: uuid.NewString()})
	if err != nil {
		return err
	}

	// write encoded value and update prev pointer
	prev, err := ls.prev()
	if err != nil {
		return err
	}
	err = ls.ldb.Put(prev, encoded, nil)
	if err != nil {
		return err
	}
	ls._prev = encoded

	return nil
}

// Pop and return the item at the front of the queue
func (ls *Queue) Pop() ([]byte, error) {
	ls.l.Lock()
	defer ls.l.Unlock()

	// get item at front of the queue that will be removed
	frontEncoded, err := ls.peek()
	if err != nil {
		return nil, err
	}
	if frontEncoded == nil {
		return nil, errors.New("cannot pop from empty queue")
	}

	// get the item that will be the new front of the queue
	next, err := ls.get(frontEncoded)
	if err != nil {
		return nil, err
	}

	// start batch updates
	batch := new(leveldb.Batch)

	// include delete for the item at the front of the queue
	batch.Delete(frontEncoded)

	// if there was a second item in the queue, update the root pointer
	if next != nil {
		r, err := ls.root()
		if err != nil {
			return nil, err
		}
		batch.Put(r, next)
	}

	if err := ls.ldb.Write(batch, nil); err != nil {
		return nil, err
	}

	// if we are popping the last item from the queue and need to update the prev pointer
	// this update comes after the batch write in case it fails
	if next == nil {
		r, err := ls.root()
		if err != nil {
			return nil, err
		}
		ls._prev = r
	}

	// decode and parse originally pushed value
	frontDecoded, err := decode(frontEncoded)
	if err != nil {
		return nil, err
	}

	return []byte(frontDecoded.val), nil
}

/*
  convenience accessors that respect the queue namespace
*/
func (ls Queue) get(key []byte) ([]byte, error) {
	front, err := ls.ldb.Get(key, nil)

	if err == leveldb.ErrNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return front, nil
}

func (ls Queue) peek() ([]byte, error) {
	root, err := ls.root()
	if err != nil {
		return nil, err
	}
	return ls.get(root)
}

func (ls Queue) root() ([]byte, error) {
	return encode(queueValue{ns: string(ls.ns), val: _root})
}

func (ls Queue) prev() ([]byte, error) {
	// ensure that prev works for a new instance
	if ls._prev == nil {
		return ls.root()
	}
	return ls._prev, nil
}

// queueValue is a representation of a pushed queue item that can be serialized to bytes
// and helps to support deduplication and namespacing
type queueValue struct {
	ns    string
	val   string
	nonce string
}

func (qv queueValue) MarshalBinary() ([]byte, error) {
	var b bytes.Buffer
	fmt.Fprintf(&b, "%s %s %s", qv.ns, qv.val, qv.nonce)
	return b.Bytes(), nil
}

func (qv *queueValue) UnmarshalBinary(data []byte) error {
	b := bytes.NewBuffer(data)
	_, err := fmt.Fscanf(b, "%s %s %s", &qv.ns, &qv.val, &qv.nonce)
	return err
}

// encode will serialize queueValue instances
func encode(qv queueValue) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(qv); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// decode will deserialize bytes to an instance of queueValue
func decode(encoded []byte) (queueValue, error) {
	dec := gob.NewDecoder(bytes.NewBuffer(encoded))
	var frontDecoded queueValue
	if err := dec.Decode(&frontDecoded); err != nil {
		return queueValue{}, err
	}
	return frontDecoded, nil
}
