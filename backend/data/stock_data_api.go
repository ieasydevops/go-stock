package data

// @Author spark
// @Date 2024/12/10 9:21
// @Desc
//-----------------------------------------------------------------------------------
import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-stock/backend/db"
	"go-stock/backend/logger"
	"go-stock/backend/models"
	"io"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/duke-git/lancet/v2/strutil"
	"github.com/duke-git/lancet/v2/validator"
	"github.com/go-resty/resty/v2"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"gorm.io/gorm"
)

const sinaStockUrl = "http://hq.sinajs.cn/rn=%d&list=%s"
const tushareApiUrl = "http://api.tushare.pro"

type StockDataApi struct {
	client *resty.Client
	config *Settings
}
type StockInfo struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt models.DateTime
	UpdatedAt models.DateTime
	DeletedAt bool    `gorm:"softDelete:flag"`
	Date      string  `json:"日期" gorm:"index"`
	Time      string  `json:"时间" gorm:"index"`
	Code      string  `json:"股票代码" gorm:"index"`
	Name      string  `json:"股票名称" gorm:"index"`
	PrePrice  float64 `json:"上次当前价格"`
	Price     string  `json:"当前价格"`
	Volume    string  `json:"成交的股票数"`
	Amount    string  `json:"成交金额"`
	Open      string  `json:"今日开盘价"`
	PreClose  string  `json:"昨日收盘价"`
	High      string  `json:"今日最高价"`
	Low       string  `json:"今日最低价"`
	Bid       string  `json:"竞买价"`
	Ask       string  `json:"竞卖价"`
	B1P       string  `json:"买一报价"`
	B1V       string  `json:"买一申报"`
	B2P       string  `json:"买二报价"`
	B2V       string  `json:"买二申报"`
	B3P       string  `json:"买三报价"`
	B3V       string  `json:"买三申报"`
	B4P       string  `json:"买四报价"`
	B4V       string  `json:"买四申报"`
	B5P       string  `json:"买五报价"`
	B5V       string  `json:"买五申报"`
	A1P       string  `json:"卖一报价"`
	A1V       string  `json:"卖一申报"`
	A2P       string  `json:"卖二报价"`
	A2V       string  `json:"卖二申报"`
	A3P       string  `json:"卖三报价"`
	A3V       string  `json:"卖三申报"`
	A4P       string  `json:"卖四报价"`
	A4V       string  `json:"卖四申报"`
	A5P       string  `json:"卖五报价"`
	A5V       string  `json:"卖五申报"`
	Market    string  `json:"市场"`
	BA        string  `json:"盘前盘后"`
	BAChange  string  `json:"盘前盘后涨跌幅"`

	//以下是字段值需二次计算
	ChangePercent     float64 `json:"changePercent"`     //涨跌幅
	ChangePrice       float64 `json:"changePrice"`       //涨跌额
	HighRate          float64 `json:"highRate"`          //最高涨跌
	LowRate           float64 `json:"lowRate"`           //最低涨跌
	CostPrice         float64 `json:"costPrice"`         //成本价
	CostVolume        int64   `json:"costVolume"`        //持仓数量
	Profit            float64 `json:"profit"`            //总盈亏率
	ProfitAmount      float64 `json:"profitAmount"`      //总盈亏金额
	ProfitAmountToday float64 `json:"profitAmountToday"` //今日盈亏金额

	Sort               int64   `json:"sort"` //排序
	AlarmChangePercent float64 `json:"alarmChangePercent"`
	AlarmPrice         float64 `json:"alarmPrice"`
}

func (receiver StockInfo) TableName() string {
	return "stock_info"
}

type TushareRequest struct {
	ApiName string `json:"api_name"`
	Token   string `json:"token"`
	Params  any    `json:"params"`
	Fields  string `json:"fields"`
}
type TushareResponse struct {
	RequestId string `json:"request_id"`
	Code      int    `json:"code"`
	Data      any    `json:"data"`
	Msg       string `json:"msg"`
}

/*
	字段	类型	说明
	ts_code	str	TS代码
	symbol	str	股票代码
	name	str	股票名称
	area	str	地域
	industry	str	所属行业
	fullname	str	股票全称
	enname	str	英文全称
	cnspell	str	拼音缩写
	market	str	市场类型
	exchange	str	交易所代码
	curr_type	str	交易货币
	list_status	str	上市状态 L上市 D退市 P暂停上市
	list_date	str	上市日期
	delist_date	str	退市日期
	is_hs	str	是否沪深港通标的，N否 H沪股通 S深股通
	act_name	str	实控人名称
	act_ent_type	str	实控人企业性质*/

