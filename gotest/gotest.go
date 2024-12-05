package main

import (
	"fmt"

	"github.com/adityasaky/gittuf-lua-prototype/common"
	lua "github.com/yuin/gopher-lua"
)

var program = `
function runTests()
	return goTest()
end
`

func main() {
	lState, err := common.NewLuaEnvironment()
	if err != nil {
		panic(err)
	}
	defer lState.Close()

	if err := lState.DoString(program); err != nil {
		panic(err)
	}

	if err := lState.CallByParam(lua.P{
		Fn:      lState.GetGlobal("runTests"),
		NRet:    1,
		Protect: true,
	}, nil); err != nil {
		panic(err)
	}

	ret := lState.Get(-1)
	if ret != nil {
		lState.Pop(1)
		fmt.Println(ret)
	}
}
