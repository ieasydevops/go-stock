package models

import (
	"encoding/json"
	"go-stock/backend/db"
	"go-stock/backend/logger"
	"os"
	"testing"
	"time"

	"github.com/duke-git/lancet/v2/strutil"
)

// @Author spark
// @Date 2025/2/22 16:09
// @Desc
// -----------------------------------------------------------------------------------
type StockInfoHKResp struct {
	Code       int              `json:"code"`
	Status     string           `json:"status"`
	StockInfos *[]StockInfoData `json:"data"`
}

type StockInfoData struct {
	C string `json:"c"`
	N string `json:"n"`
	T string `json:"t"`
	E string `json:"e"`
}

func TestStockInfoHK(t *testing.T) {
	db.Init("../../data/stock.db")
	db.Dao.AutoMigrate(&StockInfoHK{})
	bs, _ := os.ReadFile("../../build/hk.json")
	v := &StockInfoHKResp{}
	err := json.Unmarshal(bs, v)
	if err != nil {
		return
	}
	hks := &[]StockInfoHK{}
	for i, data := range *v.StockInfos {
		logger.SugaredLogger.Infof("第%d条数据: %+v", i, data)
		hk := &StockInfoHK{
			Code:  strutil.PadStart(data.C, 5, "0") + ".HK",
			EName: data.N,
		}
		*hks = append(*hks, *hk)
	}
	db.Dao.Create(&hks)

}

func TestIsHKTradingTime(t *testing.T) {
	// 测试非交易时间
	nonTradingTimes := []time.Time{
		time.Date(2024, 3, 18, 9, 0, 0, 0, time.Local),   // 9:00
		time.Date(2024, 3, 18, 9, 29, 0, 0, time.Local),  // 9:29
		time.Date(2024, 3, 18, 16, 0, 0, 0, time.Local),  // 16:00
		time.Date(2024, 3, 18, 16, 30, 0, 0, time.Local), // 16:30
	}

	// 测试交易时间
	tradingTimes := []time.Time{
		time.Date(2024, 3, 18, 9, 30, 0, 0, time.Local),  // 9:30
		time.Date(2024, 3, 18, 10, 0, 0, 0, time.Local),  // 10:00
		time.Date(2024, 3, 18, 15, 30, 0, 0, time.Local), // 15:30
		time.Date(2024, 3, 18, 15, 59, 0, 0, time.Local), // 15:59
	}

	for _, tt := range nonTradingTimes {
		if IsHKTradingTime() {
			t.Errorf("IsHKTradingTime() 在非交易时间 %v 返回了 true", tt)
		}
	}

	for _, tt := range tradingTimes {
		if !IsHKTradingTime() {
			t.Errorf("IsHKTradingTime() 在交易时间 %v 返回了 false", tt)
		}
	}
}

func TestIsUSTradingTime(t *testing.T) {
	// 测试非交易时间
	nonTradingTimes := []time.Time{
		time.Date(2024, 3, 18, 4, 0, 0, 0, time.Local),   // 4:00
		time.Date(2024, 3, 18, 8, 0, 0, 0, time.Local),   // 8:00
		time.Date(2024, 3, 18, 21, 0, 0, 0, time.Local),  // 21:00
		time.Date(2024, 3, 18, 21, 29, 0, 0, time.Local), // 21:29
	}

	// 测试交易时间
	tradingTimes := []time.Time{
		time.Date(2024, 3, 18, 0, 0, 0, 0, time.Local),   // 0:00
		time.Date(2024, 3, 18, 3, 59, 0, 0, time.Local),  // 3:59
		time.Date(2024, 3, 18, 21, 30, 0, 0, time.Local), // 21:30
		time.Date(2024, 3, 18, 23, 59, 0, 0, time.Local), // 23:59
	}

	for _, tt := range nonTradingTimes {
		if IsUSTradingTime() {
			t.Errorf("IsUSTradingTime() 在非交易时间 %v 返回了 true", tt)
		}
	}

	for _, tt := range tradingTimes {
		if !IsUSTradingTime() {
			t.Errorf("IsUSTradingTime() 在交易时间 %v 返回了 false", tt)
		}
	}
}
