package modules

import (
	"fmt"
)

type ScreenshotModule struct{}

func NewScreenshotModule() *ScreenshotModule {
	return &ScreenshotModule{}
}

func (m *ScreenshotModule) Name() string {
	return "screenshot"
}

func (m *ScreenshotModule) Execute(command string) (string, error) {
	
	return fmt.Sprintf("Screenshot saved to %s", "m"), nil
}
