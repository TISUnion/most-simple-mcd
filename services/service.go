package services

import (
	"fmt"
	"github.com/TISUnion/most-simple-mcd/constant"
	"github.com/TISUnion/most-simple-mcd/grpc/api"
	"github.com/TISUnion/most-simple-mcd/modules"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 登陆中间件
func AuthMiddle(c *gin.Context) {
	token := c.GetHeader(constant.TOKEN_HEADER_NAME)
	dbtoken := modules.GetFromDatabase(constant.DEFAULT_TOKEN_DB_KEY)
	if token == dbtoken && dbtoken != "" {
		c.Next()
	} else {
		c.JSON(http.StatusOK, getResponse(constant.TOKEN_FAILED, constant.TOKEN_FAILED_MESSAGE, ""))
		c.Abort()
	}
}

func RegisterServices() {
	ginEngine := modules.GetGinServerInstanceRouter()
	api.RegisterServerAuthMiddleware(AuthMiddle)
	api.RegisterAdminAuthMiddleware(AuthMiddle)
	api.RegisterUserAuthMiddleware(AuthMiddle)
	api.RegisterServerMcServerGinServer(ginEngine, &ServerService{})
	api.RegisterAdminAdminGinServer(ginEngine, &AdminService{})
	api.RegisterUserUserGinServer(ginEngine, &UserService{})
}

// 返回数据格式化
func getResponse(code int, message string, data interface{}) gin.H {
	responseData := make(gin.H)
	responseData["code"] = code
	responseData["message"] = message
	responseData["data"] = data
	return responseData
}

// 格式化错误信息
func errorFormat(err error) string {
	return fmt.Sprintf("GRPC error: %v", err)
}