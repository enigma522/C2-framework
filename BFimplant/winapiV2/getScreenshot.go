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
	"os"
	"bufio"
)


func GetScreenshot(filePath string) error {
	var x1, y1, x2, y2, w, h int32

	x1,_ = GetSystemMetrics(SM_XVIRTUALSCREEN)
	y1,_ = GetSystemMetrics(SM_YVIRTUALSCREEN)
	x2,_ = GetSystemMetrics(SM_CXVIRTUALSCREEN)
	y2,_ = GetSystemMetrics(SM_CYVIRTUALSCREEN)
	w = x2 - x1
	h = y2 - y1
	fmt.Println(x1, y1, x2, y2, w, h)

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
	C.testfunc(C.CString("Hello from go"))
	C.SaveBitmapToFile(C.HBITMAP(unsafe.Pointer(hBitmap)), C.HDC(unsafe.Pointer(hDC)),C.int(w),C.int(h),C.CString(filePath))
	// Convert bitmap to base64
	var bmp bytes.Buffer
	buffer := C.SaveBitmapToBuffer(C.HBITMAP(unsafe.Pointer(hBitmap)), C.HDC(unsafe.Pointer(hDC)), C.int(w), C.int(h))
	
	bmp.Write(C.GoBytes(unsafe.Pointer(buffer), C.int(w*h*4)))
	bitmapBytes := bmp.Bytes()
	//base64String := base64.StdEncoding.EncodeToString(bitmapBytes)
	file, err := os.Create("output.txt")
    if err != nil {
        fmt.Println("Error creating file:", err)
        return nil
    }
    defer file.Close()

    // Wrap the file handle with a buffered writer
    writer := bufio.NewWriter(file)

    // Write the base64String to the file
    fmt.Fprintln(writer, bitmapBytes)

    // Flush the buffered writer to ensure all data is written to the file
    writer.Flush()

    fmt.Println("Data written to file successfully.")
	fmt.Println("done")
	fmt.Println("image saved successfully!")
	return nil
}