package common

import (
	"os/exec"

	lua "github.com/yuin/gopher-lua"
)

var allGoHelpers = map[string]lua.LGFunction{
	// Go Development Helpers
	"goTest": GoTest,
	"goFmt":  GoFmt,
}

func GoTest(lState *lua.LState) int {
	cmd := exec.Command("go", "test", "./...")
	if output, err := cmd.CombinedOutput(); err != nil {
		lState.Push(lua.LString(string(output)))
		return 1
	}

	return 0
}

func GoFmt(lState *lua.LState) int {
	cmd := exec.Command("go", "fmt")
	if output, err := cmd.CombinedOutput(); err != nil {
		lState.Push(lua.LString(string(output)))
		return 1
	}
	return 0
}
