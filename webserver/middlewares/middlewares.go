package middlewares

import (
	"encoding/json"
	"github.com/edgehook/ithings/transport/isync"
	responce "github.com/edgehook/ithings/webserver/types"
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
	"net/http"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session, Content-Type, accesstoken, timeout, Srptoken")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			c.Header("Access-Control-Max-Age", "3600")
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
			c.Abort()
			return
		}
		// path := c.FullPath()
		// if strings.Contains(path, "v1") && !strings.Contains(path, "export") && !strings.Contains(path, "download") && !strings.Contains(path, "SimpleJsons") {
		// 	token := c.Request.Header.Get("accesstoken")
		// 	srpToken := c.Request.Header.Get("Srptoken")
		// 	if !verifyToken("accesstoken", token) && !verifyToken("Srptoken", srpToken) {
		// 		responce.FailWithCodeAndMessage(401, "illegal user", c)
		// 		//stop context
		// 		c.Abort()
		// 		return
		// 	}
		// }

		defer func() {
			if err := recover(); err != nil {
				klog.Errorf("Panic info is: %v", err)
				responce.FailWithMessage("Server error", c)
			}
		}()

		c.Next()
	}
}

func verifyToken(header, token string) bool {
	type tokenData struct {
		Header string `form:"header" json:"header,omitempty"`
		Token  string `form:"token" json:"token,omitempty"`
	}

	data := &tokenData{
		Token:  token,
		Header: header,
	}

	body, err := json.Marshal(data)
	if err != nil {
		klog.Errorf("Json Marshal with  err %v", err)
		return false
	}
	appHubResponse, err := isync.SendMsgToAppHub("verifyToken", string(body))

	if err != nil {
		klog.Errorf("Send msg to AppHub with  err %v", err)
		return false
	}

	if appHubResponse.StatusCode == "200" {
		return true
	}
	return false
}
