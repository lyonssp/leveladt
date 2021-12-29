package queue

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/commands"
	"github.com/leanovate/gopter/gen"
	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"
)

func TestQueueModel(t *testing.T) {
	assert := assert.New(t)

	test := &commands.ProtoCommands{
		NewSystemUnderTestFunc: func(initialState commands.State) commands.SystemUnderTest {
			dir, err := ioutil.TempDir("", "queue-*")
			assert.Nil(err)

			db, err := leveldb.OpenFile(dir, nil)
			assert.Nil(err)

			return NewQueue([]byte("test"), db)
		},
		InitialStateGen: gen.Const(makeQueueModel()),
		InitialPreConditionFunc: func(_ commands.State) bool {
			return true
		},
		GenCommandFunc: func(st commands.State) gopter.Gen {
			return gen.OneGenOf(genPushCommand, genPopCommand(st))
		},
	}

	properties := gopter.NewProperties(gopter.DefaultTestParameters())
	properties.Property("model", commands.Prop(test))
	properties.TestingRun(t)
}

var genPushCommand gopter.Gen = func(params *gopter.GenParameters) *gopter.GenResult {
	return gopter.NewGenResult(
		pushCommand{
			x: []byte(gen.Identifier()(params).Result.(string)),
		},
		gopter.NoShrinker,
	)
}

var genPopCommand = func(st commands.State) gopter.Gen {
	return func(params *gopter.GenParameters) *gopter.GenResult {
		return gopter.NewGenResult(
			popCommand{},
			gopter.NoShrinker,
		)
	}
}

type pushCommand struct {
	x []byte
}

func (cmd pushCommand) Run(sut commands.SystemUnderTest) commands.Result {
	q := sut.(*Queue)
	err := q.Push(cmd.x)
	if err != nil {
		return commands.Result(err)
	}
	return nil
}

func (cmd pushCommand) NextState(state commands.State) commands.State {
	st := state.(queueModel).clone()
	st.Push(cmd.x)
	return st
}

func (cmd pushCommand) PreCondition(_ commands.State) bool {
	return true
}

func (cmd pushCommand) PostCondition(st commands.State, result commands.Result) *gopter.PropResult {
	if e, ok := result.(error); ok {
		return &gopter.PropResult{Error: e}
	}

	return gopter.NewPropResult(true, "")
}

func (cmd pushCommand) String() string {
	return fmt.Sprintf("push(%s)", string(cmd.x))
}

type popCommand struct{}

func (cmd popCommand) Run(sut commands.SystemUnderTest) commands.Result {
	q := sut.(*Queue)
	front, err := q.Pop()
	if err != nil {
		return commands.Result(err)
	}
	return front
}

func (cmd popCommand) NextState(state commands.State) commands.State {
	st := state.(queueModel).clone()
	st.Pop()
	return st
}

func (cmd popCommand) PostCondition(st commands.State, result commands.Result) *gopter.PropResult {
	if e, ok := result.(error); ok {
		return &gopter.PropResult{Error: e}
	}

	got := result.([]byte)
	want := st.(queueModel).lastPopped
	if !bytes.Equal(got, want) {
		return gopter.NewPropResult(false, fmt.Sprintf("%s != %s", got, want))
	}

	return gopter.NewPropResult(true, "")
}

func (cmd popCommand) PreCondition(st commands.State) bool {
	return st.(queueModel).size() > 0
}

func (cmd popCommand) String() string {
	return "pop()"
}

var (
	_ commands.Command = pushCommand{}
	_ commands.Command = popCommand{}
)
