package win32

// This file is a go implementation of Windows Version Helper APIs versionhelpers.h
// https://learn.microsoft.com/en-us/windows/win32/sysinfo/version-helper-apis

import (
	"structs"
	"syscall"
	"unsafe"
)

type VerSetTypeMask DWORD

const (
	VER_BUILDNUMBER      VerSetTypeMask = 0x0000004
	VER_MAJORVERSION     VerSetTypeMask = 0x0000002
	VER_MINORVERSION     VerSetTypeMask = 0x0000001
	VER_PLATFORMID       VerSetTypeMask = 0x0000008
	VER_PRODUCT_TYPE     VerSetTypeMask = 0x0000080
	VER_SERVICEPACKMAJOR VerSetTypeMask = 0x0000020
	VER_SERVICEPACKMINOR VerSetTypeMask = 0x0000010
	VER_SUITENAME        VerSetTypeMask = 0x0000040
)

const (
	VER_NT_WORKSTATION       = 0x0000001
	VER_NT_DOMAIN_CONTROLLER = 0x0000002
	VER_NT_SERVER            = 0x0000003

	VER_PLATFORM_WIN32s        = 0
	VER_PLATFORM_WIN32_WINDOWS = 1
	VER_PLATFORM_WIN32_NT      = 2
)

const (
	VER_SERVER_NT                      = 0x80000000
	VER_WORKSTATION_NT                 = 0x40000000
	VER_SUITE_SMALLBUSINESS            = 0x00000001
	VER_SUITE_ENTERPRISE               = 0x00000002
	VER_SUITE_BACKOFFICE               = 0x00000004
	VER_SUITE_COMMUNICATIONS           = 0x00000008
	VER_SUITE_TERMINAL                 = 0x00000010
	VER_SUITE_SMALLBUSINESS_RESTRICTED = 0x00000020
	VER_SUITE_EMBEDDEDNT               = 0x00000040
	VER_SUITE_DATACENTER               = 0x00000080
	VER_SUITE_SINGLEUSERTS             = 0x00000100
	VER_SUITE_PERSONAL                 = 0x00000200
	VER_SUITE_BLADE                    = 0x00000400
	VER_SUITE_EMBEDDED_RESTRICTED      = 0x00000800
	VER_SUITE_SECURITY_APPLIANCE       = 0x00001000
	VER_SUITE_STORAGE_SERVER           = 0x00002000
	VER_SUITE_COMPUTE_SERVER           = 0x00004000
	VER_SUITE_WH_SERVER                = 0x00008000
	VER_SUITE_MULTIUSERTS              = 0x00020000
)

type VerSetCondition BYTE

const (
	VER_EQUAL         VerSetCondition = 1
	VER_GREATER       VerSetCondition = 2
	VER_GREATER_EQUAL VerSetCondition = 3
	VER_LESS          VerSetCondition = 4
	VER_LESS_EQUAL    VerSetCondition = 5
	VER_AND           VerSetCondition = 6
	VER_OR            VerSetCondition = 7
)

var lzVerSetConditionMask = lzKernel32.NewProc("VerSetConditionMask")

func VerSetConditionMask(conditionMask ULONGLONG, typeMask VerSetTypeMask, condition VerSetCondition) ULONGLONG {
	ret, _, _ := lzVerSetConditionMask.Call(uintptr(conditionMask), uintptr(typeMask), uintptr(condition))
	return ULONGLONG(ret)
}

type OSVERSIONINFOEXW struct {
	structs.HostLayout
	OSVersionInfoSize DWORD
	MajorVersion      DWORD
	MinorVersion      DWORD
	BuildNumber       DWORD
	PlatformId        DWORD
	CSDVersion        [128]WCHAR
	ServicePackMajor  WORD
	ServicePackMinor  WORD
	SuiteMask         WORD
	ProductType       BYTE
	Reserved          BYTE
}

const ERROR_OLD_WIN_VERSION = 1150

var lzVerifyVersionInfoW = lzKernel32.NewProc("VerifyVersionInfoW")

