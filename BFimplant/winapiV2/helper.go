package winapiV2

import (
	"syscall"
	"unsafe"
)

var (
	kernel32                = syscall.NewLazyDLL("kernel32.dll")
	procCreateProcessW      = kernel32.NewProc("CreateProcessW")
	procWaitForSingleObject = kernel32.NewProc("WaitForSingleObject")
	procTerminateProcess =    kernel32.NewProc("TerminateProcess")
	procOpenProcess =         kernel32.NewProc("OpenProcess")
)

const (
	CREATE_NEW_CONSOLE = 0x00000010
	CREATE_NO_WINDOW   = 0x08000000
	PROCESS_ALL_ACCESS  = 0x1F0FFF
	PROCESS_TERMINATE   = 0x0001
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