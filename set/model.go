package set

type setModel struct {
	m map[string]struct{}
}

func makeSetModel() setModel {
	return setModel{
		m: make(map[string]struct{}),
	}
}

func (mod *setModel) Add(x []byte) error {
	mod.m[string(x)] = struct{}{}
	return nil
}

func (mod *setModel) Remove(x []byte) error {
	delete(mod.m, string(x))
	return nil
}

func (mod setModel) Contains(x []byte) (bool, error) {
	_, contains := mod.m[string(x)]
	return contains, nil
}

func (mod setModel) clone() setModel {
	cp := makeSetModel()
	for x := range mod.m {
		cp.Add([]byte(x))
	}
	return cp
}

func stringify(x []byte) string {
	return string(x)
}
