//go:build windows
// +build windows

package windows

import (
	"fmt"

	"go-stock/internal/platform/interfaces"
	"github.com/sqweek/dialog"
)

var _ interfaces.DialogService = (*dialogService)(nil)

type dialogService struct{}

// NewDialogService creates a new service for Windows dialogs.
func NewDialogService() interfaces.DialogService {
	return &dialogService{}
}

func (s *dialogService) ShowMessageDialog(dialogType interfaces.DialogType, title, message string, buttons ...string) (string, error) {
	// sqweek/dialog primarily supports standard message box types (info, warning, error)
	// and doesn't directly support custom button labels through a simple API.
	// It returns bool for Yes/No or Ok/Cancel type dialogs.
	// We need to adapt this. For simplicity, we'll use its basic message types.
	// The `buttons` parameter from the interface is not easily supported here.
	// The returned string will be a best-effort mapping from the boolean result of sqweek/dialog.

	switch dialogType {
	case interfaces.DialogTypeInfo:
		dialog.Message("%s", message).Title(title).Info()
		return "OK", nil // Info typically just has an OK button.
	case interfaces.DialogTypeWarning:
		dialog.Message("%s", message).Title(title).Warning()
		return "OK", nil // Warning also typically has an OK button.
	case interfaces.DialogTypeError:
		dialog.Message("%s", message).Title(title).Error()
		return "OK", nil // Error also typically has an OK button.
	case interfaces.DialogTypeQuestion:
		// Example: Yes/No question. sqweek/dialog provides .YesNo()
		if dialog.Message("%s", message).Title(title).YesNo() {
			return "Yes", nil // Assuming "Yes" was clicked
		} else {
			return "No", nil // Assuming "No" was clicked or dialog was closed
		}
	default:
		return "", fmt.Errorf("unsupported dialog type: %s", dialogType)
	}
	// Note: Error handling from sqweek/dialog itself is minimal for basic messages.
	// It panics on critical graphical server issues, not for user actions.
}

func (s *dialogService) ShowOpenFileDialog(options interfaces.FileDialogOptions) (string, error) {
	builder := dialog.File().Title(options.Title).Load()
	if options.DefaultPath != "" {
		// sqweek/dialog doesn't directly support setting a default start directory for Load.
		// It typically starts in the last used directory or a system default.
	}
	// Adding filters
	for _, pattern := range options.Filters {
		builder.Filter(pattern, pattern) // Label and spec are the same for simplicity here
	}

	file, err := builder.Browse()
	if err == dialog.ErrCancelled {
		return "", nil // User cancelled
	}
	return file, err
}

func (s *dialogService) ShowOpenMultipleFilesDialog(options interfaces.FileDialogOptions) ([]string, error) {
	// sqweek/dialog does not directly support opening multiple files in its current API.
	// It would require a more complex library or platform-specific APIs for this.
	return nil, fmt.Errorf("opening multiple files is not supported by the current dialog service implementation")
}

func (s *dialogService) ShowSaveFileDialog(options interfaces.FileDialogOptions) (string, error) {
	builder := dialog.File().Title(options.Title).Save()
	if options.DefaultPath != "" {
		// Similar to Load, setting a default path/filename might be limited.
		// Often, the Save dialog might take a default filename input if the API supports it.
		// sqweek/dialog appears to take the default name as part of the Save() call if needed, or as a starting point.
	}
	for desc, pattern := range options.Filters {
		builder.Filter(desc, pattern)
	}

	file, err := builder.Browse()
	if err == dialog.ErrCancelled {
		return "", nil // User cancelled
	}
	return file, err
}

func (s *dialogService) ShowOpenDirectoryDialog(options interfaces.FileDialogOptions) (string, error) {
	builder := dialog.Directory().Title(options.Title)
	// DefaultPath for directory dialogs is also typically system-handled.
	dir, err := builder.Browse()
	if err == dialog.ErrCancelled {
		return "", nil // User cancelled
	}
	return dir, err
}
