//go:build darwin
// +build darwin

package main

import (
	"context"
	"embed"
	"fmt"
	"go-stock/backend/data"
	"go-stock/backend/db"
	"go-stock/backend/logger"
	"go-stock/backend/models"
	"os"
	"strings"
	"time"

	"encoding/base64"

	"github.com/PuerkitoBio/goquery"
	"github.com/coocood/freecache"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/go-resty/resty/v2"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var appIcon []byte

var (
	Version       = "1.0.0"
	VersionCommit = "dev"
	icon          []byte
	alipay        []byte
	wxpay         []byte
)

// App struct
type App struct {
	ctx   context.Context
	cache *freecache.Cache
}

// NewApp creates a new App application struct
func NewApp() *App {
	cacheSize := 512 * 1024
	cache := freecache.NewCache(cacheSize)
	return &App{
		cache: cache,
	}
}

func init() {
	var err error
	icon, err = os.ReadFile("frontend/dist/icon.png")
	if err != nil {
		logger.SugaredLogger.Errorf("Failed to read icon: %v", err)
	}
	alipay, err = os.ReadFile("frontend/dist/alipay.png")
	if err != nil {
		logger.SugaredLogger.Errorf("Failed to read alipay: %v", err)
	}
	wxpay, err = os.ReadFile("frontend/dist/wxpay.png")
	if err != nil {
		logger.SugaredLogger.Errorf("Failed to read wxpay: %v", err)
	}
	appIcon, err = os.ReadFile("build/appicon.png")
	if err != nil {
		logger.SugaredLogger.Errorf("Failed to read appIcon: %v", err)
	}
}

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:             "go-stock",
		Width:             1024,
		Height:            768,
		MinWidth:          1024,
		MinHeight:         768,
		MaxWidth:          1920,
		MaxHeight:         1080,
		DisableResize:     false,
		Fullscreen:        false,
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: false,
		BackgroundColour:  &options.RGBA{R: 255, G: 255, B: 255, A: 255},
		Assets:            assets,
		Menu:              nil,
		Logger:            nil,
		LogLevel:          0,
		OnStartup:         app.startup,
		OnDomReady:        app.domReady,
		OnBeforeClose:     app.beforeClose,
		OnShutdown:        app.shutdown,
		WindowStartState:  options.Normal,
		Bind: []interface{}{
			app,
		},
		// Windows platform specific options
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
			// DisableFramelessWindowDecorations: false,
			WebviewUserDataPath: "",
			ZoomFactor:          1.0,
		},
		// Mac platform specific options
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: true,
				HideTitle:                  false,
				HideTitleBar:               false,
				FullSizeContent:            false,
				UseToolbar:                 false,
				HideToolbarSeparator:       true,
			},
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			About: &mac.AboutInfo{
				Title:   "go-stock",
				Message: "",
				Icon:    appIcon,
			},
		},
	})

	if err != nil {
		logger.SugaredLogger.Fatal(err)
	}
}

// startup is called at application startup
func (a *App) startup(ctx context.Context) {
	logger.SugaredLogger.Infof("Version:%s", Version)
	// Perform your setup here
	a.ctx = ctx
}

func checkUpdate(a *App) {
	releaseVersion := &models.GitHubReleaseVersion{}
	_, err := resty.New().R().
		SetResult(releaseVersion).
		Get("https://api.github.com/repos/ArvinLovegood/go-stock/releases/latest")
	if err != nil {
		logger.SugaredLogger.Errorf("get github release version error:%s", err.Error())
		return
	}
	logger.SugaredLogger.Infof("releaseVersion:%+v", releaseVersion.TagName)
	if releaseVersion.TagName != Version {
		go runtime.EventsEmit(a.ctx, "updateVersion", releaseVersion)
	}
}

