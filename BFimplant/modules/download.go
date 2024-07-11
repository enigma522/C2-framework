package modules

import (
    "fmt"
    "os"
)

type DownloadModule struct{}

func NewDownloadModule() *DownloadModule {
    return &DownloadModule{}
}

func (m *DownloadModule) Name() string {
    return "screenshot"
}

func (m *DownloadModule) Execute(filepath string, data []byte) (string, error) {
    file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0644)
    if err != nil {
        fmt.Println("Error opening file:", err)
        return "", err
    }
    defer file.Close()

    // Write data to file
    _, err = file.Write(data)
    if err != nil {
        fmt.Println("Error writing to file:", err)
        return "", err
    }

    fmt.Println("File written successfully with path:", filepath)
    return fmt.Sprintf("File written successfully with path: %s", filepath), nil
}
