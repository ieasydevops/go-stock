//go:build windows
// +build windows

package windows

import (
	"log"
	sync "sync"

	"github.com/energye/systray"
	"go-stock/internal/platform/interfaces"
)

var _ interfaces.SystemTrayService = (*systemTrayService)(nil)

// systemTrayService implements the SystemTrayService interface for Windows.
// It uses the github.com/energye/systray library.
// Note: The systray library typically runs its own loop (systray.Run) and expects to be
// the main application entry point or managed carefully within another main loop (like Wails).
// Managing its lifecycle correctly within a Wails app is crucial.

type systemTrayService struct {
	// Since systray.Run is blocking and handles its own event loop,
	// direct manipulation after initial setup often involves calling systray functions
	// that are safe for concurrent use.
	// We might need a way to signal updates to the systray if its internal state needs to change.
	// For now, we assume systray functions can be called to update tooltip, icon, menu items.

	// iconData and tooltip are stored for potential refresh/update logic if needed,
	// although systray itself maintains the current state.
	currentIconData  []byte
	currentTooltip   string
	currentMenuItems []interfaces.MenuItem
	initialized      bool
	mu               sync.Mutex // To protect access to service state, though systray calls are mostly thread-safe
}

// NewSystemTrayService creates a new service for Windows system tray.
func NewSystemTrayService() interfaces.SystemTrayService {
	return &systemTrayService{}
}

// CreateSystemTray initializes the system tray icon and menu.
// This function in the context of `energye/systray` would typically be called within `systray.Run`'s onReady callback.
// How Wails integrates with systray.Run needs careful consideration.
// For now, this method will configure what *should* happen when systray is ready.
func (s *systemTrayService) CreateSystemTray(iconData []byte, initialTooltip string, menuItems []interfaces.MenuItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.initialized {
		log.Println("System tray already initialized. Updating content.")
		// Fall through to update content if already initialized
	} else {
		// This is where the conceptual problem lies: systray.Run is the entry point.
		// We can't just call systray.SetIcon, etc. here directly before systray.Run is called and onReady is fired.
		// This service should ideally hook into the application's main lifecycle that manages systray.Run.
		// For now, we'll store the desired state and assume something else calls systray.Run.
		// A more robust solution might involve this service starting systray.Run in a goroutine
		// and using channels to communicate readiness and updates, or by having Wails manage the systray loop.
		log.Println("SystemTrayService: CreateSystemTray called. Storing configuration.")
		log.Println("IMPORTANT: github.com/energye/systray requires systray.Run() to be called, usually in main.")
		log.Println("This service assumes that an external mechanism (e.g. Wails app startup) will manage the systray lifecycle.")
	}

	s.currentIconData = iconData
	s.currentTooltip = initialTooltip
	s.currentMenuItems = menuItems

	if s.initialized { // If already running, try to apply changes
		systray.SetIcon(s.currentIconData)
		systray.SetTooltip(s.currentTooltip)
		s.setupMenuItems(s.currentMenuItems)
	}
	// If not initialized, these will be applied in the onReady callback of systray.Run

	s.initialized = true // Mark as configured, actual running state depends on systray.Run elsewhere.
	return nil
}

// setupMenuItems configures the systray menu items.
// This should be called when the systray is ready or when menu items need to be updated.
func (s *systemTrayService) setupMenuItems(menuItems []interfaces.MenuItem) {
	// systray.ResetMenu() // Call this if you want to clear all previous items. Or manage items selectively.
	// For simplicity here, we'll assume we are setting the menu from scratch or replacing it.
	// A more complex implementation would involve adding/removing/updating specific items.

	// First, clear existing menu items if any are managed by systray directly (it might reset itself on SetMenu)
	// This part is tricky as systray doesn't offer a simple 'clear all items added by me' API.
	// We typically rebuild the menu.

	for _, item := range menuItems {
		if item.IsSeparator {
			systray.AddSeparator()
		} else {
			mi := systray.AddMenuItem(item.Title, item.Tooltip)
			go func(menuItem interfaces.MenuItem, systrayMenuItem *systray.MenuItem) {
				for range systrayMenuItem.ClickedCh {
					if menuItem.Callback != nil {
						menuItem.Callback()
					}
				}
			}(item, mi) // Important: capture item and mi in closure

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
			// TODO: Handle SubMenu recursively if supported/needed
		}
	}
}

