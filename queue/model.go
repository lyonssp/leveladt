package queue

import "errors"

type queueModel struct {
	ls         []string
	lastPopped []byte
}

func makeQueueModel() queueModel {
	return queueModel{ls: make([]string, 0)}
}

func (mod *queueModel) Push(x []byte) error {
	mod.ls = append(mod.ls, string(x))
	return nil
}

func (mod *queueModel) Pop() ([]byte, error) {
	if len(mod.ls) <= 0 {
		return nil, errors.New("cannot pop from empty queue")
	}

	front := mod.ls[0]
	mod.lastPopped = make([]byte, len(front))
	copy(mod.lastPopped, front)
	mod.ls = mod.ls[1:]

	return []byte(front), nil
}

func (mod queueModel) size() int {
	return len(mod.ls)
}

func (mod queueModel) clone() queueModel {
	cp := make([]string, len(mod.ls))
	copy(cp, mod.ls)
	return queueModel{ls: cp, lastPopped: mod.lastPopped}
}
