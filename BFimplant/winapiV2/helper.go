package winapiV2

import (
	"syscall"
	"unsafe"
)

var (
	kernel32                   = syscall.NewLazyDLL("kernel32.dll")
	user32               	   = syscall.NewLazyDLL("User32.dll")
	wingdi					   = syscall.NewLazyDLL("Gdi32.dll")
	procCreateProcessW         = kernel32.NewProc("CreateProcessW")
	procWaitForSingleObject    = kernel32.NewProc("WaitForSingleObject")
	procTerminateProcess 	   = kernel32.NewProc("TerminateProcess")
	procOpenProcess 		   = kernel32.NewProc("OpenProcess")
	procGetSystemMetrics 	   = user32.NewProc("GetSystemMetrics")
	procGetDC 				   = user32.NewProc("GetDC")
	procCreateCompatibleDC     = wingdi.NewProc("CreateCompatibleDC")
	procCreateCompatibleBitmap = wingdi.NewProc("CreateCompatibleBitmap")
	procSelectObject 		   = wingdi.NewProc("SelectObject")
	procBitBlt 				   = wingdi.NewProc("BitBlt")
	procDeleteDC 			   = wingdi.NewProc("DeleteDC")
	procReleaseDC 			   = user32.NewProc("ReleaseDC")
	procDeleteObject 		   = wingdi.NewProc("DeleteObject")
	procGetDIBits			   = wingdi.NewProc("GetDIBits")
	procGetObject			   = wingdi.NewProc("GetObjectW")
	
)

const (
	CREATE_NEW_CONSOLE = 0x00000010
	CREATE_NO_WINDOW   = 0x08000000
	PROCESS_ALL_ACCESS  = 0x1F0FFF
	PROCESS_TERMINATE   = 0x0001
	SM_XVIRTUALSCREEN   = 76
	SM_YVIRTUALSCREEN   = 77
	SM_CXVIRTUALSCREEN  = 78
	SM_CYVIRTUALSCREEN  = 79
	SRCCOPY			 = 0x00CC0020
	DIB_RGB_COLORS 	 = 0
	BI_RGB			 = 0	
)

func CreateProcessW(appName *uint16, cmdLine *uint16, procAttrs *syscall.SecurityAttributes, threadAttrs *syscall.SecurityAttributes, inheritHandles bool, creationFlags uint32, env *uint16, currentDir *uint16, startupInfo *syscall.StartupInfo, procInfo *syscall.ProcessInformation) (bool, error) {
	r1, _, e1 := syscall.SyscallN(procCreateProcessW.Addr(),
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
	r1, _, e1 := syscall.SyscallN(procWaitForSingleObject.Addr(),
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
	r1, _, e1 := syscall.SyscallN(procOpenProcess.Addr(),
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
	r1, _, e1 := syscall.SyscallN(procTerminateProcess.Addr(),
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
	r1, _, e1 := syscall.SyscallN(procGetSystemMetrics.Addr(),
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
	r1, _, e1 := syscall.SyscallN(procGetDC.Addr(),
		uintptr(hwnd),
	)
	if r1 == 0 {
		return 0, e1
	}
	return syscall.Handle(r1), nil
}

func CreateCompatibleDC(hdc syscall.Handle) (syscall.Handle, error) {
	r1, _, e1 := syscall.SyscallN(procCreateCompatibleDC.Addr(),
		uintptr(hdc),
	)
	if r1 == 0 {
		return 0, e1
	}
	return syscall.Handle(r1), nil
}

func CreateCompatibleBitmap(hdc syscall.Handle, width, height int32) (syscall.Handle, error) {
	r1, _, e1 := syscall.SyscallN(procCreateCompatibleBitmap.Addr(),
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
	r1, _, e1 := syscall.SyscallN(procSelectObject.Addr(),
		uintptr(hdc),
		uintptr(hgdiobj),
	)
	if r1 == 0 {
		return 0, e1
	}
	return syscall.Handle(r1), nil
}

func BitBlt(hdc syscall.Handle, xDest, yDest, width, height int32, hdcSrc syscall.Handle, xSrc, ySrc int32, rop uint32) (bool, error) {
	r1, _, e1 := syscall.SyscallN(procBitBlt.Addr(),
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
	r1, _, e1 := syscall.SyscallN(procDeleteDC.Addr(),
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
	r1, _, e1 := syscall.SyscallN(procReleaseDC.Addr(),
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
	r1, _, e1 := syscall.SyscallN(procDeleteObject.Addr(),
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
