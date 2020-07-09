package reactor

import (
	"github.com/jinzhu/gorm"
	jsoniter "github.com/json-iterator/go"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/models"
)

type ClientStore struct {
	db *gorm.DB
}

func NewClientStore() *ClientStore {
	return &ClientStore{mainDB}
}

func (s *ClientStore) toClientInfo(data string) (oauth2.ClientInfo, error) {
	var cm models.Client
	err := jsoniter.Unmarshal([]byte(data), &cm)
	return &cm, err
}

// GetByID retrieves and returns client information by id
func (s *ClientStore) GetByID(id string) (oauth2.ClientInfo, error) {
	if id == "" {
		return nil, nil
	}
	if id != "c83b3a16"{
		return nil, errors.ErrServerError
	}
	// 在这里进行校验client_id
	//item := structs.OAuthApp{}
	//if result := mainDB.Where("client_id=? and status_id=?", id, 1).Order("created_at desc").First(&item); result.Error != nil {
	//	return nil, result.Error
	//}

	// 校验通过后，使用item中的数据，来生成oauth结构体
	//return s.toClientInfo(item.Data)
	data := `{"ID":"c83b3a16","Secret":"3a19c419","Domain":"http://localhost:19090/","UserID":null}`
	return s.toClientInfo(data)
}

// Create creates and stores the new client information
//func (s *ClientStore) Create(info oauth2.ClientInfo) error {
//	data, err := jsoniter.Marshal(info)
//	if err != nil {
//		return err
//	}
//
//	c := structs.ClientStoreItem{}
//
//	c.ClientID = info.GetID()
//	c.Secret = info.GetSecret()
//	c.Domain = info.GetDomain()
//	c.Data = string(data)
//
//	dbr := mainDB.Create(&c)
//
//	return dbr.Error
//}
