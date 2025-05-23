package models

import "gorm.io/gorm"

// AIResponseResult represents an AI response result
type AIResponseResult struct {
	gorm.Model
	StockCode string `json:"stockCode"`
	StockName string `json:"stockName"`
	Result    string `json:"result"`
	ChatId    string `json:"chatId"`
	Question  string `json:"question"`
	ModelName string `json:"modelName"`
	Content   string `json:"content"`
}
type GitHubReleaseVersion struct {
	TagName string `json:"tag_name"`
	Tag     Tag    `json:"-"`
	Commit  Commit `json:"-"`
}
type Tag struct {
	Object TagObject `json:"object"`
}
type TagObject struct {
	Url string `json:"url"`
}
type Commit struct {
	TreeInfo Tree `json:"tree"`
}
type Tree struct{}
type StockInfoHK struct {
	gorm.Model
	Code string `json:"code"`
	Name string `json:"name"`
}
type StockInfoUS struct {
	gorm.Model
	Code string `json:"code"`
	Name string `json:"name"`
}
type Telegraph struct {
	gorm.Model
	Title         string          `json:"title"`
	Source        string          `json:"source"`
	Content       string          `json:"content"`
	Time          string          `json:"time"`
	Url           string          `json:"url"`
	IsRed         bool            `json:"isRed"`
	SubjectTags   []string        `json:"subjectTags" gorm:"-"`
	StocksTags    []string        `json:"stocksTags" gorm:"-"`
	TelegraphTags []TelegraphTags `json:"telegraphTags" gorm:"-"`
}
type TelegraphTags struct {
	gorm.Model
	TelegraphId uint `json:"telegraphId"`
	TagId       uint `json:"tagId"`
}
type PromptTemplate struct {
	gorm.Model
	Name    string `json:"name"`
	Content string `json:"content"`
	Type    string `json:"type"`
	ID      uint   `json:"id"`
}
type Prompt struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Type    string `json:"type"`
}
type Tags struct {
	gorm.Model
	Name string `json:"name"`
	Type string `json:"type"`
}
type VersionInfo struct {
	Version       string `json:"version"`
	VersionCommit string `json:"versionCommit"`
	GoVersion     string `json:"goVersion"`
	Icon          string `json:"icon"`
	Alipay        string `json:"alipay"`
	Wxpay         string `json:"wxpay"`
	Content       string `json:"content"`
}

// SinaStockInfo represents Sina stock information
type SinaStockInfo struct {
	gorm.Model
	Code       string `json:"code"`
	Name       string `json:"name"`
	Symbol     string `json:"symbol"`
	Open       string `json:"open"`
	Close      string `json:"close"`
	Now        string `json:"now"`
	High       string `json:"high"`
	Low        string `json:"low"`
	Buy        string `json:"buy"`
	Sell       string `json:"sell"`
	Turnover   string `json:"turnover"`
	Volume     string `json:"volume"`
	Bid1Volume string `json:"bid1Volume"`
	Bid1       string `json:"bid1"`
	Bid2Volume string `json:"bid2Volume"`
	Bid2       string `json:"bid2"`
	Bid3Volume string `json:"bid3Volume"`
	Bid3       string `json:"bid3"`
	Bid4Volume string `json:"bid4Volume"`
	Bid4       string `json:"bid4"`
	Bid5Volume string `json:"bid5Volume"`
	Bid5       string `json:"bid5"`
	Ask1Volume string `json:"ask1Volume"`
	Ask1       string `json:"ask1"`
	Ask2Volume string `json:"ask2Volume"`
	Ask2       string `json:"ask2"`
	Ask3Volume string `json:"ask3Volume"`
	Ask3       string `json:"ask3"`
	Ask4Volume string `json:"ask4Volume"`
	Ask4       string `json:"ask4"`
	Ask5Volume string `json:"ask5Volume"`
	Ask5       string `json:"ask5"`
	Date       string `json:"date"`
	Time       string `json:"time"`
}

// Resp is used for OpenAI API responses
type Resp struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}
