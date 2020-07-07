package main

import (
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

const client_id string = "c83b3a16"
const client_secret string  = "3a19c419"


var redirect_ui = "http://localhost:19090/oauth/redirect"
var server_oauth_user_info = "http://127.0.0.1:1102/oauth/user"
//var server_oauth_user_info = "http://47.100.111.232:1102/oauth/user"
var server_get_token = ""

type AccessToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn uint `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope string `json:"scope"`
	TokenType string `json:"token_type"`
}

type UserData struct {
	Avatar string `json:"avatar"`
	Email string `json:"email"`
	Id uint `json:"id"`
	Intro string `json:"intro"`
	Mobile string `json:"mobile"`
	Username string `json:"username"`
}

type UserInfo struct {
	Result uint `json:"result"`
	UserData `json:"data"`
}

func main()  {
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")

	store := cookie.NewStore([]byte("loginuser"))
	router.Use(sessions.Sessions("mysession", store))

	router.GET("/", index)
	router.GET("/oauth/redirect", oauth)
	router.GET("/welcome", welcome)

	router.Run(":19090")
}

func index(c *gin.Context)  {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Login",
	})
}

func welcome(c *gin.Context) {
	session := sessions.Default(c)
	v := session.Get("username")
	c.HTML(http.StatusOK, "welcome.html", map[string]interface{}{
		"username": v,
	})
}

func oauth(c *gin.Context)  {
	if authorizationCode, ok := c.GetQuery("code"); ok {
		if authorizationCode == "" {
			c.JSON(http.StatusOK, gin.H{
				"error": "no authorization code",
			})
		} else {
			server_get_token = ""
			//server_get_token += "http://47.100.111.232:1102/token?grant_type=authorization_code&"
			server_get_token += "http://127.0.0.1:1102/token?grant_type=authorization_code&"
			server_get_token += "client_id="
			server_get_token += client_id
			server_get_token += "&client_secret="
			server_get_token += client_secret
			server_get_token += "&code="
			server_get_token += authorizationCode
			server_get_token += "&scope=read&"
			server_get_token += "redirect_uri="
			server_get_token += redirect_ui

			response, err := http.Get(server_get_token)
			defer response.Body.Close()

			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"error": "GET Request error",
				})
			}

			var response_bytes []byte

			if response_bytes, err = ioutil.ReadAll(response.Body); err != nil {
				log.Println(err)
			}

			var access_token_struct = &AccessToken{}

			err = json.Unmarshal(response_bytes, access_token_struct)
			if err != nil {
				log.Println("Error to decode access_token", err)
			}
			access_token := access_token_struct.AccessToken

			var bearer = "Bearer " + access_token

			userInfoReq, err := http.NewRequest("GET", server_oauth_user_info, nil)
			userInfoReq.Header.Add("Authorization", bearer)

			// Send req using http Client
			client := &http.Client{}
			resp, err := client.Do(userInfoReq)
			if err != nil {
				log.Println("Error on response.\n[ERRO] -", err)
			}

			var userInfo = &UserInfo{}
			var response_bytes2 []byte

			if response_bytes2, err = ioutil.ReadAll(resp.Body); err != nil {
				log.Println("Error on get user info ", err)
			}
			//log.Println(string([]byte(body)))

			err = json.Unmarshal(response_bytes2, userInfo)
			if err != nil {
				log.Println("Error on userinfo json data ", err)
			}

			//c.JSON(http.StatusOK, gin.H{
			//	"Welcome": userInfo.Username, })
			session := sessions.Default(c)
			session.Set("username", userInfo.Username)
			session.Save()
			c.Redirect(http.StatusPermanentRedirect, "http://localhost:19090/welcome")
		}
	}

}

