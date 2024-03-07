package main

import (
	"fmt"

	"github.com/adityasaky/gittuf-lua-prototype/common"
	lua "github.com/yuin/gopher-lua"
)

var program = `
function commitMessageHook(commitMessage)
	local lines = splitString(commitMessage, "\n")
	local foundDCO = false
	for _, line in pairs(lines) do
		location = string.find(line, "Signed%-off%-by:", 1)
		if location ~= nil then
			foundDCO = true
			break
		end
	end

	if not foundDCO then
		commitMessage = commitMessage.."\nSigned-off-by: "..repositoryInformation["user.name"].." <"..repositoryInformation["user.email"]..">\n"
	end

	return commitMessage
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
		commitMessage  string
		expectedReturn string
	}{
		{
			commitMessage: `Initial commit
`,
			expectedReturn: `Initial commit

Signed-off-by: Jane Doe <jane.doe@example.com>
`,
		},
		{
			commitMessage: `Initial commit

Signed-off-by: Jane Doe <jane.doe@example.com>
`,
			expectedReturn: `Initial commit

Signed-off-by: Jane Doe <jane.doe@example.com>
`,
		},
		{
			commitMessage: `repository: DRY loading and writing updated root

We had a lot of repeated code for loading, updating, and saving root metadata.
This commit pulls all of that in the repository package into a single
helper.
`,
			expectedReturn: `repository: DRY loading and writing updated root

We had a lot of repeated code for loading, updating, and saving root metadata.
This commit pulls all of that in the repository package into a single
helper.

Signed-off-by: Jane Doe <jane.doe@example.com>
`,
		},
	}

	for _, test := range tests {
		if err := lState.CallByParam(lua.P{
			Fn:      lState.GetGlobal("commitMessageHook"),
			NRet:    1,
			Protect: true,
		},
			lua.LString(test.commitMessage),
		); err != nil {
			panic(err)
		}
		commitMessageRet := lState.Get(-1)
		if string(commitMessageRet.(lua.LString)) == test.expectedReturn {
			fmt.Println("Success!")
		} else {
			fmt.Println("Failure!")
		}
	}
}