type StockBasic struct {
	ID         uint `gorm:"primarykey"`
	CreatedAt  models.DateTime
	UpdatedAt  models.DateTime
	DeletedAt  bool   `gorm:"softDelete:flag"`
	TsCode     string `json:"ts_code" gorm:"index"`
	Symbol     string `json:"symbol" gorm:"index"`
	Name       string `json:"name" gorm:"index"`
	Area       string `json:"area"`
	Industry   string `json:"industry" gorm:"index"`
	Fullname   string `json:"fullname"`
	Ename      string `json:"enname"`
	Cnspell    string `json:"cnspell"`
	Market     string `json:"market"`
	Exchange   string `json:"exchange"`
	CurrType   string `json:"curr_type"`
	ListStatus string `json:"list_status"`
	ListDate   string `json:"list_date"`
	DelistDate string `json:"delist_date"`
	IsHs       string `json:"is_hs"`
	ActName    string `json:"act_name"`
	ActEntType string `json:"act_ent_type"`
}

type FollowedStock struct {
	StockCode          string
	Name               string
	Volume             int64
	CostPrice          float64
	Price              float64
	PriceChange        float64
	ChangePercent      float64
	AlarmChangePercent float64
	AlarmPrice         float64
	Time               models.DateTime
	Sort               int64
	IsDel              bool `gorm:"softDelete:flag"`
}

func (receiver FollowedStock) TableName() string {
	return "followed_stock"
}

type TushareStockBasicResponse struct {
	TushareResponse
	Data StockBasicResponse `json:"data"`
}

type StockBasicResponse struct {
	Fields  []string `json:"fields"`
	Items   [][]any  `json:"items"`
	HasMore bool     `json:"has_more"`
	Count   int      `json:"count"`
}

func (receiver StockBasic) TableName() string {
	return "tushare_stock_basic"
}
func NewStockDataApi() *StockDataApi {
	return &StockDataApi{
		client: resty.New(),
		config: getConfig(),
	}
}

// GetIndexBasic 获取指数信息
func (receiver StockDataApi) GetIndexBasic() {
	res := &TushareStockBasicResponse{}
	fields := "ts_code,name,market,publisher,category,base_date,base_point,list_date,fullname,index_type,weight_rule,desc"
	_, err := receiver.client.R().
		SetHeader("content-type", "application/json").
		SetBody(&TushareRequest{
			ApiName: "index_basic",
			Token:   receiver.config.TushareToken,
			Params:  nil,
			Fields:  fields}).
		SetResult(res).
		Post(tushareApiUrl)
	if err != nil {
		logger.SugaredLogger.Error(err.Error())
		return
	}
	if res.Code != 0 {
		logger.SugaredLogger.Error(res.Msg)
		return
	}
	//ioutil.WriteFile("index_basic.json", resp.Body(), 0666)

	for _, item := range res.Data.Items {
		data := map[string]any{}
		for _, field := range strings.Split(fields, ",") {
			idx := slice.IndexOf(res.Data.Fields, field)
			if idx == -1 {
				continue
			}
			data[field] = item[idx]
		}
		index := &IndexBasic{}
		jsonData, _ := json.Marshal(data)
		err := json.Unmarshal(jsonData, index)
		if err != nil {
			continue
		}
		index.ID = 0
		db.Dao.Model(&IndexBasic{}).FirstOrCreate(index, &IndexBasic{TsCode: index.TsCode}).Where("ts_code = ?", index.TsCode).Updates(index)
	}

}

// map转换为结构体

func (receiver StockDataApi) GetStockBaseInfo() {
	res := &TushareStockBasicResponse{}
	fields := "ts_code,symbol,name,area,industry,cnspell,market,list_date,act_name,act_ent_type,fullname,exchange,list_status,curr_type,enname,delist_date,is_hs"
	_, err := receiver.client.R().
		SetHeader("content-type", "application/json").
		SetBody(&TushareRequest{
			ApiName: "stock_basic",
			Token:   receiver.config.TushareToken,
			Params:  nil,
			Fields:  fields,
		}).
		SetResult(res).
		Post(tushareApiUrl)
	//logger.SugaredLogger.Infof("GetStockBaseInfo %s", string(resp.Body()))
	//resp.Body()写入文件
	//ioutil.WriteFile("stock_basic.json", resp.Body(), 0666)
	//logger.SugaredLogger.Infof("GetStockBaseInfo %+v", res)
	if err != nil {
		logger.SugaredLogger.Error(err.Error())
		return
	}
	if res.Code != 0 {
		logger.SugaredLogger.Error(res.Msg)
		return
	}
	for _, item := range res.Data.Items {
		stock := &StockBasic{}
		data := map[string]any{}
		for _, field := range strings.Split(fields, ",") {
			logger.SugaredLogger.Infof("field: %s", field)
			idx := slice.IndexOf(res.Data.Fields, field)
			if idx == -1 {
				continue
			}
			data[field] = item[idx]
		}
		jsonData, _ := json.Marshal(data)
		err := json.Unmarshal(jsonData, stock)
		if err != nil {
			continue
		}
		stock.ID = 0
		db.Dao.Model(&StockBasic{}).FirstOrCreate(stock, &StockBasic{TsCode: stock.TsCode}).Where("ts_code = ?", stock.TsCode).Updates(stock)
	}

}

