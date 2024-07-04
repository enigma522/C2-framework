package modules

import (
	"encoding/json"
	"fmt"
	"os"
)

type UploadModule struct{}

func NewUploadModule() *UploadModule {
	return &UploadModule{}
}

func (m *UploadModule) Name() string {
	return "Upload"
}

func (m *UploadModule) Execute(filePath string,data []byte) (string, error) {
	b, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return "", err
	}
	osInfo := map[string]interface{}{
		"file_path": filePath,
		"file_size": len(b),
		"file_data": b,
	}
	file_data, err := json.Marshal(osInfo)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return "", err
	}

	return string(file_data), nil
}
