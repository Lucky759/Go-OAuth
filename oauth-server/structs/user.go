package structs

import (
	"time"
)

// User User
type User struct {
	Model            `table-comment:"用户"`
	Username         string    `json:"username"      sql:"index" gorm:"type:varchar(100)"  comment:"用户名"`
	Mobile           string    `json:"mobile"        sql:"index" gorm:"type:varchar(100)"  comment:"手机号"`
	Email            string    `json:"email"         sql:"index" gorm:"type:varchar(100)"  comment:"邮箱"`
	Password         string    `json:"-"             sql:"index" gorm:"type:varchar(200)"  comment:"密码"`
	JoinAt           time.Time `json:"joinAt"                                              comment:"注册时间"`
	Avatar           string    `json:"avatar"                                              comment:"头像"`
	Intro            string    `json:"intro"                     grom:"type:varchar(8000)" comment:"简介"`
	StatusID         uint      `json:"statusID"                                            comment:"状态"`
	GiteaSSOToken    string    `json:"giteaSSOToken"                                       comment:"gitea sso token"`
	GiteaUsername    string    `json:"giteaUsername"                                       comment:"gitea username"`
	Experience       string    `json:"experience" 	gorm:"type:varchar(100) 				comment:"项目经验"`
	ParttimeStatusID uint      `json:"parttimeStatusId" comment:"空闲状态"`
}

// TableName TableName
func (*User) TableName() string {
	return "users"
}
