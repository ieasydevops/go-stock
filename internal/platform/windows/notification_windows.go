//go:build windows
// +build windows

package windows

import (
	"log"

	"github.com/go-toast/toast"
	"go-stock/internal/platform/interfaces"
)

var _ interfaces.NotificationService = (*notificationService)(nil)

type notificationService struct{}

// NewNotificationService creates a new service for Windows notifications.
func NewNotificationService() interfaces.NotificationService {
	return &notificationService{}
}

// SendNotification displays a system notification on Windows.
func (s *notificationService) SendNotification(config interfaces.NotificationConfig) error {
	notification := toast.Notification{
		AppID:   "Go-Stock", // TODO: Make AppID configurable, perhaps via AppLifecycleService
		Title:   config.Title,
		Message: config.Message,
		Icon:    config.IconPath, // IconPath should be an absolute path to an .ico or .png file
		// TODO: Handle Sound, OnClick (toast actions might be complex)
	}

	if config.Sound {
		notification.Audio = toast.Default
	} else {
		notification.Audio = toast.Silent
	}

	// Example of adding an action (toast library might have limitations or specific ways to handle callbacks)
	// notification.Actions = []toast.Action{
	//  {Type: "protocol", Label: "Open App", Arguments: "go-stock://open"},
	// }

	err := notification.Push()
	if err != nil {
		log.Printf("Error sending windows notification: %v", err)
		return err
	}

	// Handling OnClick is tricky with the current toast library directly.
	// It might require a different approach, like a global handler for toast activations
	// or using a more feature-rich notification library if complex actions are needed.
	// For now, the OnClick from interfaces.NotificationConfig is not directly mapped here.
	if config.OnClick != nil {
		log.Println("Windows NotificationService: OnClick is configured but not directly supported by current toast implementation for custom callbacks. The notification will appear, but the click action might not trigger the specific Go callback without further integration.")
	}

	return nil
}

// SupportsActions checks if the platform notification system supports actions.
// The go-toast library has some support for actions, but complex callback handling is not straightforward.
func (s *notificationService) SupportsActions() bool {
	return true // Tentatively true, but with caveats on callback complexity.
}
