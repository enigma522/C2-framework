package modules

import (
	"fmt"

)

type PingModule struct{}

func NewPingModule() *PingModule {
	return &PingModule{}
}

func (m *PingModule) Name() string {
	return "ping"
}

func (m *PingModule) Execute(command string, data []byte) (string, error) {
	fmt.Println("Execccccc pinnnnnnn")
	return "PONG", nil
}

