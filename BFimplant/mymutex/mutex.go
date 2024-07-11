package mymutex

import (
	"syscall"
	"unsafe"
)

var (
	kernel32                = syscall.NewLazyDLL("kernel32.dll")
	procCreateMutexW        = kernel32.NewProc("CreateMutexW")
	procReleaseMutex        = kernel32.NewProc("ReleaseMutex")
	procWaitForSingleObject = kernel32.NewProc("WaitForSingleObject")
)

func CreateMutex(name string) (syscall.Handle, error) {
	handle, _, err := procCreateMutexW.Call(
		0,
		0,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(name))),
	)
	if handle == 0 {
		return 0, err
	}
	return syscall.Handle(handle), nil
}

func ReleaseMutex(handle syscall.Handle) error {
	ret, _, err := procReleaseMutex.Call(uintptr(handle))
	if ret == 0 {
		return err
	}
	return nil
}

func WaitForSingleObject(handle syscall.Handle, milliseconds uint32) (uint32, error) {
	ret, _, err := procWaitForSingleObject.Call(
		uintptr(handle),
		uintptr(milliseconds),
	)
	return uint32(ret), err
}
