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

// special key that always points to the pFront of the queue
const (
	pFront = "front"
	pBack  = "back"
)

// Queue is a FIFO queue backed by LevelDB
type Queue struct {
	ns  []byte
	ldb *leveldb.DB
	l   sync.Mutex
}

func NewQueue(ns []byte, ldb *leveldb.DB) *Queue {
	return &Queue{
		ns:  ns,
		ldb: ldb,
	}
}

// Push the value x to the back of the queue
func (ls *Queue) Push(v []byte) error {
	ls.l.Lock()
	defer ls.l.Unlock()

	// encode input with namespace and deduplicating nonce
	encoded, err := encode(queueValue{ns: string(ls.ns), val: string(v), nonce: uuid.NewString()})
	if err != nil {
		return err
	}

	// get encoded queueValue at back of queue
	encBack, err := ls.peekBack()
	if err != nil {
		return err
	}

	batch := new(leveldb.Batch)

	// if the back pointer points to nothing, this is the first write to the queue and the front pointer must be updated
	if encBack == nil {
		batch.Put(ls.pFront(), encoded)
	}
	batch.Put(encBack, encoded)
	batch.Put(ls.pBack(), encoded)
	return ls.ldb.Write(batch, nil)
}

// Pop and return the item at the front of the queue
func (ls *Queue) Pop() ([]byte, error) {
	ls.l.Lock()
	defer ls.l.Unlock()

	// get encoded queueValue at front of the queue that will be removed
	encFront, err := ls.peek()
	if err != nil {
		return nil, err
	}
	if encFront == nil {
		return nil, errors.New("cannot pop from empty queue")
	}

	// get the item that will be the new front of the queue
	newEncFront, err := ls.get(encFront)
	if err != nil {
		return nil, err
	}

	// start batch updates
	batch := new(leveldb.Batch)

	// include delete for the item at the front of the queue
	batch.Delete(encFront)

	// if there was a second item in the queue, update the front pointer
	// otherwise, we need to update the back pointer to point to front
	if newEncFront != nil {
		batch.Put(ls.pFront(), newEncFront)
	} else {
		batch.Put(ls.pBack(), ls.pFront())
	}

	if err := ls.ldb.Write(batch, nil); err != nil {
		return nil, err
	}

	// decode and parse originally pushed value
	frontDecoded, err := decode(encFront)
	if err != nil {
		return nil, err
	}

	return []byte(frontDecoded.val), nil
}

/*
  convenience accessors that respect the queue namespace
*/
func (ls *Queue) get(key []byte) ([]byte, error) {
	front, err := ls.ldb.Get(key, nil)

	if err == leveldb.ErrNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return front, nil
}

// pFront encodes the pFront constant, respecting the namespace of the queue
func (ls *Queue) pFront() []byte {
	var b bytes.Buffer
	fmt.Fprint(&b, string(ls.ns), pFront)
	return b.Bytes()
}

// pBack encodes the pBack constant, respecting the namespace of the queue
func (ls *Queue) pBack() []byte {
	var b bytes.Buffer
	fmt.Fprint(&b, string(ls.ns), pBack)
	return b.Bytes()
}

// peek returns the encoded queueValue at the front of the queue
func (ls *Queue) peek() ([]byte, error) {
	return ls.get(ls.pFront())
}

// peekBack returns the encoded queueValue at the back of the queue
func (ls *Queue) peekBack() ([]byte, error) {
	return ls.get(ls.pBack())
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
