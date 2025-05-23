package models

import (
	"time"
)

// Settings represents user application settings.
type Settings struct {
	// General settings
	Theme               string `json:"theme"`               // UI theme (light, dark, system)
	Language            string `json:"language"`            // UI language
	ShowPriceInTray     bool   `json:"showPriceInTray"`     // Whether to show stock prices in system tray
	MinimizeToTray      bool   `json:"minimizeToTray"`      // Whether to minimize to system tray
	StartWithSystem     bool   `json:"startWithSystem"`     // Whether to start with system
	CheckUpdatesOnStart bool   `json:"checkUpdatesOnStart"` // Whether to check for updates on start

	// Data settings
	AutoRefreshInterval int    `json:"autoRefreshInterval"` // Data auto-refresh interval in seconds (0 = disabled)
	DefaultExchange     string `json:"defaultExchange"`     // Default stock exchange
	StockDataSource     string `json:"stockDataSource"`     // Source for stock data

	// Notification settings
	EnableNotifications  bool `json:"enableNotifications"`  // Whether to enable notifications
	EnableSoundAlerts    bool `json:"enableSoundAlerts"`    // Whether to enable sound alerts
	EnablePriceAlerts    bool `json:"enablePriceAlerts"`    // Whether to enable price alerts
	ShowNotificationTime int  `json:"showNotificationTime"` // How long to show notifications (in seconds)

	// Export/Import settings
	ExportDirectory string `json:"exportDirectory"` // Default directory for data export

	// Window settings
	WindowWidth  int `json:"windowWidth"`  // Window width
	WindowHeight int `json:"windowHeight"` // Window height
	WindowX      int `json:"windowX"`      // Window X position
	WindowY      int `json:"windowY"`      // Window Y position

	// Metadata
	LastUpdateTime time.Time `json:"lastUpdateTime"` // Last time settings were updated
}

// DefaultSettings returns default application settings
func DefaultSettings() *Settings {
	return &Settings{
		Theme:                "system",
		Language:             "auto",
		ShowPriceInTray:      true,
		MinimizeToTray:       true,
		StartWithSystem:      false,
		CheckUpdatesOnStart:  true,
		AutoRefreshInterval:  60,
		DefaultExchange:      "NASDAQ",
		StockDataSource:      "default",
		EnableNotifications:  true,
		EnableSoundAlerts:    true,
		EnablePriceAlerts:    true,
		ShowNotificationTime: 5,
		ExportDirectory:      "",
		WindowWidth:          1000,
		WindowHeight:         700,
		WindowX:              100,
		WindowY:              100,
		LastUpdateTime:       time.Now(),
	}
}

// ColorScheme represents a custom color scheme for the application UI
type ColorScheme struct {
	ID          string `json:"id"`          // Unique identifier
	Name        string `json:"name"`        // Scheme name
	Description string `json:"description"` // Scheme description
	IsDefault   bool   `json:"isDefault"`   // Whether this is the default scheme

	// Colors
	PrimaryColor    string `json:"primaryColor"`    // Primary color
	SecondaryColor  string `json:"secondaryColor"`  // Secondary color
	BackgroundColor string `json:"backgroundColor"` // Background color
	TextColor       string `json:"textColor"`       // Text color
	PriceUpColor    string `json:"priceUpColor"`    // Color for increasing prices
	PriceDownColor  string `json:"priceDownColor"`  // Color for decreasing prices
	ChartLineColor  string `json:"chartLineColor"`  // Chart line color
	ChartBgColor    string `json:"chartBgColor"`    // Chart background color
	AlertColor      string `json:"alertColor"`      // Alert color

	// Metadata
	CreatedAt time.Time `json:"createdAt"` // When the scheme was created
	UpdatedAt time.Time `json:"updatedAt"` // When the scheme was last updated
}
