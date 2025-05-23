//go:build linux
// +build linux

package linux

import (
	"log"
	"sync"

	"github.com/getlantern/systray"
	"go-stock/internal/platform/interfaces"
)

var _ interfaces.SystemTrayService = (*systemTrayService)(nil)

// systemTrayService implements the SystemTrayService interface for Linux.
// It uses the github.com/getlantern/systray library.
// Note: The systray library typically runs its own loop (systray.Run) and expects to be
// the main application entry point or managed carefully within another main loop (like Wails).
// Managing its lifecycle correctly within a Wails app is crucial.
type systemTrayService struct {
	currentIconData  []byte
	currentTooltip   string
	currentMenuItems []interfaces.MenuItem
	menuItems        []*systray.MenuItem // Track all menu items for potential later hiding/showing
	initialized      bool
	mu               sync.Mutex
}

// NewSystemTrayService creates a new service for Linux system tray.
func NewSystemTrayService() interfaces.SystemTrayService {
	return &systemTrayService{
		menuItems: make([]*systray.MenuItem, 0),
	}
}

// CreateSystemTray initializes the system tray icon and menu.
// This is conceptually what should happen in systray.Run's onReady callback.
func (s *systemTrayService) CreateSystemTray(iconData []byte, initialTooltip string, menuItems []interfaces.MenuItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.initialized {
		log.Println("Linux SystemTrayService: Already initialized. Updating content.")
	} else {
		log.Println("Linux SystemTrayService: CreateSystemTray called. Storing configuration.")
		log.Println("IMPORTANT: github.com/getlantern/systray requires systray.Run() to be called, usually in main.")
	}

	s.currentIconData = iconData
	s.currentTooltip = initialTooltip
	s.currentMenuItems = menuItems

	if s.initialized { // If already running (i.e., systray.Run is active and onReady was called)
		systray.SetIcon(s.currentIconData)
		systray.SetTooltip(s.currentTooltip)
		s.clearSystrayMenuItems() // Custom helper to attempt to clear existing items
		s.setupSystrayMenuItems(s.currentMenuItems)
	}
	// If not initialized, these will be applied in the onReady callback of systray.Run managed elsewhere.

	s.initialized = true // Mark as configured.
	return nil
}

// clearSystrayMenuItems attempts to remove all items from the systray menu.
func (s *systemTrayService) clearSystrayMenuItems() {
	// The getlantern/systray library doesn't provide a direct way to clear all menu items.
	// The workaround is to hide all existing items.
	for _, item := range s.menuItems {
		item.Hide()
	}
	// Reset the tracking slice
	s.menuItems = make([]*systray.MenuItem, 0)
	log.Println("Linux SystemTrayService: Cleared menu items (hidden)")
}

// setupSystrayMenuItems configures the systray menu items.
func (s *systemTrayService) setupSystrayMenuItems(menuItems []interfaces.MenuItem) {
	for _, item := range menuItems {
		if item.IsSeparator {
			systray.AddSeparator()
			// Note: separator doesn't return an item to track, so we can't hide it later
			// This means menu rebuilding may accumulate separators
		} else {
			mi := systray.AddMenuItem(item.Title, item.Tooltip)
			s.menuItems = append(s.menuItems, mi)

			// Handle clicks in a goroutine as ClickedCh is a channel
			go func(menuItem interfaces.MenuItem, systrayMenuItem *systray.MenuItem) {
				for range systrayMenuItem.ClickedCh {
					if menuItem.Callback != nil {
						menuItem.Callback()
					}
				}
			}(item, mi) // Capture item and mi in closure

			if item.Disabled {
				mi.Disable()
			} else {
				mi.Enable()
			}
			if item.Checked {
				mi.Check()
			} else {
				mi.Uncheck()
			}

			// Handle submenu items if any
			if len(item.SubMenu) > 0 {
				log.Println("Warning: SubMenu items are not fully supported in the current implementation")
				// Future: If systray package develops better submenu support, implement it here
			}
		}
	}
}

func (s *systemTrayService) UpdateTooltip(tooltip string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.currentTooltip = tooltip
	if s.initialized { // Check if systray is actually running, not just configured
		systray.SetTooltip(tooltip)
	}
	return nil
}

func (s *systemTrayService) UpdateIcon(iconData []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.currentIconData = iconData
	if s.initialized {
		systray.SetIcon(iconData) // Ensure iconData is in a format systray expects (e.g., .png bytes)
	}
	return nil
}

func (s *systemTrayService) SetMenuItems(menuItems []interfaces.MenuItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.currentMenuItems = menuItems
	if s.initialized {
		log.Println("Linux SystemTrayService: SetMenuItems called. Rebuilding menu.")
		s.clearSystrayMenuItems()                   // Hide existing items
		s.setupSystrayMenuItems(s.currentMenuItems) // Re-add all items
	}
	return nil
}

func (s *systemTrayService) Quit() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.initialized {
		log.Println("Linux SystemTrayService: Quit called. Signaling systray to exit.")
		systray.Quit() // Signals the systray.Run loop to terminate.
		s.initialized = false
	}
}
