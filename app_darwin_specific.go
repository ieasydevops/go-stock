//go:build darwin
// +build darwin

package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	log "go-stock/backend/logger"

	"github.com/wailsapp/wails/v2/pkg/options"
)

// getScreenResolution 获取macOS屏幕分辨率
func getScreenResolution() (int, int, error) {
	cmd := exec.Command("system_profiler", "SPDisplaysDataType")
	output, err := cmd.Output()
	if err != nil {
		log.SugaredLogger.Errorf("获取屏幕分辨率失败: %s", err)
		// 返回默认值
		return 1456, 768, nil
	}

	outputStr := string(output)
	// 查找分辨率信息
	lines := strings.Split(outputStr, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Resolution:") {
			// 提取分辨率部分
			parts := strings.Split(strings.TrimSpace(line), ":")
			if len(parts) < 2 {
				continue
			}

			resolutionStr := strings.TrimSpace(parts[1])
			// 处理 "1920 x 1080" 格式
			dimensions := strings.Split(resolutionStr, " x ")
			if len(dimensions) == 2 {
				width, err1 := strconv.Atoi(strings.TrimSpace(dimensions[0]))
				height, err2 := strconv.Atoi(strings.TrimSpace(dimensions[1]))

				if err1 == nil && err2 == nil {
					return width, height, nil
				}
			}
		}
	}

	// 如果无法解析，返回默认值
	return 1456, 768, nil
}

// OnSecondInstanceLaunch 处理第二个应用实例启动时的行为
func OnSecondInstanceLaunch(secondInstanceData options.SecondInstanceData) {
	// macOS版本使用AppleScript显示通知
	notifyScript := fmt.Sprintf(`osascript -e 'display notification "程序已经在运行了" with title "go-stock"'`)
	cmd := exec.Command("bash", "-c", notifyScript)
	if err := cmd.Run(); err != nil {
		log.SugaredLogger.Errorf("显示通知失败: %s", err)
	}

	time.Sleep(time.Second * 3)
}
