package repositories

import (
	"context"

	"go-stock/internal/domain/models"
)

// SettingsRepository defines the interface for settings data access.
type SettingsRepository interface {
	// Get retrieves the current user settings.
	Get(ctx context.Context) (*models.Settings, error)

	// Save saves the user settings.
	Save(ctx context.Context, settings *models.Settings) error

	// Reset resets the settings to default values.
	Reset(ctx context.Context) error

	// GetSetting retrieves a specific setting by key.
	GetSetting(ctx context.Context, key string) (interface{}, error)

	// SetSetting sets a specific setting by key.
	SetSetting(ctx context.Context, key string, value interface{}) error
}

// ColorSchemeRepository defines the interface for color scheme data access.
type ColorSchemeRepository interface {
	// Create creates a new color scheme.
	Create(ctx context.Context, scheme *models.ColorScheme) error

	// Update updates a color scheme.
	Update(ctx context.Context, scheme *models.ColorScheme) error

	// Delete deletes a color scheme.
	Delete(ctx context.Context, schemeID string) error

	// GetByID retrieves a color scheme by ID.
	GetByID(ctx context.Context, schemeID string) (*models.ColorScheme, error)

	// GetDefault retrieves the default color scheme.
	GetDefault(ctx context.Context) (*models.ColorScheme, error)

	// ListAll lists all available color schemes.
	ListAll(ctx context.Context) ([]*models.ColorScheme, error)

	// SetDefault sets a color scheme as the default.
	SetDefault(ctx context.Context, schemeID string) error
}
