package main

import (
	"fmt"

	"github.com/adityasaky/gittuf-lua-prototype/common"
	lua "github.com/yuin/gopher-lua"
)

var program = `
function fileHook(fileContents)
	local lines = splitString(fileContents, "\n")
	for _, line in pairs(lines) do
		if string.len(line) > 80 then
			return false
		end
	end
	return true
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

	tests := []struct {
		contents       string
		expectedReturn bool
	}{
		{
			contents: `Hello, world!
`,
			expectedReturn: true,
		},
		{
			contents: `This line is very long. Like really really long. So long it will trigger the line length check.
`,
			expectedReturn: false,
		},
		{
			contents: `Hello, world!
This line is very long. Like really really long. So long it will trigger the line length check.
`,
			expectedReturn: false,
		},
		{
			contents: `Hello, world!
Hello, world!
Hello, world!
Hello, world!
Hello, world!
`,
			expectedReturn: true,
		},
	}

	for _, test := range tests {
		if err := lState.CallByParam(lua.P{
			Fn:      lState.GetGlobal("fileHook"),
			NRet:    1,
			Protect: true,
		},
			lua.LString(test.contents),
		); err != nil {
			panic(err)
		}
		ret := lState.Get(-1)
		if bool(ret.(lua.LBool)) == test.expectedReturn {
			fmt.Println("Success!")
		} else {
			fmt.Println("Failure!")
		}
	}
}
