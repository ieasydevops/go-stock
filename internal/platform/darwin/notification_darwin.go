//go:build darwin
// +build darwin

package darwin

import (
	"fmt"
	"log"
	"os/exec"

	"go-stock/internal/platform/interfaces"
)

var _ interfaces.NotificationService = (*notificationService)(nil)

type notificationService struct {
	appName string // Store appName if needed for future use
}

// NewNotificationService creates a new service for macOS notifications.
func NewNotificationService(appName string) interfaces.NotificationService {
	return &notificationService{appName: appName}
}

// SendNotification displays a system notification on macOS using osascript.
func (s *notificationService) SendNotification(config interfaces.NotificationConfig) error {
	// Prepare the osascript command
	script := fmt.Sprintf(`display notification "%s" with title "%s"`, config.Message, config.Title)

	// Add sound if requested
	if config.Sound {
		script += ` sound name "default"`
	}

	// Run the command
	cmd := exec.Command("osascript", "-e", script)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to send notification via osascript: %w", err)
	}

	// For OnClick callbacks, log that it's not directly supported
	if config.OnClick != nil {
		log.Println("macOS NotificationService: OnClick callback is configured but not directly supported by osascript notifications")
		// In a full implementation, consider using a custom URL scheme that the app can register to handle
	}

	return nil
}

// SupportsActions checks if the platform notification system supports actions.
// This simple implementation doesn't support custom actions.
func (s *notificationService) SupportsActions() bool {
	return false // Simple osascript implementation doesn't support actions
}
