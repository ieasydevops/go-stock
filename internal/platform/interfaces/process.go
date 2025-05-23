package interfaces

import "os/exec"

// ProcessManagerService defines the contract for managing system processes.
type ProcessManagerService interface {
	// FindProcessByName finds running processes by name.
	// Returns a list of PIDs or an error.
	FindProcessByName(name string) ([]int, error)

	// FindProcessByPID finds a running process by its PID.
	// Returns true if found, false otherwise, and an error if the check fails.
	IsProcessRunning(pid int) (bool, error)

	// KillProcess attempts to terminate a process by its PID.
	KillProcess(pid int) error

	// StartProcess starts a new process with the given command and arguments.
	// Returns the started process or an error.
	StartProcess(name string, arg ...string) (*exec.Cmd, error)

	// OpenURL opens the given URL in the default web browser.
	OpenURL(url string) error

	// OpenFile opens the given file with its default application.
	OpenFile(filePath string) error

	// OpenDirectory opens the given directory in the default file explorer.
	OpenDirectory(dirPath string) error
}
