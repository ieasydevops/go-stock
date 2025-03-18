package models

import (
	"time"
)

// @Author spark
// @Date 2025/2/6 15:25
// @Desc
//-----------------------------------------------------------------------------------

type DateTime struct {
	time.Time
}

func (t DateTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + t.Format("2006-01-02 15:04:05") + `"`), nil
}

func (t *DateTime) UnmarshalJSON(data []byte) error {
	// Remove quotes
	s := string(data)
	s = s[1 : len(s)-1]

	// Parse time
	parsedTime, err := time.Parse("2006-01-02 15:04:05", s)
	if err != nil {
		return err
	}

	t.Time = parsedTime
	return nil
}

type GitHubReleaseVersion struct {
	Url       string `json:"url"`
	AssetsUrl string `json:"assets_url"`
	UploadUrl string `json:"upload_url"`
	HtmlUrl   string `json:"html_url"`
	Id        int    `json:"id"`
	Author    struct {
		Login             string `json:"login"`
		Id                int    `json:"id"`
		NodeId            string `json:"node_id"`
		AvatarUrl         string `json:"avatar_url"`
		GravatarId        string `json:"gravatar_id"`
		Url               string `json:"url"`
		HtmlUrl           string `json:"html_url"`
		FollowersUrl      string `json:"followers_url"`
		FollowingUrl      string `json:"following_url"`
		GistsUrl          string `json:"gists_url"`
		StarredUrl        string `json:"starred_url"`
		SubscriptionsUrl  string `json:"subscriptions_url"`
		OrganizationsUrl  string `json:"organizations_url"`
		ReposUrl          string `json:"repos_url"`
		EventsUrl         string `json:"events_url"`
		ReceivedEventsUrl string `json:"received_events_url"`
		Type              string `json:"type"`
		UserViewType      string `json:"user_view_type"`
		SiteAdmin         bool   `json:"site_admin"`
	} `json:"author"`
	NodeId          string   `json:"node_id"`
	TagName         string   `json:"tag_name"`
	TargetCommitish string   `json:"target_commitish"`
	Name            string   `json:"name"`
	Draft           bool     `json:"draft"`
	Prerelease      bool     `json:"prerelease"`
	CreatedAt       DateTime `json:"created_at"`
	PublishedAt     DateTime `json:"published_at"`
	Assets          []struct {
		Url      string `json:"url"`
		Id       int    `json:"id"`
		NodeId   string `json:"node_id"`
		Name     string `json:"name"`
		Label    string `json:"label"`
		Uploader struct {
			Login             string `json:"login"`
			Id                int    `json:"id"`
			NodeId            string `json:"node_id"`
			AvatarUrl         string `json:"avatar_url"`
			GravatarId        string `json:"gravatar_id"`
			Url               string `json:"url"`
			HtmlUrl           string `json:"html_url"`
			FollowersUrl      string `json:"followers_url"`
			FollowingUrl      string `json:"following_url"`
			GistsUrl          string `json:"gists_url"`
			StarredUrl        string `json:"starred_url"`
			SubscriptionsUrl  string `json:"subscriptions_url"`
			OrganizationsUrl  string `json:"organizations_url"`
			ReposUrl          string `json:"repos_url"`
			EventsUrl         string `json:"events_url"`
			ReceivedEventsUrl string `json:"received_events_url"`
			Type              string `json:"type"`
			UserViewType      string `json:"user_view_type"`
			SiteAdmin         bool   `json:"site_admin"`
		} `json:"uploader"`
		ContentType        string   `json:"content_type"`
		State              string   `json:"state"`
		Size               int      `json:"size"`
		DownloadCount      int      `json:"download_count"`
		CreatedAt          DateTime `json:"created_at"`
		UpdatedAt          DateTime `json:"updated_at"`
		BrowserDownloadUrl string   `json:"browser_download_url"`
	} `json:"assets"`
	TarballUrl string `json:"tarball_url"`
	ZipballUrl string `json:"zipball_url"`
	Body       string `json:"body"`
	Tag        Tag    `json:"tag"`
	Commit     Commit `json:"commit"`
}

type Tag struct {
	Ref    string `json:"ref"`
	NodeId string `json:"node_id"`
	Url    string `json:"url"`
	Object struct {
		Sha  string `json:"sha"`
		Type string `json:"type"`
		Url  string `json:"url"`
	} `json:"object"`
}

type Commit struct {
	Sha     string `json:"sha"`
	NodeId  string `json:"node_id"`
	Url     string `json:"url"`
	HtmlUrl string `json:"html_url"`
	Author  struct {
		Name  string   `json:"name"`
		Email string   `json:"email"`
		Date  DateTime `json:"date"`
	} `json:"author"`
	Committer struct {
		Name  string   `json:"name"`
		Email string   `json:"email"`
		Date  DateTime `json:"date"`
	} `json:"committer"`
	Tree struct {
		Sha string `json:"sha"`
		Url string `json:"url"`
	} `json:"tree"`
	Message string `json:"message"`
	Parents []struct {
		Sha     string `json:"sha"`
		Url     string `json:"url"`
		HtmlUrl string `json:"html_url"`
	} `json:"parents"`
	Verification struct {
		Verified   bool        `json:"verified"`
		Reason     string      `json:"reason"`
		Signature  interface{} `json:"signature"`
		Payload    interface{} `json:"payload"`
		VerifiedAt interface{} `json:"verified_at"`
	} `json:"verification"`
}

