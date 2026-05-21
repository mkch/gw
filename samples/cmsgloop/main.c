#include <windows.h>

#define ID_BTN_SHOW_MESSAGE 1001

static const wchar_t kWindowClassName[] = L"CMsgLoopSampleWindow";

static HMODULE hDll = NULL;

void (*GW_Show)(HWND) = NULL;

LRESULT CALLBACK WndProc(HWND hwnd, UINT msg, WPARAM wParam, LPARAM lParam) {
	switch (msg) {
		case WM_CREATE:
            GW_Show(hwnd);
			break;

		case WM_DESTROY:
			PostQuitMessage(0);
			return 0;
	}

	return DefWindowProcW(hwnd, msg, wParam, lParam);
}

int WINAPI wWinMain(HINSTANCE hInstance, HINSTANCE hPrevInstance, PWSTR pCmdLine, int nCmdShow) {
	(void)hPrevInstance;
	(void)pCmdLine;

    hDll = LoadLibraryW(L"gw.dll");
    if (hDll == NULL) {
        MessageBoxW(NULL, L"Failed to load gw.dll", L"Error", MB_OK | MB_ICONERROR);
        return 1;
    }

    GW_Show = (void (*)(HWND))GetProcAddress(hDll, "Show");
    void (*GW_Cleanup)(void) = (void (*)(void))GetProcAddress(hDll, "Cleanup");
    BOOL (*GW_PreTranslateMessage)(MSG* msg) = (BOOL (*)(MSG*))GetProcAddress(hDll, "PreTranslateMessage");
    if (GW_Show == NULL || GW_Cleanup == NULL || GW_PreTranslateMessage == NULL) {
        MessageBoxW(NULL, L"Failed to get function addresses from gw.dll", L"Error", MB_OK | MB_ICONERROR);
        FreeLibrary(hDll);
        return 1;
    }


	WNDCLASSEXW wc = {0};
	wc.cbSize = sizeof(wc);
	wc.lpfnWndProc = WndProc;
	wc.hInstance = hInstance;
	wc.lpszClassName = kWindowClassName;
	wc.hCursor = LoadCursorW(NULL, IDC_ARROW);
	wc.hbrBackground = (HBRUSH)(COLOR_WINDOW + 1);

	if (!RegisterClassExW(&wc)) {
		return 1;
	}

	HWND hwnd = CreateWindowExW(
		0,
		kWindowClassName,
		L"C Window",
		WS_OVERLAPPEDWINDOW,
		CW_USEDEFAULT,
		CW_USEDEFAULT,
		420,
		260,
		NULL,
		NULL,
		hInstance,
		NULL);

	if (hwnd == NULL) {
		return 1;
	}

	ShowWindow(hwnd, nCmdShow);
	UpdateWindow(hwnd);

	MSG msg;
	while (GetMessageW(&msg, NULL, 0, 0) > 0) {
        if(GW_PreTranslateMessage(&msg)) {
            continue;
        }
		TranslateMessage(&msg);
		DispatchMessageW(&msg);
	}

    GW_Cleanup();
    FreeLibrary(hDll);
	return (int)msg.wParam;
}