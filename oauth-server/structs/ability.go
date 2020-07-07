package structs

import "time"

// Response Response
type Response struct {
	Result  int         `json:"result"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// Model Model
type Model struct {
	ID        uint       `json:"id"                    gorm:"primary_key" comment:"自增主键"`
	CreatedAt time.Time  `json:"createdAt"                                comment:"创建于"`
	UpdatedAt time.Time  `json:"updatedAt"                                comment:"更新于"`
	DeletedAt *time.Time `json:"deletedAt" sql:"index"                    comment:"删除于"`
}
