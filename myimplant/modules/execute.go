package modules

import (
	"fmt"
)

type ExecuteModule struct{}

func NewExecuteModule() *ExecuteModule {
	return &ExecuteModule{}
}

func (m *ExecuteModule) Name() string {
	return "exec"
}

func (m *ExecuteModule) Execute(command string) (string, error) {
	fmt.Println("Executing command:", command)
	return "",nil
}
