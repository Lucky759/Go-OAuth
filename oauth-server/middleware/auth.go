package middleware

import (
	"encoding/json"
	"github.com/clintwan/armory"
	"net/http"
	"net/url"
	"oauth-server/processor"
	"oauth-server/reactor"
	"oauth-server/structs"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"

)

var whiteList = []string{
	"/constants",
	"/auth/join",
	"/auth/login",
}

func allowOrigin(c *gin.Context) {
	originURL, _ := url.Parse(c.Request.Header.Get("Origin"))
	c.Header("Access-Control-Allow-Origin", originURL.String())
	c.Header("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT")
	// c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
}

func logRequest(c *gin.Context) {
	// logrus.Debug("middleware auth identify")
	headerBuf, _ := json.MarshalIndent(map[string]interface{}{
		"method": c.Request.Method,
		"path":   c.Request.RequestURI,
		"header": c.Request.Header,
	}, "", "  ")
	logrus.Debug(string(headerBuf))
}

func inWhiteList(c *gin.Context) bool {
	r := false
	idx, _ := armory.Slice.IndexOf(whiteList, c.Request.RequestURI)
	if idx > -1 {
		r = true
	}
	return r
}

func loadUser(c *gin.Context) interface{} {
	requestToken := c.GetHeader("Authorization")
	m, err := processor.ParseJWTToken(requestToken)
	if err != nil {
		return nil
	}
	userID := m["userID"].(float64)

	data, rErr := reactor.Redis.Exec(0, func(conn redis.Conn) (interface{}, error) {
		return conn.Do("HGET", "users", userID)
	})
	if rErr != nil {
		logrus.Error(rErr)
		return nil
	}
	if data == nil {
		return nil
	}
	d := map[string]interface{}{}
	json.Unmarshal(data.([]byte), &d)
	cachedToken := d["token"]
	if requestToken != cachedToken {
		return nil
	}

	user := structs.User{}
	buf, _ := json.Marshal(d["user"])
	json.Unmarshal(buf, &user)
	return &user
}

// AuthIdentify AuthIdentify
func AuthIdentify(c *gin.Context) {
	// buf, _ := ioutil.ReadAll(c.Request.Body)
	// c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

	allowOrigin(c)

	// skip options
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}

	// logRequest(c)

	// if inWhiteList(c) {
	// 	c.Next()
	// 	return
	// }

	user := loadUser(c)
	if user != nil {
		c.Set("user", user)
	}

	c.Next()
}
