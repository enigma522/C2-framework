package winapiV2

import (
	"debug/pe"
	"fmt"
	"unsafe"
	"golang.org/x/sys/windows"
)


type ImageFileHeader struct {
	Machine              uint16
	NumberOfSections     uint16
	TimeDateStamp        uint32
	PointerToSymbolTable uint32
	NumberOfSymbols      uint32
	SizeOfOptionalHeader uint16
	Characteristics      uint16
}

type ImageNTHeaders struct {
	Signature      uint32
	FileHeader     ImageFileHeader
	OptionalHeader pe.OptionalHeader64 // 64 bit.
}

type ImageDosHeader struct {
	EMagic    uint16
	ECblp     uint16
	ECp       uint16
	ECrlc     uint16
	ECparhdr  uint16
	EMinalloc uint16
	EMaxalloc uint16
	ESs       uint16
	ESp       uint16
	ECsum     uint16
	EIp       uint16
	ECs       uint16
	ELfarlc   uint16
	EOvno     uint16
	ERes      [4]uint16
	EOemid    uint16
	EOeminfo  uint16
	ERes2     [10]uint16
	ELfanew   uint32
}

type ImageExportDirectory struct {
	Characteristics       uint32
	TimeDateStamp         uint32
	MajorVersion          uint16
	MinorVersion          uint16
	Name                  uint32
	Base                  uint32
	NumberOfFunctions     uint32
	NumberOfNames         uint32
	AddressOfFunctions    uint32
	AddressOfNames        uint32
	AddressOfNameOrdinals uint32
}




func LoadLibraryA(LibFileName string) (uintptr, error) {
	handle, _, err := loadLibraryA.Call(uintptr(unsafe.Pointer(StringToCharPtr(LibFileName))))
	return handle, err
}

func StringToCharPtr(str string) *uint8 {
	chars := append([]byte(str), 0) // null terminated
	return &chars[0]
}

func GetHashFromString(input string) int {

	hash := 0x1505
	for i := 0; i < len(input); i++ {
		hash = (hash * 0x21) & 0xFFFFFFFF
		hash = (hash + (int(input[i]) & 0xFFFFFFDF)) & 0xFFFFFFFF
	}

	return hash
}

func GetFunctionAddressbyHash(library string, hash int) uintptr {

	// Get library base.
	libraryBase, err := LoadLibraryA(library)
	if err != windows.Errno(0) {
		fmt.Printf("%s\n", err)
	}

	dosHeader := (*ImageDosHeader)(unsafe.Pointer(&(*[64]byte)(unsafe.Pointer(libraryBase))[:][0]))

	offset := (libraryBase) + uintptr(dosHeader.ELfanew)
	imageNTHeaders := (*ImageNTHeaders)(unsafe.Pointer(&(*[264]byte)(unsafe.Pointer(offset))[:][0]))

	exportDirectoryRVA := imageNTHeaders.OptionalHeader.DataDirectory[pe.IMAGE_DIRECTORY_ENTRY_EXPORT].VirtualAddress

	offset = (libraryBase) + uintptr(exportDirectoryRVA)
	imageExportDirectory := (*ImageExportDirectory)(unsafe.Pointer(&(*[256]byte)(unsafe.Pointer(offset))[:][0]))

	offset = (libraryBase) + uintptr(imageExportDirectory.AddressOfFunctions)
	addresOfFunctionsRVA := (*uint)(unsafe.Pointer(&(*[4]byte)(unsafe.Pointer(offset))[:][0]))

	offset = (libraryBase) + uintptr(imageExportDirectory.AddressOfNames)
	addressOfNamesRVA := (*uint32)(unsafe.Pointer(&(*[4]byte)(unsafe.Pointer(offset))[:][0]))

	for i := (0); i < int(imageExportDirectory.NumberOfFunctions); i += 1 {
		
		offset = (libraryBase) + uintptr(imageExportDirectory.AddressOfNames) + uintptr(i*4)
		addressOfNamesRVA = (*uint32)(unsafe.Pointer(&(*[4]byte)(unsafe.Pointer(offset))[:][0]))
		functionNameRVA := (*uint32)(unsafe.Pointer(&(*[4]byte)(unsafe.Pointer(addressOfNamesRVA))[:][0]))
		offset = (libraryBase) + uintptr(*functionNameRVA)
		functionNameVA := (uintptr)(unsafe.Pointer(&(*[32]byte)(unsafe.Pointer(offset))[:][0]))

		// Read until null byte, strings should be null terminated.
		functionName := ""
		for k := 0; k < 1000; k++ { // This for loop should be improved to not have an arbitrary high number.
			nextChar := (*byte)(unsafe.Pointer(&(*[64]byte)(unsafe.Pointer(functionNameVA))[:][k]))
			if *nextChar == 0x00 {
				break
			}
			functionName += string(*nextChar)
		}
		functionNameHash := GetHashFromString(functionName)

		if functionNameHash == hash {

			// addressOfNameOrdinalsRVA[i]
			offset = (libraryBase) + uintptr(imageExportDirectory.AddressOfNameOrdinals) + uintptr(i*2) // We multiply by 2 because each element is 2 bytes in the array.
			ordinalRVA := (*uint16)(unsafe.Pointer(&(*[2]byte)(unsafe.Pointer(offset))[:][0]))

			offset = uintptr(unsafe.Pointer(*&addresOfFunctionsRVA)) + uintptr(uint32(*ordinalRVA)*4) // We multiply by 4 because each element is 4 bytes in the array.
			functionAddressRVA := (*uint32)(unsafe.Pointer(&(*[4]byte)(unsafe.Pointer(offset))[:][0]))

			offset = (libraryBase) + uintptr(*functionAddressRVA) // 0x1b5a0
			functionAddress := (uintptr)(unsafe.Pointer(&(*[4]byte)(unsafe.Pointer(offset))[:][0]))

			return functionAddress
		}
	}

	return 0x00

}