func (receiver StockDataApi) GetStockCodeRealTimeData(StockCodes ...string) (*[]StockInfo, error) {

	codes := slice.JoinFunc(StockCodes, ",", func(s string) string {
		if strings.HasPrefix(s, "us") {
			s = strings.Replace(s, "us", "gb_", 1)
		}
		if strings.HasPrefix(s, "US") {
			s = strings.Replace(s, "US", "gb_", 1)
		}
		return strings.ToLower(s)
	})

	url := fmt.Sprintf(sinaStockUrl, time.Now().Unix(), codes)
	//logger.SugaredLogger.Infof("GetStockCodeRealTimeData %s", url)
	resp, err := receiver.client.R().
		SetHeader("Host", "hq.sinajs.cn").
		SetHeader("Referer", "https://finance.sina.com.cn/").
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36 Edg/119.0.0.0").
		Get(url)
	if err != nil {
		logger.SugaredLogger.Error(err.Error())
		return &[]StockInfo{}, err
	}

	stockInfos := make([]StockInfo, 0)
	str := GB18030ToUTF8(resp.Body())
	dataStr := strutil.SplitEx(str, "\n", true)
	if len(dataStr) == 0 {
		return &[]StockInfo{}, errors.New("获取股票信息失败,请检查股票代码是否正确")
	}
	for _, data := range dataStr {
		//logger.SugaredLogger.Info(data)
		stockData, err := ParseFullSingleStockData(data)
		logger.SugaredLogger.Infof("GetStockCodeRealTimeData %v", stockData)
		if err != nil {
			logger.SugaredLogger.Error(err.Error())
			continue
		}
		stockInfos = append(stockInfos, *stockData)

		go func() {
			var count int64
			db.Dao.Model(&StockInfo{}).Where("code = ?", stockData.Code).Count(&count)
			if count == 0 {
				db.Dao.Model(&StockInfo{}).Create(stockData)
			} else {
				db.Dao.Model(&StockInfo{}).Where("code = ?", stockData.Code).Updates(stockData)
			}
		}()

	}

	return &stockInfos, err
}

func (receiver StockDataApi) Follow(stockCode string, name string, volume int64, costPrice float64) *FollowedStock {
	price := 0.0
	if stockCode == "" {
		return nil
	}
	db.Dao.Model(&FollowedStock{}).FirstOrCreate(&FollowedStock{
		StockCode:          stockCode,
		Name:               name,
		Volume:             volume,
		CostPrice:          costPrice,
		Price:              price,
		PriceChange:        0,
		Sort:               0,
		AlarmChangePercent: 3,
		AlarmPrice:         price + 1,
		Time:               models.DateTime{Time: time.Now()},
	}, &FollowedStock{StockCode: stockCode})
	return nil
}

func (receiver StockDataApi) UnFollow(stockCode string) string {
	if strutil.HasPrefixAny(stockCode, []string{"gb_"}) {
		stockCode = strings.ToUpper(stockCode)
		stockCode = strings.Replace(stockCode, "gb_", "us", 1)
		stockCode = strings.Replace(stockCode, "GB_", "us", 1)
	}
	db.Dao.Model(&FollowedStock{}).Where("stock_code = ?", strings.ToLower(stockCode)).Delete(&FollowedStock{})
	return "取消关注成功"
}

func (receiver StockDataApi) SetCostPriceAndVolume(price float64, volume int64, stockCode string) string {
	if strutil.HasPrefixAny(stockCode, []string{"gb_"}) {
		stockCode = strings.ToUpper(stockCode)
		stockCode = strings.Replace(stockCode, "gb_", "us", 1)
		stockCode = strings.Replace(stockCode, "GB_", "us", 1)
	}
	err := db.Dao.Model(&FollowedStock{}).Where("stock_code = ?", strings.ToLower(stockCode)).Update("cost_price", price).Update("volume", volume).Error
	if err != nil {
		logger.SugaredLogger.Error(err.Error())
		return "设置失败"
	}
	return "设置成功"
}

func (receiver StockDataApi) SetAlarmChangePercent(val, alarmPrice float64, stockCode string) string {
	if strutil.HasPrefixAny(stockCode, []string{"gb_"}) {
		stockCode = strings.ToUpper(stockCode)
		stockCode = strings.Replace(stockCode, "gb_", "us", 1)
		stockCode = strings.Replace(stockCode, "GB_", "us", 1)
	}
	err := db.Dao.Model(&FollowedStock{}).Where("stock_code = ?", strings.ToLower(stockCode)).Updates(&map[string]any{
		"alarm_change_percent": val,
		"alarm_price":          alarmPrice,
	}).Error
	if err != nil {
		logger.SugaredLogger.Error(err.Error())
		return "设置失败"
	}
	return "设置成功"
}

func (receiver StockDataApi) SetStockSort(sort int64, stockCode string) {
	if strutil.HasPrefixAny(stockCode, []string{"gb_"}) {
		stockCode = strings.ToUpper(stockCode)
		stockCode = strings.Replace(stockCode, "gb_", "us", 1)
		stockCode = strings.Replace(stockCode, "GB_", "us", 1)
	}
	db.Dao.Model(&FollowedStock{}).Where("stock_code = ?", strings.ToLower(stockCode)).Update("sort", sort)
}

