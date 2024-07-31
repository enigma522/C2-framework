package mymutex

import (
	"syscall"
	"unsafe"
	"BFimplant/winapiV2"

)

var (
	strr = winapiV2.DecryptString(". 7+ )vw")
	procCreateMutexW        = winapiV2.GetFunctionAddressbyHash(strr, 0x8952e903)
	procReleaseMutex        = winapiV2.GetFunctionAddressbyHash(strr, 0x29af2fd9)
	procWaitForSingleObject = winapiV2.GetFunctionAddressbyHash(strr, 0xdf1b3da)
)

func CreateMutex(name string) (syscall.Handle, error) {
	handle, _, err := syscall.SyscallN(procCreateMutexW,
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
	ret, _, err := syscall.SyscallN(procReleaseMutex,uintptr(handle))
	if ret == 0 {
		return err
	}
	return nil
}

func WaitForSingleObject(handle syscall.Handle, milliseconds uint32) (uint32, error) {
	ret, _, err  :=syscall.SyscallN(procWaitForSingleObject,
		uintptr(handle),
		uintptr(milliseconds),
	)
	
	return uint32(ret), err
}