// domReady is called after front-end resources have been loaded
func (a *App) domReady(ctx context.Context) {
	// Add your action here
	//定时更新数据
	go func() {
		config := data.NewSettingsApi(&data.Settings{}).GetConfig()
		interval := config.RefreshInterval
		if interval <= 0 {
			interval = 1
		}
		ticker := time.NewTicker(time.Second * time.Duration(interval))
		defer ticker.Stop()
		for range ticker.C {
			if isTradingTime(time.Now()) {
				MonitorStockPrices(a)
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(60))
		defer ticker.Stop()
		for range ticker.C {
			telegraph := refreshTelegraphList()
			if telegraph != nil {
				go runtime.EventsEmit(a.ctx, "telegraph", telegraph)
			}
		}
	}()
	go runtime.EventsEmit(a.ctx, "telegraph", refreshTelegraphList())
	go MonitorStockPrices(a)

	//检查新版本
	go func() {
		checkUpdate(a)
	}()
}

func refreshTelegraphList() *[]string {
	url := "https://www.cls.cn/telegraph"
	response, err := resty.New().R().
		SetHeader("Referer", "https://www.cls.cn/").
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36 Edg/117.0.2045.60").
		Get(fmt.Sprintf(url))
	if err != nil {
		return &[]string{}
	}
	document, err := goquery.NewDocumentFromReader(strings.NewReader(string(response.Body())))
	if err != nil {
		return &[]string{}
	}
	var telegraph []string
	document.Find("div.telegraph-content-box").Each(func(i int, selection *goquery.Selection) {
		telegraph = append(telegraph, selection.Text())
	})
	return &telegraph
}

// isTradingDay 判断是否是交易日
func isTradingDay(date time.Time) bool {
	weekday := date.Weekday()
	// 判断是否是周末
	if weekday == time.Saturday || weekday == time.Sunday {
		return false
	}
	// 这里可以添加具体的节假日判断逻辑
	// 例如：判断是否是春节、国庆节等
	return true
}

// isTradingTime 判断是否是交易时间
func isTradingTime(date time.Time) bool {
	if !isTradingDay(date) {
		return false
	}

	hour, minute, _ := date.Clock()

	// 判断是否在9:15到11:30之间
	if (hour == 9 && minute >= 15) || (hour == 10) || (hour == 11 && minute <= 30) {
		return true
	}

	// 判断是否在13:00到15:00之间
	if (hour == 13) || (hour == 14) || (hour == 15 && minute <= 0) {
		return true
	}

	return false
}

func MonitorStockPrices(a *App) {
	dest := &[]data.FollowedStock{}
	db.Dao.Model(&data.FollowedStock{}).Find(dest)
	total := float64(0)
	stockInfos := GetStockInfos(*dest...)
	if stockInfos == nil {
		return
	}
	for _, stockInfo := range *stockInfos {
		total += stockInfo.ProfitAmountToday
		price, _ := convertor.ToFloat(stockInfo.Price)
		if stockInfo.PrePrice != price {
			go runtime.EventsEmit(a.ctx, "stock_price", stockInfo)
		}
	}
	if total != 0 {
		go runtime.EventsEmit(a.ctx, "realtime_profit", fmt.Sprintf("  %.2f", total))
	}
}

func GetStockInfos(follows ...data.FollowedStock) *[]data.StockInfo {
	stockCodes := make([]string, 0)
	for _, follow := range follows {
		stockCodes = append(stockCodes, follow.StockCode)
	}
	stockData, err := data.NewStockDataApi().GetStockCodeRealTimeData(stockCodes...)
	if err != nil {
		logger.SugaredLogger.Errorf("get stock code real time data error:%s", err.Error())
		return nil
	}
	stockInfos := make([]data.StockInfo, 0)
	for _, info := range *stockData {
		v, ok := slice.FindBy(follows, func(idx int, follow data.FollowedStock) bool {
			return follow.StockCode == info.Code
		})
		if ok {
			addStockFollowData(v, &info)
			stockInfos = append(stockInfos, info)
		}
	}
	return &stockInfos
}

func getStockInfo(follow data.FollowedStock) *data.StockInfo {
	stockCode := follow.StockCode
	stockDatas, err := data.NewStockDataApi().GetStockCodeRealTimeData(stockCode)
	if err != nil || len(*stockDatas) == 0 {
		return &data.StockInfo{}
	}
	stockData := (*stockDatas)[0]
	addStockFollowData(follow, &stockData)
	return &stockData
}

func addStockFollowData(follow data.FollowedStock, stockData *data.StockInfo) {
	stockData.PrePrice = follow.Price //上次当前价格
	stockData.Sort = follow.Sort
	stockData.CostPrice = follow.CostPrice //成本价
	stockData.CostVolume = follow.Volume   //成本量
	stockData.AlarmChangePercent = follow.AlarmChangePercent
	stockData.AlarmPrice = follow.AlarmPrice
}

func (a *App) GetStockList(keyword string) []data.StockBasic {
	return data.NewStockDataApi().GetStockList(keyword)
}

func (a *App) GetFollowList() []data.FollowedStock {
	dest := &[]data.FollowedStock{}
	db.Dao.Model(&data.FollowedStock{}).Find(dest)
	return *dest
}

func (a *App) Follow(code string) string {
	dest := &data.FollowedStock{}
	db.Dao.Model(&data.FollowedStock{}).Where("stock_code = ?", code).First(dest)
	if dest.StockCode != "" {
		return "已经关注了"
	}
	stockInfo := getStockInfo(data.FollowedStock{StockCode: code})
	if stockInfo.Code == "" {
		return "股票代码不存在"
	}
	price, _ := convertor.ToFloat(stockInfo.Price)
	dest = &data.FollowedStock{
		StockCode: code,
		Price:     price,
		Sort:      999,
	}
	db.Dao.Create(dest)
	return "关注成功"
}

func (a *App) UnFollow(code string) string {
	db.Dao.Where("stock_code = ?", code).Delete(&data.FollowedStock{})
	return "取消关注成功"
}

func (a *App) GetConfig() data.Settings {
	return *data.NewSettingsApi(&data.Settings{}).GetConfig()
}

func (a *App) UpdateConfig(config data.Settings) string {
	data.NewSettingsApi(&config).UpdateConfig()
	return "更新成功"
}

func (a *App) GetVersionInfo() models.VersionInfo {
	return models.VersionInfo{
		Version: Version,
		Icon:    base64.StdEncoding.EncodeToString(icon),
		Alipay:  base64.StdEncoding.EncodeToString(alipay),
		Wxpay:   base64.StdEncoding.EncodeToString(wxpay),
	}
}

func (a *App) CheckUpdate() {
	checkUpdate(a)
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	dialog, err := runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
		Type:         runtime.QuestionDialog,
		Title:        "go-stock",
		Message:      "确定关闭吗？",
		Buttons:      []string{"确定"},
		Icon:         appIcon,
		CancelButton: "取消",
	})

	if err != nil {
		logger.SugaredLogger.Errorf("dialog error:%s", err.Error())
		return false
	}
	logger.SugaredLogger.Debugf("dialog:%s", dialog)
	if dialog == "No" {
		return true
	}
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
}
