package main

import (
	"flag"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	oredis "gopkg.in/go-oauth2/redis.v3"
	"gopkg.in/oauth2.v3/errors"
	"gopkg.in/oauth2.v3/generates"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/server"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"oauth-server/config"
	"oauth-server/middleware"
	"oauth-server/reactor"
	"oauth-server/routers"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

//go:generate go run ./cmd/constants-generator.go -s constants

// generate doc
// go get github.com/swaggo/swag/cmd/swag
//go:generate swag init -g docs/doc-generator.go

func loadConfig() []byte {
	configPtr := flag.String("c", "", "config file path")
	flag.Parse()

	configPath := strings.TrimSpace(*configPtr)
	if len(configPath) == 0 {
		panic("miss config file")
	}

	configPath, _ = filepath.Abs(configPath)
	buf, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	return buf
}


func OAuthInit() *server.Server {
	manager := manage.NewDefaultManager()
	//manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	redisAddr := config.Get("persistence.redis.host").String()
	redisDB := int(config.Get("persistence.redis.db").Int())

	manager.MapTokenStorage(oredis.NewRedisStore(&redis.Options{
		Addr: redisAddr,
		DB: redisDB,
	}))

	//clientStore := store.NewClientStore()
	clientStore := reactor.NewClientStore()

	manager.MapAccessGenerate(generates.NewJWTAccessGenerate([]byte("0350f134416954fd0529f191d3eb6f8a"), jwt.SigningMethodHS512))

	manager.MapClientStorage(clientStore)

	ginserver := server.NewDefaultServer(manager)

	ginserver.SetAllowGetAccessRequest(true)
	ginserver.SetClientInfoHandler(server.ClientFormHandler)

	ginserver.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})
	ginserver.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	ginserver.SetAccessTokenExpHandler(func(w http.ResponseWriter, r *http.Request) (exp time.Duration, err error) {
		return time.Hour * 24 * 14, nil
	})

	return ginserver
}

func init() {
	rand.Seed(time.Now().UnixNano())
	config.Init(loadConfig())
	reactor.Init()
}

func main() {

	ginserver := OAuthInit()

	r := gin.Default()
	r.Use(middleware.AuthIdentify)
	r.GET("/test", func(context *gin.Context) {
		context.JSON(200, "请求成功")
	})
	r.GET("/authorize", routers.GetAuthorized(ginserver))
	r.GET("/token", routers.GetToken(ginserver))
	//r.GET("/callback", routers.GetProtectedInfo)
	r.POST("/oauth/login", routers.OAuthUserLogin(ginserver))
	r.GET("/oauth/user", routers.OAuthGetUserInfo(ginserver))
	listen := config.Get("app.listen").String()
	logrus.Info("listen on http://" + listen)
	r.Run(listen)
}
