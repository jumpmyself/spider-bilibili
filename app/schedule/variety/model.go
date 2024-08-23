package variety

import "time"

// Variety undefined
type Variety struct {
	ID          int64     `json:"id" gorm:"id"`
	Title       string    `json:"title" gorm:"title"`
	URL         string    `json:"url" gorm:"url"`
	Updatedver  int64     `json:"updatedver" gorm:"updatedver"`
	TotalPlay   string    `json:"total_play" gorm:"total_play"`
	TotalDanmu  string    `json:"total_danmu" gorm:"total_danmu"`
	TotalPeople string    `json:"total_people" gorm:"total_people"`
	Style       string    `json:"style" gorm:"style"`
	PlayTime    string    `json:"play_time" gorm:"play_time"`
	Status      string    `json:"status" gorm:"status"`
	Avator      string    `json:"avator" gorm:"avator"`
	Desc        string    `json:"desc" gorm:"desc"`
	CreatedTime time.Time `json:"created_time" gorm:"created_time"`
	UpdatedTime time.Time `json:"updated_time" gorm:"updated_time"`
}

// TableName 表名称
func (*Variety) TableName() string {
	return "variety"
}

type AutoGenerated struct {
	Code int `json:"code"`
	Data struct {
		HasNext int `json:"has_next"`
		List    []struct {
			Badge     string `json:"badge"`
			BadgeInfo struct {
				BgColor      string `json:"bg_color"`
				BgColorNight string `json:"bg_color_night"`
				Text         string `json:"text"`
			} `json:"badge_info"`
			BadgeType int    `json:"badge_type"`
			Cover     string `json:"cover"`
			FirstEp   struct {
				Cover string `json:"cover"`
				EpID  int    `json:"ep_id"`
			} `json:"first_ep"`
			IndexShow    string `json:"index_show"`
			IsFinish     int    `json:"is_finish"`
			Link         string `json:"link"`
			MediaID      int    `json:"media_id"`
			Order        string `json:"order"`
			OrderType    string `json:"order_type"`
			Score        string `json:"score"`
			SeasonID     int    `json:"season_id"`
			SeasonStatus int    `json:"season_status"`
			SeasonType   int    `json:"season_type"`
			SubTitle     string `json:"subTitle"`
			Title        string `json:"title"`
			TitleIcon    string `json:"title_icon"`
		} `json:"list"`
		Num   int `json:"num"`
		Size  int `json:"size"`
		Total int `json:"total"`
	} `json:"data"`
	Message string `json:"message"`
}
