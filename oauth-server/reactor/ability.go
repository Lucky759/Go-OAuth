package reactor

import (
	_ "github.com/jinzhu/gorm/dialects/postgres" //
	"github.com/sirupsen/logrus"
	"oauth-server/config"
	"oauth-server/processor"
	"time"

	"github.com/jinzhu/gorm"
)

var mainDB *gorm.DB

var Redis processor.Redis

func initPostgres() {
	postgresqlURL := config.Get("persistence.postgresql.url").String()
	postgresqlShowLog := config.Get("persistence.postgresql.showLog").Bool()
	var err error
	mainDB, err = InitPostgres(postgresqlURL, postgresqlShowLog)
	if err != nil {
		logrus.Error(err)
		time.Sleep(5 * time.Second)
		initPostgres()
	}
}

// InitPostgres InitPostgres
func InitPostgres(dbURL string, showLog bool) (*gorm.DB, error) {
	db, err := gorm.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	db.DB().SetMaxOpenConns(20)
	db.DB().SetMaxIdleConns(20)
	db.LogMode(showLog)
	return db, err
}

// Init Init
func Init() {
	initPostgres()
	redisHost := config.Get("persistence.redis.host").String()
	Redis.Init(redisHost)

}
