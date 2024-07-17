package winapiV2

import (
	"bytes"
	"fmt"
	"syscall"
	"unsafe"
	"encoding/hex"
	"encoding/base64"
)

func Exec(command string) (bytes.Buffer,error){
	// Create pipes for stdout and stderr
	var stdoutRead, stdoutWrite syscall.Handle
	var stderrRead, stderrWrite syscall.Handle

	syscall.CreatePipe(&stdoutRead, &stdoutWrite, &syscall.SecurityAttributes{InheritHandle: 1}, 0)
	syscall.CreatePipe(&stderrRead, &stderrWrite, &syscall.SecurityAttributes{InheritHandle: 1}, 0)

	defer syscall.CloseHandle(stdoutRead)
	defer syscall.CloseHandle(stderrRead)

	// Set up the startup info with the pipes for stdout and stderr
	var startupInfo syscall.StartupInfo
	startupInfo.Cb = uint32(unsafe.Sizeof(startupInfo))
	startupInfo.Flags = syscall.STARTF_USESTDHANDLES
	startupInfo.StdOutput = stdoutWrite
	startupInfo.StdErr = stderrWrite

	var procInfo syscall.ProcessInformation

	// Format the command line correctly
	cmd := "cmd /C "+command;
	cmdLine, err := syscall.UTF16PtrFromString(cmd)

	if err != nil {
		fmt.Println("Failed to convert command line to UTF16:", err)
		return bytes.Buffer{},err
	}

	str1:="C:\\Windows"
	str2,_:=hex.DecodeString("5c5c53797374656d33325c5c")
	str3,_:=base64.StdEncoding.DecodeString("Y21kLmV4ZQ==")


	appName, _ := syscall.UTF16PtrFromString(str1+string(str2)+string(str3))
	fmt.Println("a")
	success, err := CreateProcessW(appName, cmdLine, nil, nil, true, CREATE_NO_WINDOW, nil, nil, &startupInfo, &procInfo)
	if !success || err != nil {
		fmt.Println("CreateProcessW failed:", err)
		return bytes.Buffer{},err
	}
	fmt.Println("b")
	// Close the write end of the pipes to signal EOF to the reading process
	syscall.CloseHandle(stdoutWrite)
	syscall.CloseHandle(stderrWrite)
	var outputBuf bytes.Buffer
	// Wait for the process to complete
	_, err = WaitForSingleObject(syscall.Handle(procInfo.Process), 120000)
	if err != nil {
		fmt.Println("WaitForSingleObject failed:", err)
		return outputBuf,err
	}
	fmt.Println("c")
	// Read the output from stdout
	
	buf := make([]byte, 1024)
	
	n, _ := syscall.Read(stdoutRead, buf)
	if n > 0 {
		outputBuf.Write(buf[:n])
	}
	
	syscall.CloseHandle(syscall.Handle(procInfo.Process))
	syscall.CloseHandle(syscall.Handle(procInfo.Thread))

	
	return outputBuf,nil;
}
