package data

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/duke-git/lancet/v2/strutil"
	"github.com/go-resty/resty/v2"
	"github.com/robertkrimen/otto"
	"github.com/samber/lo"
	"go-stock/backend/db"
	"go-stock/backend/logger"
	"go-stock/backend/models"
	"strconv"
	"strings"
	"time"
)

// @Author spark
// @Date 2025/4/23 14:54
// @Desc
// -----------------------------------------------------------------------------------
type MarketNewsApi struct {
}

func NewMarketNewsApi() *MarketNewsApi {
	return &MarketNewsApi{}
}

func (m MarketNewsApi) GetNewTelegraph(crawlTimeOut int64) *[]models.Telegraph {
	url := "https://www.cls.cn/telegraph"
	response, _ := resty.New().SetTimeout(time.Duration(crawlTimeOut)*time.Second).R().
		SetHeader("Referer", "https://www.cls.cn/").
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36 Edg/117.0.2045.60").
		Get(fmt.Sprintf(url))
	var telegraphs []models.Telegraph
	//logger.SugaredLogger.Info(string(response.Body()))
	document, _ := goquery.NewDocumentFromReader(strings.NewReader(string(response.Body())))

	document.Find(".telegraph-list").Each(func(i int, selection *goquery.Selection) {
		//logger.SugaredLogger.Info(selection.Text())
		telegraph := models.Telegraph{Source: "财联社电报"}
		spans := selection.Find("div.telegraph-content-box span")
		if spans.Length() == 2 {
			telegraph.Time = spans.First().Text()
			telegraph.Content = spans.Last().Text()
			if spans.Last().HasClass("c-de0422") {
				telegraph.IsRed = true
			}
		}

		labels := selection.Find("div a.label-item")
		labels.Each(func(i int, selection *goquery.Selection) {
			if selection.HasClass("link-label-item") {
				telegraph.Url = selection.AttrOr("href", "")
			} else {
				tag := &models.Tags{
					Name: selection.Text(),
					Type: "subject",
				}
				db.Dao.Model(tag).Where("name=? and type=?", selection.Text(), "subject").FirstOrCreate(&tag)
				telegraph.SubjectTags = append(telegraph.SubjectTags, selection.Text())
			}
		})
		stocks := selection.Find("div.telegraph-stock-plate-box a")
		stocks.Each(func(i int, selection *goquery.Selection) {
			telegraph.StocksTags = append(telegraph.StocksTags, selection.Text())
		})

		//telegraph = append(telegraph, ReplaceSensitiveWords(selection.Text()))
		if telegraph.Content != "" {
			cnt := int64(0)
			db.Dao.Model(telegraph).Where("time=? and source=?", telegraph.Time, telegraph.Source).Count(&cnt)
			if cnt == 0 {
				db.Dao.Create(&telegraph)
				telegraphs = append(telegraphs, telegraph)
				for _, tag := range telegraph.SubjectTags {
					tagInfo := &models.Tags{}
					db.Dao.Model(models.Tags{}).Where("name=? and type=?", tag, "subject").First(&tagInfo)
					if tagInfo.ID > 0 {
						db.Dao.Model(models.TelegraphTags{}).Where("telegraph_id=? and tag_id=?", telegraph.ID, tagInfo.ID).FirstOrCreate(&models.TelegraphTags{
							TelegraphId: telegraph.ID,
							TagId:       tagInfo.ID,
						})
					}
				}
			}

		}
	})
	return &telegraphs
}
func (m MarketNewsApi) GetNewsList(source string, limit int) *[]*models.Telegraph {
	news := &[]*models.Telegraph{}
	if source != "" {
		db.Dao.Model(news).Preload("TelegraphTags").Where("source=?", source).Order("id desc").Limit(limit).Find(news)
	} else {
		db.Dao.Model(news).Preload("TelegraphTags").Order("id desc").Limit(limit).Find(news)
	}
	for _, item := range *news {
		tags := &[]models.Tags{}
		db.Dao.Model(&models.Tags{}).Where("id in ?", lo.Map(item.TelegraphTags, func(item models.TelegraphTags, index int) uint {
			return item.TagId
		})).Find(&tags)
		tagNames := lo.Map(*tags, func(item models.Tags, index int) string {
			return item.Name
		})
		item.SubjectTags = tagNames
		logger.SugaredLogger.Infof("tagNames %v ，SubjectTags：%s", tagNames, item.SubjectTags)
	}
	return news
}
func (m MarketNewsApi) GetTelegraphList(source string) *[]*models.Telegraph {
	news := &[]*models.Telegraph{}
	if source != "" {
		db.Dao.Model(news).Preload("TelegraphTags").Where("source=?", source).Order("id desc").Limit(20).Find(news)
	} else {
		db.Dao.Model(news).Preload("TelegraphTags").Order("id desc").Limit(20).Find(news)
	}
	for _, item := range *news {
		tags := &[]models.Tags{}
		db.Dao.Model(&models.Tags{}).Where("id in ?", lo.Map(item.TelegraphTags, func(item models.TelegraphTags, index int) uint {
			return item.TagId
		})).Find(&tags)
		tagNames := lo.Map(*tags, func(item models.Tags, index int) string {
			return item.Name
		})
		item.SubjectTags = tagNames
		logger.SugaredLogger.Infof("tagNames %v ，SubjectTags：%s", tagNames, item.SubjectTags)
	}
	return news
}

