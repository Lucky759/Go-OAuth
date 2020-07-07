package routers

import (
	"fmt"
	"net/http"
	"oauth-server/processor"
	"oauth-server/reactor"
	"oauth-server/structs"
	"strconv"
	"strings"
	"github.com/gin-gonic/gin"

	"gopkg.in/oauth2.v3/server"
)

func GetAuthorized(srv *server.Server) gin.HandlerFunc {

	fn := func(c *gin.Context) {
		err := srv.HandleAuthorizeRequest(c.Writer, c.Request)
		if err != nil {
			c.PureJSON(http.StatusBadRequest, structs.Response{Result: 4, Message: err.Error(), Data: c.Writer})
			return
			//http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		}
	}
	return fn
}

func GetToken(srv *server.Server) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		err := srv.HandleTokenRequest(c.Writer, c.Request)
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		}
	}
	return fn
}

func GetProtectedInfo(c *gin.Context) {
	if queryParam, ok := c.GetQuery("code"); ok {
		if queryParam == "" {
			c.JSON(http.StatusOK, gin.H{
				"error": "no authorization code",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{"code": queryParam})
		}
	}
}

// -------------------- OAuthLogin --------------------

type oAuthLoginReq struct {
	Mobile   string `json:"mobile"   example:"18621586899"`
	Password string `json:"password"   example:"xY27EqHpMsXkrFuz"`
}

func oAuthLoginVerify(c *gin.Context) (*oAuthLoginReq, error) {
	req := oAuthLoginReq{}
	c.ShouldBindJSON(&req)
	req.Mobile = strings.TrimSpace(req.Mobile)
	req.Password = strings.TrimSpace(req.Password)
	//req.Captcha = strings.TrimSpace(req.Captcha)
	v := len(req.Mobile)*len(req.Password) > 0
	//v := len(req.Mobile) * (len(req.Password)+len(req.Captcha)) > 0
	if v {
		return &req, nil
	}
	return &req, fmt.Errorf("参数异常")
}

func OAuthUserLogin(srv *server.Server) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		req, err := oAuthLoginVerify(c)
		if err != nil {
			c.PureJSON(http.StatusOK, structs.Response{Result: 4, Message: "参数异常"})
			return
		}
		user := (*structs.User)(nil)

		user, _ = reactor.UserFindWithMobile(req.Mobile)
		if user == nil {
			c.PureJSON(http.StatusOK, structs.Response{Result: 101001, Message: "用户不存在"})
			return
		}
		if !processor.CompareEncodedPassword(user.Password, req.Password) {
			c.PureJSON(http.StatusOK, structs.Response{Result: 101002, Message: "密码不正确"})
			return
		}

		if user.StatusID == 2 {
			c.PureJSON(http.StatusOK, structs.Response{Result: 2, Message: "用户被禁用"})
			return
		}
		oAuthLogin(c, user, srv)
	}

	return gin.HandlerFunc(fn)
}

func oAuthLogin(c *gin.Context, user *structs.User, srv *server.Server) {
	// Set UserAuthorizationHandler
	srv.SetUserAuthorizationHandler(
		func(w http.ResponseWriter, r *http.Request) (userID string, err error) {
			return strconv.FormatUint(uint64(user.ID), 10), nil
		})

	//token, err := processor.GenerateJWTToken(map[string]interface{}{
	//	"userID": user.ID,
	//	// "time":   time.Now().Unix(),
	//})
	//if err != nil {
	//	logrus.Error(err)
	//}z`
	//
	//_, err = reactor.Redis.Exec(0, func(conn redis.Conn) (interface{}, error) {
	//	d := map[string]interface{}{
	//		"user":  user,
	//		"token": token,
	//	}
	//	buf, _ := json.Marshal(d)
	//	return conn.Do("HSET", "users", user.ID, buf)
	//})
	//if err != nil {
	//	logrus.Error(err)
	//}

	c.PureJSON(http.StatusOK, structs.Response{
		Result: 0,
		Data: map[string]interface{}{
			"user_id": user.ID,
		},
	})
}

// -------------------- OAuthGetUserInfo --------------------

type oAuthGetUserInfoReq struct {
	//UserID uint `json:"user_id" binding:"required"`
	UserID   uint   `form:"userID" binding:"required" example:"1"`
}

func oAuthGetUserInfoVerify(c *gin.Context) (*oAuthGetUserInfoReq, error) {
	var req oAuthGetUserInfoReq
	err := c.ShouldBind(&req)
	if err != nil {
		return nil, err
	}
	v := req.UserID > 0
	if v {
		return &req, nil
	}
	return &req, fmt.Errorf("参数异常")
}

// OAuthGetUserInfo OAuthGetUserInfo
// @ID OAuthGetUserInfo
// @Tags OAuthApp
// @Summary OAuth获取用户信息
// @Description
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} structs.Response
// @Router /oauth/user [get]
func OAuthGetUserInfo(srv *server.Server) gin.HandlerFunc {
	fn := func(c *gin.Context) {
		t, err := srv.ValidationBearerToken(c.Request)
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusBadRequest)
			return
		}
		userIDStr := t.GetUserID()
		userIDInt, err := strconv.Atoi(userIDStr)
		userIDUint := uint(userIDInt)
		if err != nil{
			return
		}

		//s := srv.GetTokenData(t)

		//req, err := oAuthGetUserInfoVerify(c)
		//if err != nil {
		//	c.PureJSON(http.StatusOK, structs.Response{Result: constants.ResultParamsError, Message: constants.ResultParamsErrorExplain})
		//	return
		//}

		user, err := reactor.UserInfo(userIDUint)
		//User StatusID is not used
		if err != nil {
			c.PureJSON(http.StatusOK, structs.Response{Result: 2, Message: "数据库异常"})
			return
		}
		//if user == nil || user.ID != req.UserID {
		//	c.PureJSON(http.StatusOK, structs.Response{Result: constants.ResultAuthenticationFailed, Message: constants.ResultAuthenticationFailedExplain})
		//	return
		//}
		c.PureJSON(http.StatusOK, structs.Response{
			Result: 0,
			Data:   user,
		})
	}
	return fn
}
