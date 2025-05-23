package interfaces

// Display represents a single display/monitor.
type Display struct {
	ID        string  // Platform-specific display ID
	Name      string  // Human-readable name, if available
	Bounds    Rect    // The bounds of the display in global screen coordinates
	WorkArea  Rect    // The usable area, excluding taskbars, docks, etc.
	Scale     float64 // The scaling factor for this display (e.g., 1.0, 1.5, 2.0)
	IsPrimary bool
}

// Rect represents a rectangle with an origin (X, Y) and dimensions (Width, Height).
type Rect struct {
	X      int
	Y      int
	Width  int
	Height int
}

// ScreenInfoService defines the contract for retrieving information about screen displays.
type ScreenInfoService interface {
	// GetPrimaryDisplay returns information about the primary display.
	GetPrimaryDisplay() (Display, error)

	// GetAllDisplays returns information about all connected displays.
	GetAllDisplays() ([]Display, error)

	// GetDisplayForWindow (if applicable, might require window handle) returns the display a window is on.
	// GetDisplayForWindow(windowHandle uintptr) (Display, error)
}
