//go:build darwin
// +build darwin

package darwin

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"go-stock/internal/platform/interfaces"
)

var _ interfaces.ScreenInfoService = (*screenInfoService)(nil)

type screenInfoService struct{}

// NewScreenInfoService creates a new service for Darwin (macOS) screen information.
func NewScreenInfoService() interfaces.ScreenInfoService {
	return &screenInfoService{}
}

func (s *screenInfoService) GetPrimaryDisplay() (interfaces.Display, error) {
	// On macOS, get the primary display info
	displays, err := s.GetAllDisplays()
	if err != nil {
		return interfaces.Display{}, fmt.Errorf("failed to get primary display: %w", err)
	}

	// Find the primary display
	for _, display := range displays {
		if display.IsPrimary {
			return display, nil
		}
	}

	// If no primary display found, return the first one (if available)
	if len(displays) > 0 {
		return displays[0], nil
	}

	return interfaces.Display{}, fmt.Errorf("no displays found")
}

func (s *screenInfoService) GetAllDisplays() ([]interfaces.Display, error) {
	// On macOS, we can use system_profiler to get display information
	// Alternatively, we could use CGDisplay functions via cgo, but that complicates the build
	cmd := exec.Command("system_profiler", "SPDisplaysDataType", "-json")
	_, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get display information: %w", err)
	}

	// For simplicity in this demo, we'll use a simpler approach with just basic screen metrics
	// In a full implementation, we'd parse the JSON and extract detailed display info

	// This is just a simplified example - in production, parse the JSON output correctly
	var displays []interfaces.Display

	// Create at least one display with reasonable defaults
	// In real code, we would extract this from the system_profiler output
	primaryDisplay := interfaces.Display{
		ID:        "main",
		Name:      "Main Display",
		IsPrimary: true,
		Scale:     1.0,
		Bounds: interfaces.Rect{
			X:      0,
			Y:      0,
			Width:  1440, // Default reasonable values
			Height: 900,  // Adjust based on typical macOS displays
		},
		WorkArea: interfaces.Rect{
			X:      0,
			Y:      23, // Account for menu bar
			Width:  1440,
			Height: 877, // Adjusted for menu bar
		},
	}

	// Get screen resolution using Applescript as a fallback
	if resOutput, err := getScreenResolution(); err == nil {
		parts := strings.Split(resOutput, "x")
		if len(parts) == 2 {
			if width, err := strconv.Atoi(strings.TrimSpace(parts[0])); err == nil {
				primaryDisplay.Bounds.Width = width
				primaryDisplay.WorkArea.Width = width
			}
			if height, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil {
				primaryDisplay.Bounds.Height = height
				primaryDisplay.WorkArea.Height = height - 23 // Adjust for menu bar
			}
		}
	}

	displays = append(displays, primaryDisplay)

	return displays, nil
}

// Helper function to get screen resolution using AppleScript
func getScreenResolution() (string, error) {
	script := `tell application "Finder" to get bounds of window of desktop`
	cmd := exec.Command("osascript", "-e", script)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// Parse bounds output like "0, 0, 1440, 900"
	parts := strings.Split(strings.TrimSpace(string(output)), ", ")
	if len(parts) >= 4 {
		width := parts[2]
		height := parts[3]
		return width + "x" + height, nil
	}

	return "", fmt.Errorf("unexpected output format")
}