func (m MarketNewsApi) GetSinaNews(crawlTimeOut uint) *[]models.Telegraph {
	news := &[]models.Telegraph{}
	response, _ := resty.New().SetTimeout(time.Duration(crawlTimeOut)*time.Second).R().
		SetHeader("Referer", "https://finance.sina.com.cn").
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36 Edg/117.0.2045.60").
		Get("https://zhibo.sina.com.cn/api/zhibo/feed?callback=callback&page=1&page_size=20&zhibo_id=152&tag_id=0&dire=f&dpc=1&pagesize=20&id=4161089&type=0&_=" + strconv.FormatInt(time.Now().Unix(), 10))
	js := string(response.Body())
	js = strutil.ReplaceWithMap(js,
		map[string]string{
			"try{callback(":  "var data=",
			");}catch(e){};": ";",
		})
	//logger.SugaredLogger.Info(js)
	vm := otto.New()
	_, err := vm.Run(js)
	if err != nil {
		logger.SugaredLogger.Error(err)
	}
	vm.Run("var result = data.result;")
	//vm.Run("var resultStr =JSON.stringify(data);")
	vm.Run("var resultData = result.data;")
	vm.Run("var feed = resultData.feed;")
	vm.Run("var feedStr = JSON.stringify(feed);")

	value, _ := vm.Get("feedStr")
	//resultStr, _ := vm.Get("resultStr")

	//logger.SugaredLogger.Info(resultStr)
	feed := make(map[string]any)
	err = json.Unmarshal([]byte(value.String()), &feed)
	if err != nil {
		logger.SugaredLogger.Errorf("json.Unmarshal error:%v", err.Error())
	}
	var telegraphs []models.Telegraph

	if feed["list"] != nil {
		for _, item := range feed["list"].([]any) {
			telegraph := models.Telegraph{Source: "新浪财经"}
			data := item.(map[string]any)
			//logger.SugaredLogger.Infof("%s:%s", data["create_time"], data["rich_text"])
			telegraph.Content = data["rich_text"].(string)
			telegraph.Time = strings.Split(data["create_time"].(string), " ")[1]
			tags := data["tag"].([]any)
			telegraph.SubjectTags = lo.Map(tags, func(tagItem any, index int) string {
				name := tagItem.(map[string]any)["name"].(string)
				tag := &models.Tags{
					Name: name,
					Type: "sina_subject",
				}
				db.Dao.Model(tag).Where("name=? and type=?", name, "sina_subject").FirstOrCreate(&tag)
				return name
			})
			if _, ok := lo.Find(telegraph.SubjectTags, func(item string) bool { return item == "焦点" }); ok {
				telegraph.IsRed = true
			}
			logger.SugaredLogger.Infof("telegraph.SubjectTags:%v %s", telegraph.SubjectTags, telegraph.Content)

			if telegraph.Content != "" {
				cnt := int64(0)
				db.Dao.Model(telegraph).Where("time=? and source=?", telegraph.Time, telegraph.Source).Count(&cnt)
				if cnt == 0 {
					db.Dao.Create(&telegraph)
					telegraphs = append(telegraphs, telegraph)
					for _, tag := range telegraph.SubjectTags {
						tagInfo := &models.Tags{}
						db.Dao.Model(models.Tags{}).Where("name=? and type=?", tag, "sina_subject").First(&tagInfo)
						if tagInfo.ID > 0 {
							db.Dao.Model(models.TelegraphTags{}).Where("telegraph_id=? and tag_id=?", telegraph.ID, tagInfo.ID).FirstOrCreate(&models.TelegraphTags{
								TelegraphId: telegraph.ID,
								TagId:       tagInfo.ID,
							})
						}
					}
				}
			}
		}
		return &telegraphs
	}

	return news

}

