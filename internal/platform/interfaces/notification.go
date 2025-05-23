package interfaces

// NotificationConfig holds configuration for a notification.
// TODO: Expand with more options like Icon, Sound, Actions, Timeout.
type NotificationConfig struct {
	Title    string
	Message  string
	IconPath string // Path to an icon file
	Sound    bool   // Whether to play a sound
	// Callback to execute when the notification is clicked.
	// The string argument could be an action identifier if notifications support actions.
	OnClick func(actionID string)
}

// NotificationService defines the contract for sending system notifications.
type NotificationService interface {
	// SendNotification displays a system notification.
	SendNotification(config NotificationConfig) error

	// SupportsActions checks if the platform notification system supports actions.
	SupportsActions() bool
}