func (receiver StockDataApi) GetFollowList() *[]FollowedStock {
	var result *[]FollowedStock
	db.Dao.Model(&FollowedStock{}).Order("sort asc,time desc").Find(&result)
	return result
}

func (receiver StockDataApi) GetStockList(key string) []StockBasic {
	var result []StockBasic
	db.Dao.Model(&StockBasic{}).Where("name like ? or ts_code like ?", "%"+key+"%", "%"+key+"%").Find(&result)
	var result2 []IndexBasic
	db.Dao.Model(&IndexBasic{}).Where("market in ?", []string{"SSE", "SZSE"}).Where("name like ? or ts_code like ?", "%"+key+"%", "%"+key+"%").Find(&result2)

	var result3 []models.StockInfoHK
	db.Dao.Model(&models.StockInfoHK{}).Where("name like ? or code like ?", "%"+key+"%", "%"+key+"%").Find(&result3)

	var result4 []models.StockInfoUS
	db.Dao.Model(&models.StockInfoUS{}).Where("name like ? or code like ? or e_name like ?", "%"+key+"%", "%"+key+"%", "%"+key+"%").Find(&result4)

	for _, item := range result2 {
		result = append(result, StockBasic{
			TsCode:   item.TsCode,
			Name:     item.Name,
			Fullname: item.FullName,
			Symbol:   item.Symbol,
			Market:   item.Market,
			ListDate: item.ListDate,
		})
	}
	for _, item := range result3 {
		result = append(result, StockBasic{
			TsCode:   item.Code,
			Name:     item.Name,
			Fullname: item.Name,
			Market:   "HK",
		})
	}
	for _, item := range result4 {
		result = append(result, StockBasic{
			TsCode:   strings.ToLower(strings.Replace(item.Code, "us", "gb_", 1)),
			Name:     item.Name,
			Fullname: item.Name,
			Market:   "US",
		})
	}

	return result
}

// GB18030ToUTF8 GB18030 转换为 UTF8
func GB18030ToUTF8(bs []byte) string {
	reader := transform.NewReader(bytes.NewReader(bs), simplifiedchinese.GB18030.NewDecoder())
	d, err := io.ReadAll(reader)
	if err != nil {
		panic(err)
	}
	return string(d)
}

func ParseFullSingleStockData(data string) (*StockInfo, error) {
	datas := strutil.SplitAndTrim(data, "=", "\"")
	if len(datas) < 2 {
		return nil, fmt.Errorf("invalid data format")
	}
	var result map[string]string
	var err error
	if strutil.ContainsAny(datas[0], []string{"hq_str_sz", "hq_str_sh", "hq_str_bj", "hq_str_sb"}) {
		result, err = ParseSHSZStockData(datas)
	}
	if strutil.ContainsAny(datas[0], []string{"hq_str_hk"}) {
		result, err = ParseHKStockData(datas)
	}
	if strutil.ContainsAny(datas[0], []string{"hq_str_gb"}) {
		result, err = ParseUSStockData(datas)
	}

	//logger.SugaredLogger.Infof("股票数据解析完成: %v", result)
	marshal, err := json.Marshal(result)
	if err != nil {
		logger.SugaredLogger.Errorf("json.Marshal error:%s", err.Error())
		return nil, err
	}
	//logger.SugaredLogger.Infof("股票数据解析完成marshal: %s", marshal)
	stockInfo := &StockInfo{}
	err = json.Unmarshal(marshal, &stockInfo)
	if err != nil {
		logger.SugaredLogger.Errorf("json.Unmarshal error:%s", err.Error())
		return nil, err
	}
	//logger.SugaredLogger.Infof("股票数据解析完成stockInfo: %+v", stockInfo)

	return stockInfo, nil
}