func (m MarketNewsApi) GlobalStockIndexes(crawlTimeOut uint) map[string]any {
	response, _ := resty.New().SetTimeout(time.Duration(crawlTimeOut)*time.Second).R().
		SetHeader("Referer", "https://stockapp.finance.qq.com/mstats").
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36 Edg/117.0.2045.60").
		Get("https://proxy.finance.qq.com/ifzqgtimg/appstock/app/rank/indexRankDetail2")
	js := string(response.Body())
	res := make(map[string]any)
	json.Unmarshal([]byte(js), &res)
	return res["data"].(map[string]any)
}

func (m MarketNewsApi) GetIndustryRank(sort string, cnt int) map[string]any {

	url := fmt.Sprintf("https://proxy.finance.qq.com/ifzqgtimg/appstock/app/mktHs/rank?l=%d&p=1&t=01/averatio&ordertype=&o=%s", cnt, sort)
	response, _ := resty.New().SetTimeout(time.Duration(5)*time.Second).R().
		SetHeader("Referer", "https://stockapp.finance.qq.com/").
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36 Edg/117.0.2045.60").
		Get(url)
	js := string(response.Body())
	res := make(map[string]any)
	json.Unmarshal([]byte(js), &res)
	return res
}

func (m MarketNewsApi) GetIndustryMoneyRankSina(fenlei string) []map[string]any {
	url := fmt.Sprintf("https://vip.stock.finance.sina.com.cn/quotes_service/api/json_v2.php/MoneyFlow.ssl_bkzj_bk?page=1&num=20&sort=netamount&asc=0&fenlei=%s", fenlei)

	response, _ := resty.New().SetTimeout(time.Duration(5)*time.Second).R().
		SetHeader("Host", "vip.stock.finance.sina.com.cn").
		SetHeader("Referer", "https://finance.sina.com.cn").
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36 Edg/117.0.2045.60").
		Get(url)
	js := string(response.Body())
	res := &[]map[string]any{}
	err := json.Unmarshal([]byte(js), &res)
	if err != nil {
		logger.SugaredLogger.Error(err)
		return *res
	}
	return *res
}

func (m MarketNewsApi) GetMoneyRankSina(sort string) []map[string]any {
	if sort == "" {
		sort = "netamount"
	}
	url := fmt.Sprintf("https://vip.stock.finance.sina.com.cn/quotes_service/api/json_v2.php/MoneyFlow.ssl_bkzj_ssggzj?page=1&num=20&sort=%s&asc=0&bankuai=&shichang=", sort)
	response, _ := resty.New().SetTimeout(time.Duration(5)*time.Second).R().
		SetHeader("Host", "vip.stock.finance.sina.com.cn").
		SetHeader("Referer", "https://finance.sina.com.cn").
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36 Edg/117.0.2045.60").
		Get(url)
	js := string(response.Body())
	res := &[]map[string]any{}
	err := json.Unmarshal([]byte(js), &res)
	if err != nil {
		logger.SugaredLogger.Error(err)
		return *res
	}
	return *res
}

func (m MarketNewsApi) GetStockMoneyTrendByDay(stockCode string, days int) []map[string]any {
	url := fmt.Sprintf("http://vip.stock.finance.sina.com.cn/quotes_service/api/json_v2.php/MoneyFlow.ssl_qsfx_zjlrqs?page=1&num=%d&sort=opendate&asc=0&daima=%s", days, stockCode)

	response, _ := resty.New().SetTimeout(time.Duration(5)*time.Second).R().
		SetHeader("Host", "vip.stock.finance.sina.com.cn").
		SetHeader("Referer", "https://finance.sina.com.cn").
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36 Edg/117.0.2045.60").Get(url)
	js := string(response.Body())
	res := &[]map[string]any{}
	err := json.Unmarshal([]byte(js), &res)
	if err != nil {
		logger.SugaredLogger.Error(err)
		return *res
	}
	return *res

}