func VerifyVersionInfoW(lpVersionInfo *OSVERSIONINFOEXW, typeMask VerSetTypeMask, conditionMask ULONGLONG) (bool, error) {
	ret, _, err := lzVerifyVersionInfoW.Call(uintptr(unsafe.Pointer(lpVersionInfo)), uintptr(typeMask), uintptr(conditionMask))
	if ret == 0 {
		if err == syscall.Errno(ERROR_OLD_WIN_VERSION) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func IsWindowsVersionOrGreater(major, minor, servicePackMajor WORD) (bool, error) {
	var osvi = OSVERSIONINFOEXW{
		OSVersionInfoSize: DWORD(unsafe.Sizeof(OSVERSIONINFOEXW{})),
		MajorVersion:      DWORD(major),
		MinorVersion:      DWORD(minor),
		ServicePackMajor:  servicePackMajor,
	}
	conditionMask := VerSetConditionMask(
		VerSetConditionMask(
			VerSetConditionMask(0, VER_MAJORVERSION, VER_GREATER_EQUAL),
			VER_MINORVERSION, VER_GREATER_EQUAL),
		VER_SERVICEPACKMAJOR, VER_GREATER_EQUAL)
	return VerifyVersionInfoW(&osvi, VER_MAJORVERSION|VER_MINORVERSION|VER_SERVICEPACKMAJOR, conditionMask)
}

const (
	_WIN32_WINNT_NT4          = 0x0400
	_WIN32_WINNT_WIN2K        = 0x0500
	_WIN32_WINNT_WINXP        = 0x0501
	_WIN32_WINNT_WS03         = 0x0502
	_WIN32_WINNT_WIN6         = 0x0600
	_WIN32_WINNT_VISTA        = 0x0600
	_WIN32_WINNT_WS08         = 0x0600
	_WIN32_WINNT_LONGHORN     = 0x0600
	_WIN32_WINNT_WIN7         = 0x0601
	_WIN32_WINNT_WIN8         = 0x0602
	_WIN32_WINNT_WINBLUE      = 0x0603
	_WIN32_WINNT_WINTHRESHOLD = 0x0A00
	_WIN32_WINNT_WIN10        = 0x0A00
)

func IsWindowsXPOrGreater() (bool, error) {
	return IsWindowsVersionOrGreater(WORD(HIBYTE(uint16(_WIN32_WINNT_WINXP))), WORD(LOBYTE(uint16(_WIN32_WINNT_WINXP))), 0)
}

func IsWindowsXPSP1OrGreater() (bool, error) {
	return IsWindowsVersionOrGreater(WORD(HIBYTE(uint16(_WIN32_WINNT_WINXP))), WORD(LOBYTE(uint16(_WIN32_WINNT_WINXP))), 1)
}

func IsWindowsXPSP2OrGreater() (bool, error) {
	return IsWindowsVersionOrGreater(WORD(HIBYTE(uint16(_WIN32_WINNT_WINXP))), WORD(LOBYTE(uint16(_WIN32_WINNT_WINXP))), 2)
}

func IsWindowsXPSP3OrGreater() (bool, error) {
	return IsWindowsVersionOrGreater(WORD(HIBYTE(uint16(_WIN32_WINNT_WINXP))), WORD(LOBYTE(uint16(_WIN32_WINNT_WINXP))), 3)
}

func IsWindowsVistaOrGreater() (bool, error) {
	return IsWindowsVersionOrGreater(WORD(HIBYTE(uint16(_WIN32_WINNT_VISTA))), WORD(LOBYTE(uint16(_WIN32_WINNT_VISTA))), 0)
}

func IsWindowsVistaSP1OrGreater() (bool, error) {
	return IsWindowsVersionOrGreater(WORD(HIBYTE(uint16(_WIN32_WINNT_VISTA))), WORD(LOBYTE(uint16(_WIN32_WINNT_VISTA))), 1)
}

func IsWindowsVistaSP2OrGreater() (bool, error) {
	return IsWindowsVersionOrGreater(WORD(HIBYTE(uint16(_WIN32_WINNT_VISTA))), WORD(LOBYTE(uint16(_WIN32_WINNT_VISTA))), 2)
}

func IsWindows7OrGreater() (bool, error) {
	return IsWindowsVersionOrGreater(WORD(HIBYTE(uint16(_WIN32_WINNT_WIN7))), WORD(LOBYTE(uint16(_WIN32_WINNT_WIN7))), 0)
}

func IsWindows7SP1OrGreater() (bool, error) {
	return IsWindowsVersionOrGreater(WORD(HIBYTE(uint16(_WIN32_WINNT_WIN7))), WORD(LOBYTE(uint16(_WIN32_WINNT_WIN7))), 1)
}

func IsWindows8OrGreater() (bool, error) {
	return IsWindowsVersionOrGreater(WORD(HIBYTE(uint16(_WIN32_WINNT_WIN8))), WORD(LOBYTE(uint16(_WIN32_WINNT_WIN8))), 0)
}

func IsWindows8Point1OrGreater() (bool, error) {
	return IsWindowsVersionOrGreater(WORD(HIBYTE(uint16(_WIN32_WINNT_WINBLUE))), WORD(LOBYTE(uint16(_WIN32_WINNT_WINBLUE))), 0)
}

func IsWindowsThresholdOrGreater() (bool, error) {
	return IsWindowsVersionOrGreater(WORD(HIBYTE(uint16(_WIN32_WINNT_WINTHRESHOLD))), WORD(LOBYTE(uint16(_WIN32_WINNT_WINTHRESHOLD))), 0)
}

func IsWindows10OrGreater() (bool, error) {
	return IsWindowsVersionOrGreater(WORD(HIBYTE(uint16(_WIN32_WINNT_WINTHRESHOLD))), WORD(LOBYTE(uint16(_WIN32_WINNT_WINTHRESHOLD))), 0)
}

func IsWindowsServer() (bool, error) {
	var osvi = OSVERSIONINFOEXW{
		OSVersionInfoSize: DWORD(unsafe.Sizeof(OSVERSIONINFOEXW{})),
		ProductType:       VER_NT_WORKSTATION,
	}
	conditionMask := VerSetConditionMask(0, VER_PRODUCT_TYPE, VER_EQUAL)

	ok, err := VerifyVersionInfoW(&osvi, VER_PRODUCT_TYPE, conditionMask)
	if err != nil {
		return false, err
	}
	return !ok, nil
}

func IsActiveSessionCountLimited() (bool, error) {
	var versionInfo = OSVERSIONINFOEXW{
		OSVersionInfoSize: DWORD(unsafe.Sizeof(OSVERSIONINFOEXW{})),
		SuiteMask:         VER_SUITE_TERMINAL,
	}
	coditionMask := VerSetConditionMask(0, VER_SUITENAME, VER_AND)
	suiteTerminal, err := VerifyVersionInfoW(&versionInfo, VER_SUITENAME, coditionMask)
	if err != nil {
		return false, err
	}

	versionInfo.SuiteMask = VER_SUITE_SINGLEUSERTS
	suiteSingleUserTS, err := VerifyVersionInfoW(&versionInfo, VER_SUITENAME, coditionMask)
	if err != nil {
		return false, err
	}
	return !suiteTerminal || suiteSingleUserTS, nil
}