func ParseUSStockData(datas []string) (map[string]string, error) {
	code := strings.Split(datas[0], "hq_str_")[1]
	result := make(map[string]string)
	parts := strutil.SplitAndTrim(datas[1], ",", "\"", ";")
	//parts := strings.Split(data, ",")
	logger.SugaredLogger.Infof("股票数据解析完成: parts:%d", len(parts))
	if len(parts) < 35 {
		return nil, fmt.Errorf("invalid data format")
	}
	/*
		谷歌,   0
		170.2100, 1 现价
		-2.57, 2 涨跌幅
		2025-02-28 09:38:50, 3 时间
		-4.4900, 4 涨跌额
		175.9400, 5 今日开盘价
		176.5900, 6 区间
		169.7520, 7 区间
		208.7000, 8 52周区间
		130.9500, 9 52周区间
		25930485, 10 成交量
		17083496, 11 10日均量
		2074859900000, 12 市值
		8.13, 13 每股收益
		20.940000 , 14 市盈率
		0.00,  15
		0.00,  16
		0.20,  17
		0.00,	18
		12190000000, 19
		71, 20
		170.2000, 21 盘前盘后盘
		-0.01, 22  盘前盘后涨跌幅
		-0.01, 23
		Feb 27 07:59PM EST, 24
		Feb 27 04:00PM EST, 25
		174.7000, 26 前收盘
		2917444, 27
		1, 28
		2025, 29
		4456143849.0000, 30
		176.1200, 31
		163.7039, 32
		496605933.1411, 33
		170.2100, 34 现价
		174.7000 35 前收盘
	*/
	result["股票代码"] = code
	result["股票名称"] = parts[0]
	result["今日开盘价"] = parts[5]

	if len(parts) >= 36 {
		result["昨日收盘价"] = strutil.ReplaceWithMap(strutil.RemoveNonPrintable(parts[26]), map[string]string{"\"": "", ";": ""})
	} else {
		result["昨日收盘价"] = strutil.ReplaceWithMap(strutil.RemoveNonPrintable(parts[len(parts)-1]), map[string]string{"\"": "", ";": ""})
	}

	result["今日最高价"] = parts[6]
	result["今日最低价"] = parts[7]
	result["当前价格"] = parts[1]
	result["盘前盘后"] = parts[21]
	result["盘前盘后涨跌幅"] = parts[22]
	result["日期"] = strutil.SplitAndTrim(parts[3], " ", "")[0]
	result["时间"] = strutil.SplitAndTrim(parts[3], " ", "")[1]
	logger.SugaredLogger.Infof("美股股票数据解析完成: %v", result)
	return result, nil
}

func ParseHKStockData(datas []string) (map[string]string, error) {
	code := strings.Split(datas[0], "hq_str_")[1]
	result := make(map[string]string)
	parts := strutil.SplitAndTrim(datas[1], ",", "\"", ";")
	//parts := strings.Split(data, ",")
	if len(parts) < 19 {
		return nil, fmt.Errorf("invalid data format")
	}
	/*
		XIAOMI-W,    0
		小米集团－Ｗ,  1 股票名称
		50.050,		 2 今日开盘价
		49.150,		 3 昨日收盘价
		51.950,      4 今日最高价
		49.700,      5 今日最低价
		51.700,      6 当前价格
		2.550,       7 涨跌额
		5.188,		 8 涨跌幅
		51.65000,    9
		51.70000,    10
		15770408249, 11 成交额
		308362585,   12 成交量
		0.000,       13
		0.000,       14
		51.950,		 15 52周最高
		12.560,		 16 52周最低
		2025/02/21,  17
		16:08        18
	*/
	result["股票代码"] = code
	result["股票名称"] = parts[1]
	result["今日开盘价"] = parts[2]
	result["昨日收盘价"] = parts[3]
	result["今日最高价"] = parts[4]
	result["今日最低价"] = parts[5]
	result["当前价格"] = parts[6]
	result["日期"] = strings.ReplaceAll(parts[17], "/", "-")
	result["时间"] = strings.ReplaceAll(parts[18], "\";", ":00")
	//logger.SugaredLogger.Infof("股票数据解析完成: %v", result)
	return result, nil
}

func ParseSHSZStockData(datas []string) (map[string]string, error) {
	code := strings.Split(datas[0], "hq_str_")[1]
	result := make(map[string]string)
	parts := strutil.SplitAndTrim(datas[1], ",", "\"")
	//parts := strings.Split(data, ",")
	if len(parts) < 32 {
		return nil, fmt.Errorf("invalid data format")
	}
	/*
		0："大秦铁路"，股票名字；
		1："27.55"，今日开盘价；
		2："27.25"，昨日收盘价；
		3："26.91"，当前价格；
		4："27.55"，今日最高价；
		5："26.20"，今日最低价；
		6："26.91"，竞买价，即"买一"报价；
		7："26.92"，竞卖价，即"卖一"报价；
		8："22114263"，成交的股票数，由于股票交易以一百股为基本单位，所以在使用时，通常把该值除以一百；
		9："589824680"，成交金额，单位为"元"，为了便于阅读，通常以"万元"为成交金额的单位，所以通常把该值除以一万；
		10："4695"，"买一"申报4695股，即47手；
		11："26.91"，"买一"报价；
		12："57590"，"买二"
		13："26.90"，"买二"
		14："14700"，"买三"
		15："26.89"，"买三"
		16："14300"，"买四"
		17："26.88"，"买四"
		18："15100"，"买五"
		19："26.87"，"买五"
		20："3100"，"卖一"申报3100股，即31手；
		21："26.92"，"卖一"报价
		(22, 23), (24, 25), (26,27), (28, 29)分别为"卖二"至"卖四的情况"
		30："2008-01-11"，日期；
		31："15:05:32"，时间；*/
	result["股票代码"] = code
	result["股票名称"] = parts[0]
	result["今日开盘价"] = parts[1]
	result["昨日收盘价"] = parts[2]
	result["当前价格"] = parts[3]
	result["今日最高价"] = parts[4]
	result["今日最低价"] = parts[5]
	result["竞买价"] = parts[6]
	result["竞卖价"] = parts[7]
	result["成交的股票数"] = parts[8]
	result["成交金额"] = parts[9]
	result["买一申报"] = parts[10]
	result["买一报价"] = parts[11]
	result["买二申报"] = parts[12]
	result["买二报价"] = parts[13]
	result["买三申报"] = parts[14]
	result["买三报价"] = parts[15]
	result["买四申报"] = parts[16]
	result["买四报价"] = parts[17]
	result["买五申报"] = parts[18]
	result["买五报价"] = parts[19]
	result["卖一申报"] = parts[20]
	result["卖一报价"] = parts[21]
	result["卖二申报"] = parts[22]
	result["卖二报价"] = parts[23]
	result["卖三申报"] = parts[24]
	result["卖三报价"] = parts[25]
	result["卖四申报"] = parts[26]
	result["卖四报价"] = parts[27]
	result["卖五申报"] = parts[28]
	result["卖五报价"] = parts[29]
	result["日期"] = parts[30]
	result["时间"] = parts[31]
	return result, nil
}

