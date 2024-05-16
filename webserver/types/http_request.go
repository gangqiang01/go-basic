package types

import (
	"encoding/json"

	"github.com/edgehook/ithings/transport/isync"
	"github.com/gin-gonic/gin"
	"k8s.io/klog/v2"
)

type TokenSubject struct {
	Username string `form:"username" json:"username"`
}

func GetUserName(c *gin.Context) string {
	type tokenData struct {
		Header string `form:"header" json:"header,omitempty"`
		Token  string `form:"token" json:"token,omitempty"`
	}
	token := c.Request.Header.Get("accesstoken")
	data := &tokenData{
		Token:  token,
		Header: "accesstoken",
	}

	body, err := json.Marshal(data)
	if err != nil {
		klog.Errorf("Json Marshal with  err %v", err)
		return ""
	}
	appHubResponse, err := isync.SendMsgToAppHub("verifyToken", string(body))

	if err != nil {
		klog.Errorf("Send msg to AppHub with  err %v", err)
		return ""
	}

	if appHubResponse.StatusCode == "200" {
		subject := TokenSubject{}

		err := json.Unmarshal([]byte(appHubResponse.Msg), &subject)
		if err != nil {
			klog.Errorf("json Unmarshal with err %v", err)
			return ""
		}
		return subject.Username
	}

	return ""
}
