package winapiV2

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/png"
	"unsafe"
)

type BITMAPINFOHEADER struct {
	BiSize          uint32
	BiWidth         int32
	BiHeight        int32
	BiPlanes        uint16
	BiBitCount      uint16
	BiCompression   uint32
	BiSizeImage     uint32
	BiXPelsPerMeter int32
	BiYPelsPerMeter int32
	BiClrUsed       uint32
	BiClrImportant  uint32
}

func GetScreenshot() (string, error) {

	var x1, y1, x2, y2, width, height int32

	s, _ := SetProcessDPIAware()
	if !s {
		return "", fmt.Errorf("failllll to settttt proccccc DPI aware")
	}
	x1, _ = GetSystemMetrics(SM_XVIRTUALSCREEN)
	y1, _ = GetSystemMetrics(SM_YVIRTUALSCREEN)
	x2, _ = GetSystemMetrics(SM_CXVIRTUALSCREEN)
	y2, _ = GetSystemMetrics(SM_CYVIRTUALSCREEN)
	width = x2 - x1
	height = y2 - y1

	img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))

	hdc,_ := GetDC(0)
	if hdc == 0 {
		return "", errors.New("GetDC failed")
	}
	defer ReleaseDC(0, hdc)

	memory_device,_ := CreateCompatibleDC(hdc)
	if memory_device == 0 {
		return "", errors.New("CreateCompatibleDC failed")
	}
	defer DeleteDC(memory_device)

	bitmap,_ := CreateCompatibleBitmap(hdc, int32(width), int32(height))
	if bitmap == 0 {
		return "", errors.New("CreateCompatibleBitmap failed")
	}
	defer DeleteObject(bitmap)

	var header BITMAPINFOHEADER
	header.BiSize = uint32(unsafe.Sizeof(header))
	header.BiPlanes = 1
	header.BiBitCount = 32
	header.BiWidth = int32(width)
	header.BiHeight = int32(-height)
	header.BiCompression = BI_RGB
	header.BiSizeImage = 0


	bitmapDataSize := uintptr(((int64(width)*int64(header.BiBitCount) + 31) / 32) * 4 * int64(height))
	hmem,_ := GlobalAlloc(GMEM_MOVEABLE, bitmapDataSize)
	defer GlobalFree(hmem)
	memptr,_ := GlobalLock(hmem)
	defer GlobalUnlock(hmem)

	old,_ := SelectObject(memory_device, bitmap)
	if old == 0 {
		return "", errors.New("SelectObject failed")
	}
	defer SelectObject(memory_device, old)

	succ,_:=BitBlt(memory_device, 0, 0, int32(width), int32(height), hdc, x1, y1, SRCCOPY)

	if !succ {
		return "", errors.New("BitBlt failed")
	}

	ress, _ :=GetDIBits(hdc, bitmap, 0, uint32(height), (*uint8)(unsafe.Pointer(memptr)), (*BITMAPINFO)(unsafe.Pointer(&header)), DIB_RGB_COLORS)
	if ress == 0 {
		return "", errors.New("GetDIBits failed")
	}

	i := 0
	src := uintptr(memptr)
	for y := 0; y < int(height); y++ {
		for x := 0; x < int(width); x++ {
			v0 := *(*uint8)(unsafe.Pointer(src))
			v1 := *(*uint8)(unsafe.Pointer(src + 1))
			v2 := *(*uint8)(unsafe.Pointer(src + 2))

			// BGRA => RGBA, and set A to 255
			img.Pix[i], img.Pix[i+1], img.Pix[i+2], img.Pix[i+3] = v2, v1, v0, 255

			i += 4
			src += 4
		}
	}
	var pngBuffer bytes.Buffer
	err := png.Encode(&pngBuffer, img)
	if err != nil {
		panic(err)
	}

	imgBase64Str := base64.StdEncoding.EncodeToString(pngBuffer.Bytes())

	return imgBase64Str, nil
}
