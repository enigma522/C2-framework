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

func (m *PingModule) Execute(command string) (string, error) {
	fmt.Println("Executing ping")
	return "PONG", nil
}

