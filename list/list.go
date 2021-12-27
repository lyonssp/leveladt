package list

import (
	"encoding/binary"

	"github.com/syndtr/goleveldb/leveldb"
)

type List struct {
	ns  []byte
	ldb *leveldb.DB

	length int
}

// Append the value x to the list
func (ls *List) Append(v []byte) error {
	k := ls.key(int64(ls.length))
	ls.length++
	return ls.ldb.Put(k, v, nil)
}

// Get return the item at index i
func (ls *List) Get(i int64) ([]byte, error) {
	v, err := ls.ldb.Get(ls.key(i), nil)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (ls *List) key(i int64) []byte {
	namespaced := make([]byte, len(ls.ns)+8) // namespace length plus 64 bit integer
	copy(namespaced[:len(ls.ns)], ls.ns)
	binary.PutVarint(namespaced, i)
	return namespaced
}
