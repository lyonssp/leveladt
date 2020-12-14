package set

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/commands"
	"github.com/leanovate/gopter/gen"
	"github.com/stretchr/testify/assert"
	"github.com/syndtr/goleveldb/leveldb"
)

func TestSetModel(t *testing.T) {
	assert := assert.New(t)

	test := &commands.ProtoCommands{
		NewSystemUnderTestFunc: func(initialState commands.State) commands.SystemUnderTest {
			dir, err := ioutil.TempDir("", "test")
			assert.Nil(err)

			db, err := leveldb.OpenFile(dir, nil)
			assert.Nil(err)

			return &Set{
				ns:  []byte("test"),
				ldb: db,
			}
		},
		InitialStateGen: gen.Const(makeSetModel()),
		InitialPreConditionFunc: func(_ commands.State) bool {
			return true
		},
		GenCommandFunc: func(_ commands.State) gopter.Gen {
			return gen.OneGenOf(genAddCommand, genRemoveCommand, genContainsCommand)
		},
	}

	properties := gopter.NewProperties(gopter.DefaultTestParameters())
	properties.Property("model", commands.Prop(test))
	properties.TestingRun(t)
}

var genAddCommand gopter.Gen = func(params *gopter.GenParameters) *gopter.GenResult {
	return gopter.NewGenResult(
		addCommand{
			x: []byte(gen.Identifier()(params).Result.(string)),
		},
		gopter.NoShrinker,
	)
}

var genRemoveCommand gopter.Gen = func(params *gopter.GenParameters) *gopter.GenResult {
	return gopter.NewGenResult(
		removeCommand{
			x: []byte(gen.Identifier()(params).Result.(string)),
		},
		gopter.NoShrinker,
	)
}

var genContainsCommand gopter.Gen = func(params *gopter.GenParameters) *gopter.GenResult {
	return gopter.NewGenResult(
		containsCommand{
			x: []byte(gen.Identifier()(params).Result.(string)),
		},
		gopter.NoShrinker,
	)
}

type addCommand struct {
	x []byte
}

func (cmd addCommand) Run(sut commands.SystemUnderTest) commands.Result {
	err := sut.(*Set).Add(cmd.x)
	if err != nil {
		return commands.Result(err)
	}
	return nil
}

func (cmd addCommand) NextState(state commands.State) commands.State {
	st := state.(setModel).clone()
	st.Add(cmd.x)
	return st
}

func (cmd addCommand) PreCondition(_ commands.State) bool {
	return true
}

func (cmd addCommand) PostCondition(_ commands.State, result commands.Result) *gopter.PropResult {
	if e, ok := result.(error); ok {
		return &gopter.PropResult{Error: e}
	}
	return gopter.NewPropResult(true, "")
}

func (cmd addCommand) String() string {
	return fmt.Sprintf("add(%s)", string(cmd.x))
}

type removeCommand struct {
	x []byte
}

func (cmd removeCommand) Run(sut commands.SystemUnderTest) commands.Result {
	err := sut.(*Set).Remove(cmd.x)
	if err != nil {
		return commands.Result(err)
	}
	return nil
}

func (cmd removeCommand) NextState(state commands.State) commands.State {
	st := state.(setModel).clone()
	st.Remove(cmd.x)
	return st
}

func (cmd removeCommand) PostCondition(_ commands.State, result commands.Result) *gopter.PropResult {
	if e, ok := result.(error); ok {
		return &gopter.PropResult{Error: e}
	}
	return gopter.NewPropResult(true, "")
}

func (cmd removeCommand) PreCondition(_ commands.State) bool {
	return true
}

func (cmd removeCommand) String() string {
	return fmt.Sprintf("remove(%s)", string(cmd.x))
}

type containsCommand struct {
	x []byte
}

func (cmd containsCommand) NextState(state commands.State) commands.State {
	return state.(setModel).clone()
}

func (cmd containsCommand) PreCondition(_ commands.State) bool {
	return true
}

func (cmd containsCommand) Run(sut commands.SystemUnderTest) commands.Result {
	contains, err := sut.(*Set).Contains(cmd.x)
	if err != nil {
		return commands.Result(err)
	}
	return contains
}

func (cmd containsCommand) PostCondition(state commands.State, result commands.Result) *gopter.PropResult {
	if e, ok := result.(error); ok {
		return &gopter.PropResult{Error: e}
	}
	model, _ := state.(setModel).Contains(cmd.x)
	system := result.(bool)
	return gopter.NewPropResult(system == model, fmt.Sprintf("system != model: %v != %v", system, model))
}

func (cmd containsCommand) String() string {
	return fmt.Sprintf("contains(%s)", string(cmd.x))
}

var (
	_ commands.Command = addCommand{}
	_ commands.Command = removeCommand{}
	_ commands.Command = containsCommand{}
)
