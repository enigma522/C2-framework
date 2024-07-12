#include <windows.h>
#include <stdio.h>

__attribute__((weak)) void testfunc(char* str) {
    printf("Hello from C!\n");
    printf("You passed: %s\n", str);
}


// Function to save the screenshot to a file
__attribute__((weak)) void SaveBitmapToFile(HBITMAP hBitmap, HDC hDC, int w, int h, const char* filepath) {
    BITMAP bmp;
    BITMAPFILEHEADER bmpFileHeader;
    BITMAPINFOHEADER bmpInfoHeader;
    DWORD dwBmpSize;
    HANDLE hFile;
    DWORD dwBytesWritten;
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

    bmpFileHeader.bfOffBits = sizeof(BITMAPFILEHEADER) + sizeof(BITMAPINFOHEADER);
    bmpFileHeader.bfSize = bmpFileHeader.bfOffBits + dwBmpSize;
    bmpFileHeader.bfType = 0x4D42; // BM

    hFile = CreateFile(filepath, GENERIC_WRITE, 0, NULL, CREATE_ALWAYS, FILE_ATTRIBUTE_NORMAL, NULL);

    WriteFile(hFile, &bmpFileHeader, sizeof(BITMAPFILEHEADER), &dwBytesWritten, NULL);
    WriteFile(hFile, &bmpInfoHeader, sizeof(BITMAPINFOHEADER), &dwBytesWritten, NULL);
    WriteFile(hFile, lpBitmapData, dwBmpSize, &dwBytesWritten, NULL);

    CloseHandle(hFile);
    free(lpBitmapData);
}

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