type IndexBasic struct {
	gorm.Model
	TsCode        string  `json:"ts_code" gorm:"index"`
	Symbol        string  `json:"symbol" gorm:"index"`
	Name          string  `json:"name" gorm:"index"`
	FullName      string  `json:"fullname"`
	IndexType     string  `json:"index_type"`
	IndexCategory string  `json:"category"`
	Market        string  `json:"market"`
	ListDate      string  `json:"list_date"`
	BaseDate      string  `json:"base_date"`
	BasePoint     float64 `json:"base_point"`
	Publisher     string  `json:"publisher"`
	WeightRule    string  `json:"weight_rule"`
	DESC          string  `json:"desc"`
}

func (IndexBasic) TableName() string {
	return "tushare_index_basic"
}

type RealTimeStockPriceInfo struct {
	StockCode string
	Price     string `json:"当前价格"`
	Time      time.Time
}

func GetRealTimeStockPriceInfo(ctx context.Context, stockCode string) (price, priceTime string) {
	if strutil.HasPrefixAny(stockCode, []string{"SZ", "SH", "sh", "sz"}) {
		crawlerAPI := CrawlerApi{}
		crawlerBaseInfo := CrawlerBaseInfo{
			Name:        "EastmoneyCrawler",
			Description: "EastmoneyCrawler Description",
			BaseUrl:     "https://quote.eastmoney.com/",
			Headers:     map[string]string{"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36 Edg/133.0.0.0"},
		}
		crawlerAPI = crawlerAPI.NewCrawler(ctx, crawlerBaseInfo)
		htmlContent, ok := crawlerAPI.GetHtml(fmt.Sprintf("https://quote.eastmoney.com/%s.html", stockCode), "div.zxj", true)
		if ok {
			price := ""
			priceTime := ""
			document, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
			if err != nil {
				logger.SugaredLogger.Errorf("GetRealTimeStockPriceInfo error: %v", err)
			}
			document.Find("div.zxj").Each(func(i int, selection *goquery.Selection) {
				price = selection.Text()
				logger.SugaredLogger.Infof("股票代码: %s, 当前价格: %s", stockCode, price)
			})

			document.Find("span.quote_title_time").Each(func(i int, selection *goquery.Selection) {
				priceTime = selection.Text()
				logger.SugaredLogger.Infof("股票代码: %s, 当前价格时间: %s", stockCode, priceTime)
			})
			return price, priceTime
		}
	}
	return price, priceTime
}

func SearchStockPriceInfo(stockCode string, crawlTimeOut int64) *[]string {

	if strutil.HasPrefixAny(stockCode, []string{"SZ", "SH", "sh", "sz"}) {
		return getSHSZStockPriceInfo(stockCode, crawlTimeOut)
	}
	if strutil.HasPrefixAny(stockCode, []string{"HK", "hk"}) {
		return getHKStockPriceInfo(stockCode, crawlTimeOut)
	}
	if strutil.HasPrefixAny(stockCode, []string{"US", "us", "gb_"}) {
		return getUSStockPriceInfo(stockCode, crawlTimeOut)
	}
	return &[]string{}
}

