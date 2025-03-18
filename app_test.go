package main

import (
	"go-stock/backend/logger"
	"testing"
	"time"
)

// @Author spark
// @Date 2025/2/24 9:35
// @Desc
// -----------------------------------------------------------------------------------

// IsHKTradingTime 判断是否在港股交易时间内
func IsHKTradingTime(date time.Time) bool {
	hour, minute, _ := date.Clock()
	// 判断是否在9:30到16:00之间
	if (hour == 9 && minute >= 30) || (hour >= 10 && hour < 16) || (hour == 16 && minute <= 0) {
		return true
	}
	return false
}

// IsUSTradingTime 判断是否在美股交易时间内
func IsUSTradingTime(date time.Time) bool {
	hour, minute, _ := date.Clock()
	// 判断是否在21:30到次日4:00之间
	if (hour >= 21 && minute >= 30) || (hour >= 22 && hour < 24) || (hour >= 0 && hour < 4) || (hour == 4 && minute <= 0) {
		return true
	}
	return false
}

func TestIsHKTradingTime(t *testing.T) {
	f := IsHKTradingTime(time.Now())
	t.Log(f)
}

func TestIsUSTradingTime(t *testing.T) {
	date := time.Now()
	hour, minute, _ := date.Clock()
	logger.SugaredLogger.Infof("当前时间: %d:%d", hour, minute)

	t.Log(IsUSTradingTime(time.Now()))
}
