package modules

import (
	"fmt"
	"os/exec"
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

	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return "", err
	}

	fmt.Println(string(output))

	return string(output), nil
}
