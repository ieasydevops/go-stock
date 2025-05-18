//go:build darwin
// +build darwin

package main

import (
	"go-stock/backend/logger"
	"os"
	"os/user"
	"path/filepath"
)

// macOS特定的初始化函数
func macOSInit() {
	// 确保数据目录存在
	ensureDataDirectory()

	// 记录macOS特定的日志信息
	logger.SugaredLogger.Info("Running on macOS")
}

// 确保数据目录存在
func ensureDataDirectory() {
	// 获取当前用户的主目录
	usr, err := user.Current()
	if err != nil {
		logger.SugaredLogger.Errorf("Failed to get current user: %v", err)
		return
	}

	// 创建应用数据目录
	dataDir := filepath.Join(usr.HomeDir, ".go-stock", "data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		logger.SugaredLogger.Errorf("Failed to create data directory: %v", err)
		return
	}

	// 创建软链接（如果不存在）
	currentDir, err := os.Getwd()
	if err != nil {
		logger.SugaredLogger.Errorf("Failed to get current directory: %v", err)
		return
	}

	localDataDir := filepath.Join(currentDir, "data")
	// 检查本地data目录是否存在
	if _, err := os.Stat(localDataDir); os.IsNotExist(err) {
		// 如果不存在，创建软链接
		if err := os.Symlink(dataDir, localDataDir); err != nil {
			logger.SugaredLogger.Errorf("Failed to create symlink: %v", err)
			return
		}
	}

	logger.SugaredLogger.Infof("Data directory initialized: %s", dataDir)
}