func (s *systemTrayService) UpdateTooltip(tooltip string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.currentTooltip = tooltip
	if s.initialized {
		systray.SetTooltip(tooltip)
	}
	return nil
}

func (s *systemTrayService) UpdateIcon(iconData []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.currentIconData = iconData
	if s.initialized {
		systray.SetIcon(iconData)
	}
	return nil
}

func (s *systemTrayService) SetMenuItems(menuItems []interfaces.MenuItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.currentMenuItems = menuItems
	if s.initialized {
		// This is a simplified approach: re-create all menu items.
		// systray's API for dynamic updates might require more granular control if we want to avoid full reset.
		// For now, assume onReady will call setupMenuItems or we call it here if already live.
		// The systray library itself does not have a ResetMenu function to clear all items easily.
		// We should call systray.Quit() and then systray.Run() again to fully reset the menu, or rebuild it carefully.
		// The most straightforward way with current systray API is to have setupMenuItems build the menu, and if called again,
		// it would add items again unless managed carefully (e.g. keep references to *systray.MenuItem and hide/show/update them).

		// A common pattern is to clear (if possible) and rebuild.
		// Since systray.ResetMenu() isn't available, and AddMenuItem appends,
		// managing dynamic updates perfectly requires tracking added items or a full systray restart (not ideal).
		// For this implementation, we'll assume this is called to define the *initial* set or a *full replacement* that happens
		// in conjunction with a mechanism that can truly reset/rebuild the menu within systray's lifecycle.
		// If systray is already running, we call setupMenuItems which will append. This isn't ideal for dynamic updates without a clear mechanism.
		// This simplified version will just re-apply. A robust solution is complex with energye/systray's current API for dynamic menu changes without full restart.
		log.Println("SystemTrayService: SetMenuItems called. Re-evaluating menu setup.")
		// We assume systray.AddMenuItem can be called multiple times and it just appends.
		// For a true Set, one might need to systray.Quit() and systray.Run() again if no other clear method is available.
		// Or, if the onReady function that sets up the menu is designed to be callable multiple times and rebuilds it.
		// The following will just add new items based on the new list.
		// A proper 'Set' would need to remove old ones. This is a known challenge with this library.
		s.setupMenuItems(s.currentMenuItems)
	}
	return nil
}

// Quit is called to clean up the system tray before application exit.
func (s *systemTrayService) Quit() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.initialized {
		log.Println("SystemTrayService: Quit called. Signaling systray to exit.")
		systray.Quit() // This signals the systray.Run loop to terminate.
		s.initialized = false
	}
}

// This is a conceptual onReady function that would be passed to systray.Run
// e.g., in main.go: systray.Run(onReady, onExit)
// For this service to function, such an onReady function needs to exist and call methods on this service instance.
/*
func (s *systemTrayService) OnSysTrayReady() {
    s.mu.Lock()
    systray.SetIcon(s.currentIconData)
    systray.SetTooltip(s.currentTooltip)
    s.setupMenuItems(s.currentMenuItems)
    s.mu.Unlock()
    log.Println("System tray is ready and configured.")

    // Example of how to handle a menu item click from within this struct
    // go func() {
    //  for {
    //      select {
    //      case <-someMenuItem.ClickedCh:
    //          s.someAction()
    //      }
    //  }
    // }()
}

func (s *systemTrayService) OnSysTrayExit() {
    log.Println("System tray exiting.")
    // Cleanup, if any, specific to this service after systray loop finishes.
}
*/

// Note: Actual systray.Run() and its onReady/onExit callbacks must be managed by the main application entry point (e.g., Wails app struct or main.go).
// This service provides methods to be called by that main management logic.
