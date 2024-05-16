package router

import (
	v1 "github.com/edgehook/ithings/webserver/api/v1"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()
	apiv1 := r.Group("/v1")
	{
		apiv1.POST("/awake/:mac", v1.AwakeDevice)
	}
	return r

}
