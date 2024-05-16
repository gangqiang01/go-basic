package api

import (
	responce "github.com/edgehook/ithings/webserver/types"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	// type auth struct {
	// 	Username string `form:"username" json:"username"`
	// 	Password string `form:"password" json:"password"`
	// }

	type resp struct {
		AccessToken string `form:"accessToken" json:"accessToken"`
	}
	// var authInfo auth
	// if err := c.Bind(&authInfo); err != nil {
	// 	responce.FailWithMessage("Parameter error", c)
	// 	return
	// }

	// body, err := json.Marshal(authInfo)
	// if err != nil {
	// 	responce.FailWithMessage("Json Marshal with  err %v", c)
	// 	return
	// }
	// appHubResponse, err := isync.SendMsgToAppHub("login", string(body))

	// if err != nil {
	// 	responce.FailWithMessage(fmt.Sprintf("Send msg to AppHub with  err %v", err), c)
	// 	return
	// }

	// if appHubResponse.StatusCode != "200" {
	// 	responce.FailWithMessage(fmt.Sprintf("Send msg to AppHub with  err %v", appHubResponse.Msg), c)
	// 	return
	// }

	// responce.OkWithData(&resp{
	// 	AccessToken: appHubResponse.Msg,
	// }, c)
	responce.OkWithData(&resp{
		AccessToken: "",
	}, c)
}
