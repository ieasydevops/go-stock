//go:build darwin
// +build darwin

package darwin

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"go-stock/internal/platform/interfaces"
)

var _ interfaces.ProcessManagerService = (*processManagerService)(nil)

type processManagerService struct{}

// NewProcessManagerService creates a new service for Darwin (macOS) process management.
func NewProcessManagerService() interfaces.ProcessManagerService {
	return &processManagerService{}
}

func (s *processManagerService) FindProcessByName(name string) ([]int, error) {
	cmd := exec.Command("pgrep", "-f", name)
	output, err := cmd.Output()
	if err != nil {
		// pgrep returns error if no processes found
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return []int{}, nil
		}
		return nil, fmt.Errorf("failed to find processes by name: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	pids := make([]int, 0, len(lines))
	for _, line := range lines {
		if line == "" {
			continue
		}
		pid, err := strconv.Atoi(line)
		if err != nil {
			continue
		}
		pids = append(pids, pid)
	}
	return pids, nil
}

func (s *processManagerService) IsProcessRunning(pid int) (bool, error) {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false, err
	}

	// On Unix-like systems, FindProcess always succeeds, so we need to send signal 0
	// to actually check if the process exists
	err = process.Signal(syscall.Signal(0))
	if err != nil {
		if err.Error() == "operation not permitted" {
			// Process exists but may be owned by someone else
			return true, nil
		}
		// Process doesn't exist or other error
		return false, nil
	}
	return true, nil
}

func (s *processManagerService) KillProcess(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return process.Kill()
}

func (s *processManagerService) StartProcess(name string, arg ...string) (*exec.Cmd, error) {
	cmd := exec.Command(name, arg...)
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return cmd, nil
}

func (s *processManagerService) OpenURL(url string) error {
	cmd := exec.Command("open", url)
	return cmd.Run()
}

func (s *processManagerService) OpenFile(filePath string) error {
	cmd := exec.Command("open", filePath)
	return cmd.Run()
}

func (s *processManagerService) OpenDirectory(dirPath string) error {
	cmd := exec.Command("open", dirPath)
	return cmd.Run()
}
