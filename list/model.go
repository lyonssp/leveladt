package list

import "fmt"

type listModel struct {
	ls []string
}

func makeListModel() listModel {
	return listModel{
		ls: make([]string, 0),
	}
}

func (mod *listModel) Append(x []byte) error {
	mod.ls = append(mod.ls, string(x))
	return nil
}

func (mod listModel) Get(i int64) ([]byte, error) {
	if i >= int64(len(mod.ls)) {
		return nil, fmt.Errorf("no item found at index %d", i)
	}
	return []byte(mod.ls[i]), nil
}

func (mod listModel) size() int {
	return len(mod.ls)
}

func (mod listModel) clone() listModel {
	cp := makeListModel()
	for _, x := range mod.ls {
		cp.Append([]byte(x))
	}
	return cp
}
