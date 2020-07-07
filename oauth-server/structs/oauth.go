package structs

// OAuth struct
type OAuthApp struct {
	Model         `table-comment:"OAuth认证"`
	AppName       string `json:"app_name"    gorm:"type:varchar(100)"  comment:"应用名称"`
	Description   string `json:"description" gorm:"type:varchar(100)" comment:"描述"`
	HomePage      string `json:"home_page" gorm:"type:varchar(100)" comment:"主页"`
	CallbackUrl   string `json:"callback_url" gorm:"type:varchar(100)" comment:"回调地址"`
	ClientID      string `json:"client_id" gorm:"type:varchar(100)" comment:"应用ID"`
	ClientSecret string `json:"client_secret" gorm:"type:varchar(100)" comment:"应用Secret"`
	StatusID      uint   `json:"status_id" comment:"启用禁用"`
	Data          string `json:"data" gorm:"type:varchar(4096)"`
}

// TableName
func (o *OAuthApp) TableName() string {
	return "oauth_apps"
}

//// ClientStoreItem data item
//type ClientStoreItem struct {
//	Model    `table-comment:"OAuthClient"`
//	ClientID string `json:"client_id" gorm:"type:varchar(512)" comment:"clientid"`
//	Secret   string `json:"secret" gorm:"type:varchar(512)" comment:"secrect"`
//	Domain   string `json:"domain" gorm:"type:varchar(512)" comment:"domain"`
//	Data     string `json:"data" gorm:"type:varchar(4096) comment:"data"`
//}
//
//// TableName
//func (s *ClientStoreItem) TableName() string {
//	return "oauth_client_store"
//}
