//go:build darwin
// +build darwin

package darwin

import (
	"fmt"

	"go-stock/internal/platform/interfaces"
	"github.com/sqweek/dialog"
)

var _ interfaces.DialogService = (*dialogService)(nil)

type dialogService struct{}

// NewDialogService creates a new service for Darwin (macOS) dialogs.
func NewDialogService() interfaces.DialogService {
	return &dialogService{}
}

func (s *dialogService) ShowMessageDialog(dialogType interfaces.DialogType, title, message string, buttons ...string) (string, error) {
	switch dialogType {
	case interfaces.DialogTypeInfo:
		dialog.Message("%s", message).Title(title).Info()
		return "OK", nil
	case interfaces.DialogTypeWarning:
		dialog.Message("%s", message).Title(title).Info()
		return "OK", nil
	case interfaces.DialogTypeError:
		dialog.Message("%s", message).Title(title).Error()
		return "OK", nil
	case interfaces.DialogTypeQuestion:
		if dialog.Message("%s", message).Title(title).YesNo() {
			return "Yes", nil
		} else {
			return "No", nil
		}
	default:
		return "", fmt.Errorf("unsupported dialog type: %v", dialogType)
	}
}

func (s *dialogService) ShowOpenFileDialog(options interfaces.FileDialogOptions) (string, error) {
	builder := dialog.File()
	builder.Title(options.Title)

	// Apply filters if provided
	for _, filter := range options.Filters {
		builder.Filter(filter, filter) // Using filter as both description and pattern
	}

	// Display dialog and handle result
	path, err := builder.Load()
	if err == dialog.Cancelled {
		return "", nil // User cancelled, not an error
	}
	return path, err
}

func (s *dialogService) ShowOpenMultipleFilesDialog(options interfaces.FileDialogOptions) ([]string, error) {
	// As noted, sqweek/dialog does not support opening multiple files directly.
	return nil, fmt.Errorf("opening multiple files is not supported by the current dialog service implementation (sqweek/dialog)")
}

func (s *dialogService) ShowSaveFileDialog(options interfaces.FileDialogOptions) (string, error) {
	builder := dialog.File()
	builder.Title(options.Title)

	// Apply filters if provided
	for _, filter := range options.Filters {
		builder.Filter(filter, filter) // Using filter as both description and pattern
	}

	// Display dialog and handle result
	path, err := builder.Save()
	if err == dialog.Cancelled {
		return "", nil // User cancelled, not an error
	}
	return path, err
}

func (s *dialogService) ShowOpenDirectoryDialog(options interfaces.FileDialogOptions) (string, error) {
	builder := dialog.Directory()
	builder.Title(options.Title)

	path, err := builder.Browse()
	if err == dialog.Cancelled {
		return "", nil // User cancelled
	}
	return path, err
}
