package data

import (
	"encoding/json"
	"go-stock/backend/db"
	"go-stock/backend/logger"
	"gorm.io/gorm"
)

type Settings struct {
	gorm.Model
	TushareToken           string `json:"tushareToken"`
	LocalPushEnable        bool   `json:"localPushEnable"`
	DingPushEnable         bool   `json:"dingPushEnable"`
	DingRobot              string `json:"dingRobot"`
	UpdateBasicInfoOnStart bool   `json:"updateBasicInfoOnStart"`
	RefreshInterval        int64  `json:"refreshInterval"`

	OpenAiEnable      bool    `json:"openAiEnable"`
	OpenAiBaseUrl     string  `json:"openAiBaseUrl"`
	OpenAiApiKey      string  `json:"openAiApiKey"`
	OpenAiModelName   string  `json:"openAiModelName"`
	OpenAiMaxTokens   int     `json:"openAiMaxTokens"`
	OpenAiTemperature float64 `json:"openAiTemperature"`
	OpenAiApiTimeOut  int     `json:"openAiApiTimeOut"`
	Prompt            string  `json:"prompt"`
	CheckUpdate       bool    `json:"checkUpdate"`
	QuestionTemplate  string  `json:"questionTemplate"`
	CrawlTimeOut      int64   `json:"crawlTimeOut"`
	KDays             int64   `json:"kDays"`
	EnableDanmu       bool    `json:"enableDanmu"`
	BrowserPath       string  `json:"browserPath"`
	EnableNews        bool    `json:"enableNews"`
	DarkTheme         bool    `json:"darkTheme"`
	BrowserPoolSize   int     `json:"browserPoolSize"`
	EnableFund        bool    `json:"enableFund"`
}

func (receiver Settings) TableName() string {
	return "settings"
}

type SettingsApi struct {
	Config Settings
}

func NewSettingsApi(settings *Settings) *SettingsApi {
	return &SettingsApi{
		Config: *settings,
	}
}

func (s SettingsApi) UpdateConfig() string {
	count := int64(0)
	db.Dao.Model(s.Config).Count(&count)
	if count > 0 {
		db.Dao.Model(s.Config).Where("id=?", s.Config.ID).Updates(map[string]any{
			"local_push_enable":          s.Config.LocalPushEnable,
			"ding_push_enable":           s.Config.DingPushEnable,
			"ding_robot":                 s.Config.DingRobot,
			"update_basic_info_on_start": s.Config.UpdateBasicInfoOnStart,
			"refresh_interval":           s.Config.RefreshInterval,
			"open_ai_enable":             s.Config.OpenAiEnable,
			"open_ai_base_url":           s.Config.OpenAiBaseUrl,
			"open_ai_api_key":            s.Config.OpenAiApiKey,
			"open_ai_model_name":         s.Config.OpenAiModelName,
			"open_ai_max_tokens":         s.Config.OpenAiMaxTokens,
			"open_ai_temperature":        s.Config.OpenAiTemperature,
			"tushare_token":              s.Config.TushareToken,
			"prompt":                     s.Config.Prompt,
			"check_update":               s.Config.CheckUpdate,
			"open_ai_api_time_out":       s.Config.OpenAiApiTimeOut,
			"question_template":          s.Config.QuestionTemplate,
			"crawl_time_out":             s.Config.CrawlTimeOut,
			"k_days":                     s.Config.KDays,
			"enable_danmu":               s.Config.EnableDanmu,
			"browser_path":               s.Config.BrowserPath,
			"enable_news":                s.Config.EnableNews,
			"dark_theme":                 s.Config.DarkTheme,
			"enable_fund":                s.Config.EnableFund,
		})
	} else {
		logger.SugaredLogger.Infof("未找到配置，创建默认配置:%+v", s.Config)
		db.Dao.Model(s.Config).Create(&Settings{
			LocalPushEnable:        s.Config.LocalPushEnable,
			DingPushEnable:         s.Config.DingPushEnable,
			DingRobot:              s.Config.DingRobot,
			UpdateBasicInfoOnStart: s.Config.UpdateBasicInfoOnStart,
			RefreshInterval:        s.Config.RefreshInterval,
			OpenAiEnable:           s.Config.OpenAiEnable,
			OpenAiBaseUrl:          s.Config.OpenAiBaseUrl,
			OpenAiApiKey:           s.Config.OpenAiApiKey,
			OpenAiModelName:        s.Config.OpenAiModelName,
			OpenAiMaxTokens:        s.Config.OpenAiMaxTokens,
			OpenAiTemperature:      s.Config.OpenAiTemperature,
			TushareToken:           s.Config.TushareToken,
			Prompt:                 s.Config.Prompt,
			CheckUpdate:            s.Config.CheckUpdate,
			OpenAiApiTimeOut:       s.Config.OpenAiApiTimeOut,
			QuestionTemplate:       s.Config.QuestionTemplate,
			CrawlTimeOut:           s.Config.CrawlTimeOut,
			KDays:                  s.Config.KDays,
			EnableDanmu:            s.Config.EnableDanmu,
			BrowserPath:            s.Config.BrowserPath,
			EnableNews:             s.Config.EnableNews,
			DarkTheme:              s.Config.DarkTheme,
			EnableFund:             s.Config.EnableFund,
		})
	}
	return "保存成功！"
}
func (s SettingsApi) GetConfig() *Settings {
	var settings Settings
	db.Dao.Model(&Settings{}).First(&settings)

	if settings.OpenAiEnable {
		if settings.OpenAiApiTimeOut <= 0 {
			settings.OpenAiApiTimeOut = 60 * 5
		}
		if settings.CrawlTimeOut <= 0 {
			settings.CrawlTimeOut = 60
		}
		if settings.KDays < 30 {
			settings.KDays = 120
		}
	}
	if settings.BrowserPath == "" {
		settings.BrowserPath, _ = CheckBrowserOnWindows()
	}
	if settings.BrowserPoolSize <= 0 {
		settings.BrowserPoolSize = 1
	}
	return &settings
}

func (s SettingsApi) Export() string {
	d, _ := json.MarshalIndent(s.GetConfig(), "", "    ")
	return string(d)
}
