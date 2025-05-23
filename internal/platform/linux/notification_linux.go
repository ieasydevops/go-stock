//go:build linux
// +build linux

package linux

import (
	"fmt"
	"log"
	"os/exec"

	"go-stock/internal/platform/interfaces"
)

var _ interfaces.NotificationService = (*notificationService)(nil)

type notificationService struct {
	appName string
}

// NewNotificationService creates a new service for Linux notifications.
func NewNotificationService(appName string) interfaces.NotificationService {
	return &notificationService{appName: appName}
}

// SendNotification displays a system notification on Linux using the notify-send command.
// This uses libnotify which is compatible with most Linux desktop environments.
func (s *notificationService) SendNotification(config interfaces.NotificationConfig) error {
	args := []string{
		"--app-name", s.appName,
		config.Title,
		config.Message,
	}

	// Add icon if provided
	if config.IconPath != "" {
		args = append(args, "--icon", config.IconPath)
	}

	// notify-send doesn't have a direct sound option, but we could use a separate
	// command to play a sound if needed
	if config.Sound {
		// For a full implementation, consider spawning a goroutine to play a sound
		// via paplay, aplay, or similar
		log.Println("Linux NotificationService: Sound requested but not directly supported by notify-send")
	}

	// Execute the notify-send command
	cmd := exec.Command("notify-send", args...)
	if err := cmd.Run(); err != nil {
		// Check if notify-send is available
		_, lookErr := exec.LookPath("notify-send")
		if lookErr != nil {
			return fmt.Errorf("notify-send command not found; install libnotify-bin package: %w", lookErr)
		}
		return fmt.Errorf("failed to send notification via notify-send: %w", err)
	}

	// For OnClick callbacks, log that it's not directly supported by simple notify-send
	if config.OnClick != nil {
		log.Println("Linux NotificationService: OnClick callback is configured but not directly supported by basic notify-send")
		// For a full implementation, consider using a more sophisticated D-Bus approach
		// that can handle notification actions and callbacks
	}

	return nil
}

// SupportsActions checks if the platform notification system supports actions.
// This simple implementation doesn't support actions.
func (s *notificationService) SupportsActions() bool {
	return false // Simple notify-send implementation doesn't support actions
}
