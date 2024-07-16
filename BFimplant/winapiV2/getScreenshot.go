package winapiV2

/*
#cgo LDFLAGS: -lgdi32
#include "hello.c"
*/
import "C"

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"unsafe"
)

func GetScreenshot() (string, error) {
	var x1, y1, x2, y2, w, h int32

	s, _ := SetProcessDPIAware()
	if !s {
		return "", fmt.Errorf("failed to set process DPI aware")
	}
	x1, _ = GetSystemMetrics(SM_XVIRTUALSCREEN)
	y1, _ = GetSystemMetrics(SM_YVIRTUALSCREEN)
	x2, _ = GetSystemMetrics(SM_CXVIRTUALSCREEN)
	y2, _ = GetSystemMetrics(SM_CYVIRTUALSCREEN)
	w = x2 - x1
	h = y2 - y1

	hScreen, _ := GetDC(0)
	defer ReleaseDC(0, hScreen)

	hDC, _ := CreateCompatibleDC(hScreen)
	defer DeleteDC(hDC)

	hBitmap, _ := CreateCompatibleBitmap(hScreen, w, h)
	defer DeleteObject(hBitmap)

	oldObj, _ := SelectObject(hDC, hBitmap)
	defer SelectObject(hDC, oldObj)

	BitBlt(hDC, 0, 0, w, h, hScreen, x1, y1, SRCCOPY)

	var bmp bytes.Buffer
	buffer := C.SaveBitmapToBuffer(C.HBITMAP(unsafe.Pointer(hBitmap)), C.HDC(unsafe.Pointer(hDC)), C.int(w), C.int(h))

	bmp.Write(C.GoBytes(unsafe.Pointer(buffer), C.int(w*h*4)))
	bitmapBytes := bmp.Bytes()

	// Convert bitmap bytes to an image
	img := convertBGRAtoRGBAAndFlip(bitmapBytes, int(w), int(h))

	// Encode image to PNG
	var pngBuffer bytes.Buffer
	err := png.Encode(&pngBuffer, img)
	if err != nil {
		panic(err)
	}

	imgBase64Str := base64.StdEncoding.EncodeToString(pngBuffer.Bytes())

	return imgBase64Str, nil
}

func convertBGRAtoRGBAAndFlip(bgra []byte, width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			flipY := height - y - 1
			j := (flipY*width + x) * 4
			b, g, r, a := bgra[j], bgra[j+1], bgra[j+2], bgra[j+3]
			img.Set(x, y, color.RGBA{R: r, G: g, B: b, A: a})
		}
	}
	return img
}
