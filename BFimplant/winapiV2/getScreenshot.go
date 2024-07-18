package winapiV2

import (
	"errors"
	win "github.com/lxn/win"
	"image"
	"image/png"
	"unsafe"
	"encoding/base64"
	"bytes"
	"fmt"
)

func GetScreenshot() (string, error) {

	var x1, y1, x2, y2, width, height int32

	s, _ := SetProcessDPIAware()
	if !s {
		return "", fmt.Errorf("failed to set process DPI aware")
	}
	x1, _ = GetSystemMetrics(SM_XVIRTUALSCREEN)
	y1, _ = GetSystemMetrics(SM_YVIRTUALSCREEN)
	x2, _ = GetSystemMetrics(SM_CXVIRTUALSCREEN)
	y2, _ = GetSystemMetrics(SM_CYVIRTUALSCREEN)
	width = x2 - x1
	height = y2 - y1

	img := image.NewRGBA(image.Rect(0, 0, int(width), int(height)))

	hdc := win.GetDC(0)
	if hdc == 0 {
		return "", errors.New("GetDC failed")
	}
	defer win.ReleaseDC(0, hdc)

	memory_device := win.CreateCompatibleDC(hdc)
	if memory_device == 0 {
		return "", errors.New("CreateCompatibleDC failed")
	}
	defer win.DeleteDC(memory_device)

	bitmap := win.CreateCompatibleBitmap(hdc, int32(width), int32(height))
	if bitmap == 0 {
		return "", errors.New("CreateCompatibleBitmap failed")
	}
	defer win.DeleteObject(win.HGDIOBJ(bitmap))

	var header win.BITMAPINFOHEADER
	header.BiSize = uint32(unsafe.Sizeof(header))
	header.BiPlanes = 1
	header.BiBitCount = 32
	header.BiWidth = int32(width)
	header.BiHeight = int32(-height)
	header.BiCompression = win.BI_RGB
	header.BiSizeImage = 0


	bitmapDataSize := uintptr(((int64(width)*int64(header.BiBitCount) + 31) / 32) * 4 * int64(height))
	hmem := win.GlobalAlloc(win.GMEM_MOVEABLE, bitmapDataSize)
	defer win.GlobalFree(hmem)
	memptr := win.GlobalLock(hmem)
	defer win.GlobalUnlock(hmem)

	old := win.SelectObject(memory_device, win.HGDIOBJ(bitmap))
	if old == 0 {
		return "", errors.New("SelectObject failed")
	}
	defer win.SelectObject(memory_device, old)

	if !win.BitBlt(memory_device, 0, 0, int32(width), int32(height), hdc, x1, y1, win.SRCCOPY) {
		return "", errors.New("BitBlt failed")
	}

	if win.GetDIBits(hdc, bitmap, 0, uint32(height), (*uint8)(memptr), (*win.BITMAPINFO)(unsafe.Pointer(&header)), win.DIB_RGB_COLORS) == 0 {
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
