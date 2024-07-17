package per

import (
	"syscall"
	"unsafe"
	"fmt"

)

var (
	modadvapi32          = syscall.NewLazyDLL("advapi32.dll")
	procRegCreateKeyExW  = modadvapi32.NewProc("RegCreateKeyExW")
	procRegSetValueExW   = modadvapi32.NewProc("RegSetValueExW")
	procRegCloseKey      = modadvapi32.NewProc("RegCloseKey")
	modKernel32           = syscall.NewLazyDLL("kernel32.dll")
	modShell32            = syscall.NewLazyDLL("shell32.dll")
	procGetModuleFileName = modKernel32.NewProc("GetModuleFileNameW")
	procSHGetFolderPath   = modShell32.NewProc("SHGetFolderPathW")
	procCopyFile          = modKernel32.NewProc("CopyFileW")
)

const (
	HKEY_CURRENT_USER = 0x80000001
	KEY_ALL_ACCESS    = 0xf003f
	REG_SZ            = 1
	csidlPersonal     = 0x0005 // CSIDL_PERSONAL - My Documents folder
)

func RegCreateKeyEx(key syscall.Handle, subKey string) (syscall.Handle, error) {
	var result syscall.Handle
	var disposition uint32
	ret, _, _ := syscall.SyscallN(procRegCreateKeyExW.Addr(),
		uintptr(key),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(subKey))),
		0,
		uintptr(0),
		0,
		KEY_ALL_ACCESS,
		uintptr(0),
		uintptr(unsafe.Pointer(&result)),
		uintptr(unsafe.Pointer(&disposition)),
	)
	if ret != 0 {
		return 0, fmt.Errorf("RegCreateKeyExW failed with error code: %d", ret)
	}
	return result, nil
}

func RegSetValueEx(key syscall.Handle, valueName string, value string) error {
	ret, _, _ := syscall.SyscallN(procRegSetValueExW.Addr(),
	uintptr(key),
	uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(valueName))),
	0,
	REG_SZ,
	uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(value))),
	uintptr(len(value)*2),
	)
	
	if ret != 0 {
		return fmt.Errorf("RegSetValueExW failed with error code: %d", ret)
	}
	return nil
}

func RegCloseKey(key syscall.Handle) error {
	ret, _, _ := syscall.SyscallN(procRegCloseKey.Addr(),
		uintptr(key),
	)
	if ret != 0 {
		return fmt.Errorf("RegCloseKey failed with error code: %d", ret)
	}
	return nil
}
// File operations

func GetExecutablePath() string {
	var buf [syscall.MAX_PATH]uint16
	_,_, err := procGetModuleFileName.Call(0, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	if err != nil && err.Error() != "The operation completed successfully." {
		fmt.Println("Error getting executable path:", err)
		return ""
	}
	return syscall.UTF16ToString(buf[:])
}

func GetDocumentsPath() string {
	var buf [syscall.MAX_PATH]uint16
	hr, _, _ := procSHGetFolderPath.Call(0, csidlPersonal, 0, 0, uintptr(unsafe.Pointer(&buf[0])))
	if hr != 0 {
		fmt.Println("Error getting Documents path:", hr)
		return ""
	}
	return syscall.UTF16ToString(buf[:])
}

func FileExists(path string) bool {
	attrs, _ := syscall.GetFileAttributes(syscall.StringToUTF16Ptr(path))
	return attrs != syscall.INVALID_FILE_ATTRIBUTES
}

func CopyFile(src, dst string) error {
	srcPtr, err := syscall.UTF16PtrFromString(src)
	if err != nil {
		return err
	}
	dstPtr, err := syscall.UTF16PtrFromString(dst)
	if err != nil {
		return err
	}
	ret, _, err := procCopyFile.Call(uintptr(unsafe.Pointer(srcPtr)), uintptr(unsafe.Pointer(dstPtr)), 0)
	if ret == 0 {
		return err
	}
	return nil
}