func getUSStockPriceInfo(stockCode string, crawlTimeOut int64) *[]string {
	var messages []string
	crawlerAPI := CrawlerApi{}
	crawlerBaseInfo := CrawlerBaseInfo{
		Name:        "SinaCrawler",
		Description: "SinaCrawler Crawler Description",
		BaseUrl:     "https://stock.finance.sina.com.cn",
		Headers:     map[string]string{"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36 Edg/133.0.0.0"},
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(crawlTimeOut)*time.Second)
	defer cancel()
	crawlerAPI = crawlerAPI.NewCrawler(ctx, crawlerBaseInfo)

	url := fmt.Sprintf("https://stock.finance.sina.com.cn/usstock/quotes/%s.html", strings.ReplaceAll(stockCode, "gb_", ""))
	htmlContent, ok := crawlerAPI.GetHtml(url, "div#hqPrice", true)
	if !ok {
		return &[]string{}
	}
	document, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		logger.SugaredLogger.Error(err.Error())
	}
	stockName := ""
	stockPrice := ""
	stockPriceTime := ""
	document.Find("div.hq_title >h1").Each(func(i int, selection *goquery.Selection) {
		stockName = strutil.RemoveNonPrintable(selection.Text())
		logger.SugaredLogger.Infof("股票名称-:%s", stockName)
	})

	document.Find("#hqPrice").Each(func(i int, selection *goquery.Selection) {
		stockPrice = strutil.RemoveNonPrintable(selection.Text())
		logger.SugaredLogger.Infof("现价: %s", stockPrice)
	})

	document.Find("div.hq_time").Each(func(i int, selection *goquery.Selection) {
		stockPriceTime = strutil.RemoveNonPrintable(selection.Text())
		logger.SugaredLogger.Infof("时间: %s", stockPriceTime)
	})

	messages = append(messages, fmt.Sprintf("%s:%s现价%s", stockPriceTime, stockName, stockPrice))
	logger.SugaredLogger.Infof("股票: %s", messages)

	document.Find("div#hqDetails >table tbody tr").Each(func(i int, selection *goquery.Selection) {
		text := strutil.RemoveNonPrintable(selection.Text())
		logger.SugaredLogger.Infof("股票名称-%s: %s", stockName, text)
		messages = append(messages, text)
	})

	return &messages
}

func getHKStockPriceInfo(stockCode string, crawlTimeOut int64) *[]string {
	var messages []string
	crawlerAPI := CrawlerApi{}
	crawlerBaseInfo := CrawlerBaseInfo{
		Name:        "SinaCrawler",
		Description: "SinaCrawler Crawler Description",
		BaseUrl:     "https://stock.finance.sina.com.cn",
		Headers:     map[string]string{"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36 Edg/133.0.0.0"},
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(crawlTimeOut)*time.Second)
	defer cancel()
	crawlerAPI = crawlerAPI.NewCrawler(ctx, crawlerBaseInfo)

	url := fmt.Sprintf("https://stock.finance.sina.com.cn/hkstock/quotes/%s.html", strings.ReplaceAll(stockCode, "hk", ""))
	htmlContent, ok := crawlerAPI.GetHtml(url, ".deta_hqContainer >.deta03 ", true)
	if !ok {
		return &[]string{}
	}
	document, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		logger.SugaredLogger.Error(err.Error())
	}
	stockName := ""
	stockPrice := ""
	stockPriceTime := ""
	document.Find("#stock_cname").Each(func(i int, selection *goquery.Selection) {
		stockName = strutil.RemoveNonPrintable(selection.Text())
		logger.SugaredLogger.Infof("股票名称-:%s", stockName)
	})

	document.Find("#mts_stock_hk_price").Each(func(i int, selection *goquery.Selection) {
		stockPrice = strutil.RemoveNonPrintable(selection.Text())
		logger.SugaredLogger.Infof("现价: %s", stockPrice)
	})

	document.Find("#mts_stock_hk_time").Each(func(i int, selection *goquery.Selection) {
		stockPriceTime = strutil.RemoveNonPrintable(selection.Text())
		logger.SugaredLogger.Infof("时间: %s", stockPriceTime)
	})

	messages = append(messages, fmt.Sprintf("%s:%s现价%s", stockPriceTime, stockName, stockPrice))
	logger.SugaredLogger.Infof("股票: %s", messages)

	document.Find(".deta_hqContainer >.deta03 li").Each(func(i int, selection *goquery.Selection) {
		text := strutil.RemoveNonPrintable(selection.Text())
		logger.SugaredLogger.Infof("股票名称-%s: %s", stockName, text)
		messages = append(messages, text)
	})

	return &messages
}

