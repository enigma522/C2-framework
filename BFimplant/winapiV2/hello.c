#include <windows.h>
#include <stdio.h>


// Function to save the screenshot to a buffer
__attribute__((weak)) char* SaveBitmapToBuffer(HBITMAP hBitmap, HDC hDC, int w, int h) {
    BITMAP bmp;
    BITMAPINFOHEADER bmpInfoHeader;
    DWORD dwBmpSize;
    char *lpBitmapData;

    GetObject(hBitmap, sizeof(BITMAP), &bmp);

    bmpInfoHeader.biSize = sizeof(BITMAPINFOHEADER);
    bmpInfoHeader.biWidth = bmp.bmWidth;
    bmpInfoHeader.biHeight = bmp.bmHeight;
    bmpInfoHeader.biPlanes = 1;
    bmpInfoHeader.biBitCount = 32;
    bmpInfoHeader.biCompression = BI_RGB;
    bmpInfoHeader.biSizeImage = 0;
    bmpInfoHeader.biXPelsPerMeter = 0;
    bmpInfoHeader.biYPelsPerMeter = 0;
    bmpInfoHeader.biClrUsed = 0;
    bmpInfoHeader.biClrImportant = 0;

    dwBmpSize = ((bmp.bmWidth * bmpInfoHeader.biBitCount + 31) / 32) * 4 * bmp.bmHeight;

    lpBitmapData = (char*)malloc(dwBmpSize);
    GetDIBits(hDC, hBitmap, 0, (UINT)bmp.bmHeight, lpBitmapData, (BITMAPINFO*)&bmpInfoHeader, DIB_RGB_COLORS);

    // Copy bitmap data to buffer
    return lpBitmapData;
}