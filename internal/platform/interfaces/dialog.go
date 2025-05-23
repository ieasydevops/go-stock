package interfaces

// DialogType represents the type of dialog (e.g., info, warning, error, question).
// TODO: Define specific dialog types as constants.
type DialogType string

const (
	DialogTypeInfo     DialogType = "info"
	DialogTypeWarning  DialogType = "warning"
	DialogTypeError    DialogType = "error"
	DialogTypeQuestion DialogType = "question"
)

// FileDialogOptions holds options for file dialogs.
// TODO: Add options like DefaultPath, Filters (e.g., {"Images", "*.png;*.jpg"}).
type FileDialogOptions struct {
	Title         string
	DefaultPath   string
	Filters       map[string]string // e.g., {"Go Files": "*.go", "Text Files": "*.txt"}
	CanCreateDirs bool              // For save dialogs, allow creating directories (macOS)
}

// DialogService defines the contract for showing native system dialogs.
type DialogService interface {
	// ShowMessageDialog displays a simple message dialog (e.g., info, warning, error).
	// Returns true if "OK" or affirmative action was taken, false for "Cancel" or negative.
	ShowMessageDialog(dialogType DialogType, title, message string, buttons ...string) (string, error)

	// ShowOpenFileDialog displays a dialog to open a single file.
	// Returns the selected file path, or an empty string if canceled.
	ShowOpenFileDialog(options FileDialogOptions) (string, error)

	// ShowOpenMultipleFilesDialog displays a dialog to open multiple files.
	// Returns the selected file paths, or an empty slice if canceled.
	ShowOpenMultipleFilesDialog(options FileDialogOptions) ([]string, error)

	// ShowSaveFileDialog displays a dialog to save a file.
	// Returns the selected file path, or an empty string if canceled.
	ShowSaveFileDialog(options FileDialogOptions) (string, error)

	// ShowOpenDirectoryDialog displays a dialog to select a directory.
	// Returns the selected directory path, or an empty string if canceled.
	ShowOpenDirectoryDialog(options FileDialogOptions) (string, error)
}