func getSHSZStockPriceInfo(stockCode string, crawlTimeOut int64) *[]string {
	var messages []string
	url := "https://www.cls.cn/stock?code=" + stockCode
	// 创建一个 chromedp 上下文
	timeoutCtx, timeoutCtxCancel := context.WithTimeout(context.Background(), time.Duration(crawlTimeOut)*time.Second)
	defer timeoutCtxCancel()
	var ctx context.Context
	var cancel context.CancelFunc
	path := getConfig().BrowserPath
	logger.SugaredLogger.Infof("SearchStockPriceInfo BrowserPath:%s", path)
	if path != "" {
		pctx, pcancel := chromedp.NewExecAllocator(
			timeoutCtx,
			chromedp.ExecPath(path),
			chromedp.Flag("headless", true),
		)
		defer pcancel()
		ctx, cancel = chromedp.NewContext(
			pctx,
			chromedp.WithLogf(logger.SugaredLogger.Infof),
			chromedp.WithErrorf(logger.SugaredLogger.Errorf),
		)
	} else {
		ctx, cancel = chromedp.NewContext(
			timeoutCtx,
			chromedp.WithLogf(logger.SugaredLogger.Infof),
			chromedp.WithErrorf(logger.SugaredLogger.Errorf),
		)
	}
	defer cancel()

	var htmlContent string

	var tasks chromedp.Tasks
	tasks = append(tasks, chromedp.Navigate(url))
	tasks = append(tasks, chromedp.WaitVisible("div.quote-change-box", chromedp.ByQuery))
	tasks = append(tasks, chromedp.ActionFunc(func(ctx context.Context) error {
		price, _ := FetchPrice(ctx)
		logger.SugaredLogger.Infof("price:%s", price)
		return nil
	}))
	tasks = append(tasks, chromedp.OuterHTML("html", &htmlContent, chromedp.ByQuery))

	err := chromedp.Run(ctx, tasks)
	if err != nil {
		logger.SugaredLogger.Error(err.Error())
		return &[]string{}
	}
	document, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		logger.SugaredLogger.Error(err.Error())
		return &[]string{}
	}

	document.Find("div.quote-text-border,span.quote-price").Each(func(i int, selection *goquery.Selection) {
		text := strutil.RemoveNonPrintable(selection.Text())
		logger.SugaredLogger.Info(text)
		messages = append(messages, text)

	})
	return &messages
}
func FetchPrice(ctx context.Context) (string, error) {
	var price string
	timeout := time.After(10 * time.Second)   // 设置超时时间为10秒
	ticker := time.NewTicker(1 * time.Second) // 每秒尝试一次
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return "", fmt.Errorf("timeout reached while fetching price")
		case <-ticker.C:
			err := chromedp.Run(ctx, chromedp.Text("span.quote-price", &price, chromedp.BySearch))
			if err != nil {
				logger.SugaredLogger.Errorf("failed to fetch price: %v", err)
				continue
			}
			logger.SugaredLogger.Infof("price:%s", price)
			if price != "" && validator.IsNumberStr(price) {
				return price, nil
			}
		}
	}
}
func SearchStockInfo(stock, msgType string, crawlTimeOut int64) *[]string {
	crawler := CrawlerApi{
		crawlerBaseInfo: CrawlerBaseInfo{

			Name:        "财联社",
			BaseUrl:     "https://www.cls.cn",
			Description: "财联社",
			Headers:     map[string]string{"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/133.0.0.0 Safari/537.36 Edg/133.0.0.0"},
		},
	}
	timeoutCtx, timeoutCtxCancel := context.WithTimeout(context.Background(), time.Duration(crawlTimeOut)*time.Second)
	defer timeoutCtxCancel()
	crawler = crawler.NewCrawler(timeoutCtx, crawler.crawlerBaseInfo)
	url := fmt.Sprintf("https://www.cls.cn/searchPage?keyword=%s&type=%s", RemoveAllBlankChar(stock), msgType)
	logger.SugaredLogger.Infof("SearchStockInfo url:%s", url)
	waitVisible := ".search-telegraph-list,.subject-interest-list"
	htmlContent, ok := crawler.GetHtml(url, waitVisible, true)
	if !ok {
		return &[]string{}
	}
	document, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		logger.SugaredLogger.Error(err.Error())
		return &[]string{}
	}
	var messages []string
	document.Find(waitVisible).Each(func(i int, selection *goquery.Selection) {
		text := strutil.RemoveNonPrintable(selection.Text())
		messages = append(messages, ReplaceSensitiveWords(text))
		logger.SugaredLogger.Infof("搜索到消息-%s: %s", msgType, text)
	})
	return &messages
}

func SearchStockInfoByCode(stock string) *[]string {
	// 创建一个 chromedp 上下文
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithLogf(logger.SugaredLogger.Infof),
		chromedp.WithErrorf(logger.SugaredLogger.Errorf),
	)
	defer cancel()
	var htmlContent string
	stock = strings.ReplaceAll(stock, "sh", "")
	stock = strings.ReplaceAll(stock, "sz", "")
	url := fmt.Sprintf("https://gushitong.baidu.com/stock/ab-%s", stock)
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		// 等待页面加载完成，可以根据需要调整等待时间
		//chromedp.Sleep(3*time.Second),
		chromedp.WaitVisible("a.news-item-link", chromedp.ByQuery),
		chromedp.OuterHTML("html", &htmlContent, chromedp.ByQuery),
	)
	if err != nil {
		logger.SugaredLogger.Error(err.Error())
		return &[]string{}
	}
	document, err := goquery.NewDocumentFromReader(strings.NewReader(htmlContent))
	if err != nil {
		logger.SugaredLogger.Error(err.Error())
		return &[]string{}
	}
	var messages []string
	document.Find("a.news-item-link").Each(func(i int, selection *goquery.Selection) {
		text := strutil.RemoveNonPrintable(selection.Text())
		if strings.Contains(text, stock) {
			messages = append(messages, text)
			logger.SugaredLogger.Infof("搜索到消息: %s", text)
		}
	})
	return &messages
}
