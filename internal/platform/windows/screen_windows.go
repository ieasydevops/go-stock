//go:build windows
// +build windows

package windows

import (
	"fmt"
	// "syscall" // No longer needed directly for UPointer
	"unsafe"

	"golang.org/x/sys/windows"

	"go-stock/internal/platform/interfaces"
)

// Constants from winuser.h for GetSystemMetrics
const (
	SM_CXSCREEN = 0
	SM_CYSCREEN = 1
	// SM_CXVIRTUALSCREEN = 78
	// SM_CYVIRTUALSCREEN = 79
	// SM_CMONITORS       = 80
)

// Constants for SystemParametersInfo
const (
	SPI_GETWORKAREA = 0x0030
)

// RECT structure - ensure this matches the Windows API RECT structure (int32)
// The one in interfaces.Rect is int, which might be different size than C.LONG (int32).
// For syscalls, we need the exact Windows RECT.
type winRECT struct {
	Left   int32
	Top    int32
	Right  int32
	Bottom int32
}

var _ interfaces.ScreenInfoService = (*screenInfoService)(nil)

type screenInfoService struct{}

// NewScreenInfoService creates a new service for Windows screen information.
func NewScreenInfoService() interfaces.ScreenInfoService {
	return &screenInfoService{}
}

// GetPrimaryDisplay returns information about the primary display.
func (s *screenInfoService) GetPrimaryDisplay() (interfaces.Display, error) {
	user32 := windows.NewLazySystemDLL("user32.dll")
	getSystemMetrics := user32.NewProc("GetSystemMetrics")
	systemParametersInfoW := user32.NewProc("SystemParametersInfoW")

	cxScreenRaw, _, errGetCx := getSystemMetrics.Call(uintptr(SM_CXSCREEN))
	// According to MSDN, GetSystemMetrics on error might return 0.
	// However, 0 can also be a valid value for some metrics if the primary monitor is not configured.
	// It's better to check the error from Call if it's non-nil, though Call itself rarely returns non-nil error for GetSystemMetrics.
	// A more robust check is to see if the value makes sense or if an extended error occurred.
	// For simplicity, we'll assume if an error is present in `errGetCx` (from the syscall itself, not the function's logical error), it's bad.
	if uintptr(cxScreenRaw) == 0 && errGetCx != nil && errGetCx.Error() != "The operation completed successfully." {
		return interfaces.Display{}, fmt.Errorf("GetSystemMetrics for SM_CXSCREEN failed: %v", errGetCx)
	}
	cxScreen := int(cxScreenRaw)

	cyScreenRaw, _, errGetCy := getSystemMetrics.Call(uintptr(SM_CYSCREEN))
	if uintptr(cyScreenRaw) == 0 && errGetCy != nil && errGetCy.Error() != "The operation completed successfully." {
		return interfaces.Display{}, fmt.Errorf("GetSystemMetrics for SM_CYSCREEN failed: %v", errGetCy)
	}
	cyScreen := int(cyScreenRaw)

	var workAreaRect winRECT
	ret, _, errSP := systemParametersInfoW.Call(
		uintptr(SPI_GETWORKAREA),
		0,                                      // uiParam, not used for SPI_GETWORKAREA
		uintptr(unsafe.Pointer(&workAreaRect)), // pvParam is a pointer to a RECT structure
		0,                                      // fWinIni, not used for SPI_GETWORKAREA
	)
	if ret == 0 { // SystemParametersInfo returns 0 (FALSE) on failure.
		// errSP might contain an error message from the syscall if one occurred.
		return interfaces.Display{}, fmt.Errorf("SystemParametersInfoW for SPI_GETWORKAREA failed. Call error: %v. Rect: %+v", errSP, workAreaRect)
	}

	// DPI scaling is complex. A basic approach might assume 96 DPI (scale 1.0) or try to get it.
	// For a robust solution: GetDpiForMonitor or (older) GetDeviceCaps with LOGPIXELSX.
	// For now, defaulting to 1.0.
	dpiScale := 1.0
	// TODO: Implement proper DPI detection.

	return interfaces.Display{
		ID:     "Primary",
		Name:   "Primary Display",
		Bounds: interfaces.Rect{X: 0, Y: 0, Width: cxScreen, Height: cyScreen},
		WorkArea: interfaces.Rect{
			X:      int(workAreaRect.Left),
			Y:      int(workAreaRect.Top),
			Width:  int(workAreaRect.Right - workAreaRect.Left),
			Height: int(workAreaRect.Bottom - workAreaRect.Top),
		},
		Scale:     dpiScale,
		IsPrimary: true,
	}, nil
}

// GetAllDisplays returns information about all connected displays.
// This requires EnumDisplayMonitors and GetMonitorInfoW, which is more involved.
// For this basic implementation, we will just return the primary display.
func (s *screenInfoService) GetAllDisplays() ([]interfaces.Display, error) {
	// TODO: Implement full multi-monitor support using EnumDisplayMonitors.
	primary, err := s.GetPrimaryDisplay()
	if err != nil {
		return nil, fmt.Errorf("failed to get primary display info for GetAllDisplays: %w", err)
	}
	return []interfaces.Display{primary}, nil
}