type AIResponseResult struct {
	ID        uint     `gorm:"primarykey"`
	CreatedAt DateTime `gorm:"column:created_at"`
	UpdatedAt DateTime `gorm:"column:updated_at"`
	DeletedAt bool     `gorm:"softDelete:flag"`
	ChatId    string   `json:"chatId"`
	ModelName string   `json:"modelName"`
	StockCode string   `json:"stockCode"`
	StockName string   `json:"stockName"`
	Question  string   `json:"question"`
	Content   string   `json:"content"`
	IsDel     bool     `gorm:"softDelete:flag"`
}

func (receiver AIResponseResult) TableName() string {
	return "ai_response_result"
}

type VersionInfo struct {
	ID             uint     `gorm:"primarykey"`
	CreatedAt      DateTime `gorm:"column:created_at"`
	UpdatedAt      DateTime `gorm:"column:updated_at"`
	DeletedAt      bool     `gorm:"softDelete:flag"`
	Version        string   `json:"version"`
	Content        string   `json:"content"`
	Icon           string   `json:"icon"`
	Alipay         string   `json:"alipay"`
	Wxpay          string   `json:"wxpay"`
	BuildTimeStamp int64    `json:"buildTimeStamp"`
	IsDel          bool     `gorm:"softDelete:flag"`
}

func (receiver VersionInfo) TableName() string {
	return "version_info"
}

type StockInfoHK struct {
	ID        uint     `gorm:"primarykey"`
	CreatedAt DateTime `gorm:"column:created_at"`
	UpdatedAt DateTime `gorm:"column:updated_at"`
	DeletedAt bool     `gorm:"softDelete:flag"`
	Code      string   `json:"code"`
	Name      string   `json:"name"`
	FullName  string   `json:"fullName"`
	EName     string   `json:"eName"`
	IsDel     bool     `gorm:"softDelete:flag"`
}

func (receiver StockInfoHK) TableName() string {
	return "stock_base_info_hk"
}

type StockInfoUS struct {
	ID        uint     `gorm:"primarykey"`
	CreatedAt DateTime `gorm:"column:created_at"`
	UpdatedAt DateTime `gorm:"column:updated_at"`
	DeletedAt bool     `gorm:"softDelete:flag"`
	Code      string   `json:"code"`
	Name      string   `json:"name"`
	FullName  string   `json:"fullName"`
	EName     string   `json:"eName"`
	Exchange  string   `json:"exchange"`
	Type      string   `json:"type"`
	IsDel     bool     `gorm:"softDelete:flag"`
}

func (receiver StockInfoUS) TableName() string {
	return "stock_base_info_us"
}

type Resp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// IsHKTradingTime 判断当前是否为港股交易时间
func IsHKTradingTime() bool {
	now := DateTime{Time: time.Now()}
	hour := now.Hour()
	minute := now.Minute()

	// 港股交易时间：
	// 开市前竞价时段：09:00 - 09:30
	if hour == 9 && minute >= 0 && minute <= 30 {
		return true
	}

	// 上午持续交易时段：09:30 - 12:00
	if (hour == 9 && minute > 30) || (hour >= 10 && hour < 12) || (hour == 12 && minute == 0) {
		return true
	}

	// 下午持续交易时段：13:00 - 16:00
	if (hour == 13 && minute >= 0) || (hour >= 14 && hour < 16) || (hour == 16 && minute == 0) {
		return true
	}

	// 收市竞价交易时段：16:00 - 16:10
	if hour == 16 && minute >= 0 && minute <= 10 {
		return true
	}

	return false
}

// IsUSTradingTime 判断当前是否为美股交易时间
func IsUSTradingTime() bool {
	// 获取美国东部时区
	est, err := time.LoadLocation("America/New_York")
	var estTime DateTime
	if err != nil {
		estTime = DateTime{Time: time.Now().Add(time.Hour * -12)}
	} else {
		// 将当前时间转换为美国东部时间
		estTime = DateTime{Time: time.Now().In(est)}
	}

	// 判断是否是周末
	weekday := estTime.Weekday()
	if weekday == time.Saturday || weekday == time.Sunday {
		return false
	}

	hour := estTime.Hour()
	minute := estTime.Minute()

	// 美股交易时间：
	// 开市前竞价时段：04:00 - 09:30
	if (hour == 4 && minute >= 0) || (hour > 4 && hour < 9) || (hour == 9 && minute <= 30) {
		return true
	}

	// 持续交易时段：09:30 - 16:00
	if (hour == 9 && minute > 30) || (hour > 9 && hour < 16) || (hour == 16 && minute == 0) {
		return true
	}

	// 收市后交易时段：16:00 - 20:00
	if (hour == 16 && minute > 0) || (hour > 16 && hour < 20) || (hour == 20 && minute == 0) {
		return true
	}

	return false
}
