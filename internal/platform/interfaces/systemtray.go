package interfaces

// MenuItem represents an item in the system tray menu.
type MenuItem struct {
	Title       string
	Tooltip     string
	Disabled    bool
	Checked     bool
	IsSeparator bool
	SubMenu     []MenuItem // For nested menus
	// Callback is executed when the menu item is clicked.
	// It might not be applicable if SubMenu is present.
	Callback func()
	// TODO: Add Icon field for menu items
}

// SystemTrayService defines the contract for managing a system tray icon and menu.
type SystemTrayService interface {
	// CreateSystemTray initializes the system tray icon and menu.
	// iconData should be the raw bytes of the icon image (e.g., PNG).
	// initialTooltip is the tooltip shown on hover.
	// menuItems define the initial menu structure.
	CreateSystemTray(iconData []byte, initialTooltip string, menuItems []MenuItem) error

	// UpdateTooltip updates the tooltip text for the system tray icon.
	UpdateTooltip(tooltip string) error

	// UpdateIcon updates the system tray icon.
	// iconData should be the raw bytes of the new icon image.
	UpdateIcon(iconData []byte) error

	// SetMenuItems completely replaces the current menu items with the new ones.
	SetMenuItems(menuItems []MenuItem) error

	// AddMenuItem adds a new item to the menu. (Consider if this is needed if SetMenuItems is comprehensive)
	// AddMenuItem(item MenuItem) error

	// UpdateMenuItem allows modifying an existing menu item (e.g., by title or a unique ID).
	// UpdateMenuItem(identifier string, updatedItem MenuItem) error

	// RemoveMenuItem removes a menu item. (Consider if this is needed)
	// RemoveMenuItem(identifier string) error

	// Quit is called to clean up the system tray before application exit.
	Quit()
}
