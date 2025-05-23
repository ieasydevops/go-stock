//go:build windows
// +build windows

package windows

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	// For listing processes, a more robust solution might involve Cgo and PSAPI or Toolhelp32Snapshot.
	// However, for a simpler, Go-native (but less efficient for finding by name) approach,
	// one could use `tasklist` command and parse its output for FindProcessByName.
	// IsProcessRunning can be done by trying to find the process or by os.FindProcess and then sending a signal 0 (though signal 0 is tricky on windows).

	"go-stock/internal/platform/interfaces"
)

var _ interfaces.ProcessManagerService = (*processManagerService)(nil)

type processManagerService struct{}

// NewProcessManagerService creates a new service for Windows process management.
func NewProcessManagerService() interfaces.ProcessManagerService {
	return &processManagerService{}
}

// FindProcessByName finds running processes by name using `tasklist`.
// Returns a list of PIDs or an error.
// This is a basic implementation; a more robust solution would use Windows APIs.
func (s *processManagerService) FindProcessByName(name string) ([]int, error) {
	cmd := exec.Command("tasklist", "/NH", "/FI", fmt.Sprintf("IMAGENAME eq %s", name))
	output, err := cmd.Output()
	if err != nil {
		// tasklist returns error if no process is found, so we check for that specific case.
		// The exact error message or exit code might vary based on system locale.
		if exitErr, ok := err.(*exec.ExitError); ok {
			// A common message when no tasks are running with the specified criteria.
			// This string might need localization or a more robust check of exit codes.
			if strings.Contains(string(exitErr.Stderr), "No tasks are running") ||
				strings.Contains(string(output), "INFO: No tasks are running") { // Check stdout too
				return []int{}, nil // No process found is not an error in this context
			}
		}
		return nil, fmt.Errorf("failed to execute tasklist for %s: %w. Output: %s", name, err, string(output))
	}

	var pids []int
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 2 {
			pidStr := fields[1]
			pid, err := strconv.Atoi(pidStr)
			if err == nil {
				pids = append(pids, pid)
			}
		}
	}
	if len(pids) == 0 && len(lines) > 0 && strings.TrimSpace(lines[0]) != "" {
		// If output was produced but no PIDs parsed, it might be an info message like "No tasks are running..."
		// which wasn't caught by the error check if tasklist exited 0 in some cases.
		if strings.Contains(strings.ToLower(lines[0]), "no tasks") {
			return []int{}, nil
		}
	}

	return pids, nil
}

// IsProcessRunning checks if a process with the given PID is running using `tasklist`.
func (s *processManagerService) IsProcessRunning(pid int) (bool, error) {
	cmd := exec.Command("tasklist", "/NH", "/FI", fmt.Sprintf("PID eq %d", pid))
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if strings.Contains(string(exitErr.Stderr), "No tasks are running") ||
				strings.Contains(string(output), "INFO: No tasks are running") {
				return false, nil
			}
		}
		return false, fmt.Errorf("failed to execute tasklist for PID %d: %w. Output: %s", pid, err, string(output))
	}
	// If output is not empty and no error, the process is running.
	return strings.TrimSpace(string(output)) != "", nil
}

// KillProcess attempts to terminate a process by its PID using `taskkill`.
func (s *processManagerService) KillProcess(pid int) error {
	cmd := exec.Command("taskkill", "/F", "/PID", strconv.Itoa(pid))
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to kill process with PID %d: %w", pid, err)
	}
	return nil
}

// StartProcess starts a new process with the given command and arguments.
func (s *processManagerService) StartProcess(name string, arg ...string) (*exec.Cmd, error) {
	cmd := exec.Command(name, arg...)
	err := cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to start process %s: %w", name, err)
	}
	return cmd, nil
}

// OpenURL opens the given URL in the default web browser.
func (s *processManagerService) OpenURL(url string) error {
	cmd := exec.Command("cmd", "/c", "start", "", url) // The empty "" is for title for start command
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to open URL %s: %w", url, err)
	}
	return nil
}

// OpenFile opens the given file with its default application.
func (s *processManagerService) OpenFile(filePath string) error {
	cmd := exec.Command("cmd", "/c", "start", "", filePath)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	return nil
}

// OpenDirectory opens the given directory in the default file explorer.
func (s *processManagerService) OpenDirectory(dirPath string) error {
	cmd := exec.Command("explorer", dirPath) // Using explorer directly is often better for directories
	err := cmd.Run()
	if err != nil {
		// Fallback to cmd /c start if explorer fails for some reason
		cmdFallback := exec.Command("cmd", "/c", "start", "", dirPath)
		if errFallback := cmdFallback.Run(); errFallback != nil {
			return fmt.Errorf("failed to open directory %s (explorer: %v; cmd start: %v)", dirPath, err, errFallback)
		}
	}
	return nil
}
