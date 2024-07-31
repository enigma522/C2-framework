package winapiV2

import (
	"fmt"
	"syscall"
	"unsafe"
)

var (
	str =DecryptString("$!3$5,vw")
	procAllocateAndInitializeSid = GetFunctionAddressbyHash(str, 0xa9174a4f)
	procFreeSid                  = GetFunctionAddressbyHash(str, 0xd47b1967)
	procCheckTokenMembership     = GetFunctionAddressbyHash(str, 0x1cf324d0)

	SECURITY_NT_AUTHORITY       = SID_IDENTIFIER_AUTHORITY{Value: [6]byte{0, 0, 0, 0, 0, 5}}
    SECURITY_BUILTIN_DOMAIN_RID = uint32(0x00000020)
	DOMAIN_ALIAS_RID_ADMINS     = uint32(0x00000220)
)

type SID_IDENTIFIER_AUTHORITY struct {
	Value [6]byte
}



func AllocateAndInitializeSid(identAuth *SID_IDENTIFIER_AUTHORITY, subAuthCount byte, subAuth0, subAuth1, subAuth2, subAuth3, subAuth4, subAuth5, subAuth6, subAuth7 uint32) (*syscall.SID, error) {
	var sid *syscall.SID
	r1, _, e1 := syscall.SyscallN(procAllocateAndInitializeSid, uintptr(unsafe.Pointer(identAuth)), uintptr(subAuthCount), uintptr(subAuth0), uintptr(subAuth1), uintptr(subAuth2), uintptr(subAuth3), uintptr(subAuth4), uintptr(subAuth5), uintptr(subAuth6), uintptr(subAuth7), uintptr(unsafe.Pointer(&sid)), 0)
	if r1 == 0 {
		return nil, e1
	}
	return sid, nil
}

func FreeSid(sid *syscall.SID) error {
	r1, _, e1 := syscall.SyscallN(procFreeSid, uintptr(unsafe.Pointer(sid)), 0, 0)
	if r1 != 0 {
		return e1
	}
	return nil
}

func CheckTokenMembership(token syscall.Token, sid *syscall.SID, isMember *bool) error {
	var b int32
	r1, _, e1 := syscall.SyscallN(procCheckTokenMembership, uintptr(token), uintptr(unsafe.Pointer(sid)), uintptr(unsafe.Pointer(&b)))
	*isMember = b != 0
	if r1 == 0 {
		return e1
	}
	return nil
}

// Check if the program has administrative privileges.
func IsAdmin() bool {
	sid, err := AllocateAndInitializeSid(&SECURITY_NT_AUTHORITY, 2, SECURITY_BUILTIN_DOMAIN_RID, DOMAIN_ALIAS_RID_ADMINS, 0, 0, 0, 0, 0, 0)
	if err != nil {
		fmt.Printf("AllocateAndInitializeSid failed: %v\n", err)
		return false
	}
	defer FreeSid(sid)

	var isAdmin bool
	err = CheckTokenMembership(0, sid, &isAdmin)
	if err != nil {
		fmt.Printf("CheckTokenMembership failed: %v\n", err)
		return false
	}

	return isAdmin
}

