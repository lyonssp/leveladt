package list

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/commands"
	"github.com/leanovate/gopter/gen"
	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"
)

func TestListModel(t *testing.T) {
	assert := assert.New(t)

	test := &commands.ProtoCommands{
		NewSystemUnderTestFunc: func(initialState commands.State) commands.SystemUnderTest {
			dir, err := ioutil.TempDir("", "test")
			assert.Nil(err)

			db, err := leveldb.OpenFile(dir, nil)
			assert.Nil(err)

			return List{
				ns:  []byte("test"),
				ldb: db,
			}
		},
		InitialStateGen: gen.Const(makeListModel()),
		InitialPreConditionFunc: func(_ commands.State) bool {
			return true
		},
		GenCommandFunc: func(st commands.State) gopter.Gen {
			return gen.OneGenOf(genAppendCommand, genGetCommand(st))
		},
	}

	properties := gopter.NewProperties(gopter.DefaultTestParameters())
	properties.Property("model", commands.Prop(test))
	properties.TestingRun(t)
}

var genAppendCommand gopter.Gen = func(params *gopter.GenParameters) *gopter.GenResult {
	return gopter.NewGenResult(
		appendCommand{
			x: []byte(gen.Identifier()(params).Result.(string)),
		},
		gopter.NoShrinker,
	)
}

var genGetCommand = func(st commands.State) gopter.Gen {
	return func(params *gopter.GenParameters) *gopter.GenResult {
		length := int64(st.(listModel).size())

		if length == 0 {
			return gopter.NewEmptyResult(reflect.TypeOf(getCommand{}))
		}

		index := params.Rng.Int63n(length)

		return gopter.NewGenResult(
			getCommand{i: index},
			gopter.NoShrinker,
		)
	}
}

type appendCommand struct {
	x []byte
}

func (cmd appendCommand) Run(sut commands.SystemUnderTest) commands.Result {
	ls := sut.(List)
	err := ls.Append(cmd.x)
	if err != nil {
		return commands.Result(err)
	}
	return nil
}

func (cmd appendCommand) NextState(state commands.State) commands.State {
	st := state.(listModel).clone()
	st.Append(cmd.x)
	return st
}

func (cmd appendCommand) PreCondition(_ commands.State) bool {
	return true
}

func (cmd appendCommand) PostCondition(_ commands.State, result commands.Result) *gopter.PropResult {
	if e, ok := result.(error); ok {
		return &gopter.PropResult{Error: e}
	}
	return gopter.NewPropResult(true, "")
}

func (cmd appendCommand) String() string {
	return fmt.Sprintf("append(%s)", string(cmd.x))
}

type getCommand struct {
	i int64
}

func (cmd getCommand) Run(sut commands.SystemUnderTest) commands.Result {
	ls := sut.(List)
	_, err := ls.Get(cmd.i)
	if err != nil {
		return commands.Result(err)
	}
	return nil
}

func (cmd getCommand) NextState(state commands.State) commands.State {
	st := state.(listModel).clone()
	st.Get(cmd.i)
	return st
}

func (cmd getCommand) PostCondition(_ commands.State, result commands.Result) *gopter.PropResult {
	if e, ok := result.(error); ok {
		return &gopter.PropResult{Error: e}
	}
	return gopter.NewPropResult(true, "")
}

func (cmd getCommand) PreCondition(_ commands.State) bool {
	return true
}

func (cmd getCommand) String() string {
	return fmt.Sprintf("get(%d)", cmd.i)
}

var (
	_ commands.Command = appendCommand{}
	_ commands.Command = getCommand{}
)
