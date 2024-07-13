package winapiV2

/*
#cgo LDFLAGS: -lgdi32
#include "hello.c"
*/
import "C"

import (
	"bytes"
	"fmt"
	"unsafe"
	"strconv"
	"strings"
)


func GetScreenshot(filePath string) (string, error) {
	var x1, y1, x2, y2, w, h int32

	s,_ := SetProcessDPIAware()
	if !s {
		return "", fmt.Errorf("failed to set process DPI aware")
	}
	x1,_ = GetSystemMetrics(SM_XVIRTUALSCREEN)
	y1,_ = GetSystemMetrics(SM_YVIRTUALSCREEN)
	x2,_ = GetSystemMetrics(SM_CXVIRTUALSCREEN)
	y2,_ = GetSystemMetrics(SM_CYVIRTUALSCREEN)
	w = x2 - x1
	h = y2 - y1

	hScreen,_ := GetDC(0)
	defer ReleaseDC(0, hScreen)

	hDC,_ := CreateCompatibleDC(hScreen)
	defer DeleteDC(hDC)

	hBitmap,_ := CreateCompatibleBitmap(hScreen, w, h)
	defer DeleteObject(hBitmap)

	oldObj,_ := SelectObject(hDC, hBitmap)
	defer SelectObject(hDC, oldObj)

	BitBlt(hDC, 0, 0, w, h, hScreen, x1, y1, SRCCOPY)

	fmt.Println("starting")
	
	var bmp bytes.Buffer
	buffer := C.SaveBitmapToBuffer(C.HBITMAP(unsafe.Pointer(hBitmap)), C.HDC(unsafe.Pointer(hDC)), C.int(w), C.int(h))
	
	bmp.Write(C.GoBytes(unsafe.Pointer(buffer), C.int(w*h*4)))
	bitmapBytes := bmp.Bytes()

	var sb strings.Builder

	for _, b := range bitmapBytes {
		sb.WriteString(strconv.Itoa(int(b)))
		sb.WriteString(" ")
	}
	result := sb.String()

	return result, nil
}