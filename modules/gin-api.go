package modules

import (
	"encoding/json"
	"fmt"
	"github.com/TISUnion/most-simple-mcd/constant"
	json_struct "github.com/TISUnion/most-simple-mcd/json-struct"
	"github.com/TISUnion/most-simple-mcd/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

const (
	TOKEN_HEADER_NAME = "X-Token"
)

func RegisterRouter() {
	router := GetGinServerInstanceRouter()
	// 用户登陆
	router.POST("/user/login", userLogin)
	// 用户注销
	router.POST("/user/logout", userLogout)
	// 需要登陆才能请求的接口
	v1 :=
		router.Group("/api/v1")
	{
		// 登陆验证中间件
		v1.Use(func(c *gin.Context) {
			token := c.GetHeader(TOKEN_HEADER_NAME)
			dbtoken := GetFromDatabase(constant.DEFAULT_TOKEN_DB_KEY)
			if token == dbtoken && dbtoken != "" {
				c.Next()
			} else {
				c.JSON(http.StatusOK, getResponse(constant.TOKEN_FAILED, constant.TOKEN_FAILED_MESSAGE, ""))
				c.Abort()
			}
		})
		// 获取用户信息
		v1.GET("/user/info", getInfo)
		// 更新用户信息
		v1.PATCH("/user/account", updateUserData)
		// 获取配置
		v1.GET("/config/list", getConfig)
		// 修改配置
		v1.PATCH("/config", updateConfig)
	}
	// websocket实时监听服务端耗费资源
	router.GET("/server/resources/listen/:serverId", serversResourcesListen)
	// websocket实时获取服务器输出
	router.GET("/server/std/listen/:serverId", serversStdListen)
}

// 用户登录
func userLogin(c *gin.Context) {
	var reqInfo json_struct.AdminUser
	if err := c.BindJSON(&reqInfo); err != nil {
		WriteLogToDefault(errorFormat(err), constant.LOG_ERROR)
		c.JSON(http.StatusOK, getResponse(constant.HTTP_PARAMS_ERROR, constant.HTTP_PARAMS_ERROR_MESSAGE, ""))
		return
	}
	var adminObj json_struct.AdminUser
	adminJson := GetFromDatabase(constant.DEFAULT_ACCOUNT_DB_KEY)
	if adminJson == "" {
		adminObj = *setDefaultAccount()
	} else {
		if err := json.Unmarshal([]byte(adminJson), &adminObj); err != nil {
			WriteLogToDefault(errorFormat(err), constant.LOG_ERROR)
			c.JSON(http.StatusOK, getResponse(constant.HTTP_SYSTEM_ERROR, constant.HTTP_SYSTEM_ERROR_MESSAGE, ""))
			return
		}
	}
	if reqInfo.Account == adminObj.Account && utils.Md5(reqInfo.Password) == adminObj.Password {
		token := GetFromDatabase(constant.DEFAULT_TOKEN_DB_KEY)
		if token == "" {
			token = utils.Md5(fmt.Sprintf("%v%s", time.Now().UnixNano(), reqInfo.Password))
			SetWiteTTLFromDatabase(constant.DEFAULT_TOKEN_DB_KEY, token, constant.DEFAULT_TOKEN_DB_KEY_EXPIRE)
		}
		c.JSON(http.StatusOK, getResponse(constant.HTTP_OK, "", json_struct.UserToken{Token: token}))
		return
	}
	c.JSON(http.StatusOK, getResponse(constant.PASSWORD_ERROR, constant.PASSWORD_ERROR_MESSAGE, ""))
}

// 用户信息获取
func getInfo(c *gin.Context) {
	var adminObj json_struct.AdminUser
	adminJson := GetFromDatabase(constant.DEFAULT_ACCOUNT_DB_KEY)
	_ = json.Unmarshal([]byte(adminJson), &adminObj)
	adminObj.Password = ""
	c.JSON(http.StatusOK, getResponse(constant.HTTP_OK, "", adminObj))
}

// 用户注销
func userLogout(c *gin.Context) {
	SetFromDatabase(constant.DEFAULT_TOKEN_DB_KEY, "")
	c.JSON(http.StatusOK, getResponse(constant.HTTP_OK, "", ""))
}

// 服务端消耗资源监听
func serversResourcesListen(c *gin.Context) {
	serverId, ok := c.Params.Get("serverId")
	if serverId == "" || !ok {

	}
	upGrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	// 校验token
	for {
		_, tokenByte, err := ws.ReadMessage()
		if err != nil {
			ws.Close()
			return
		}
		dbtoken := GetFromDatabase(constant.DEFAULT_TOKEN_DB_KEY)
		if string(tokenByte) == dbtoken {
			break
		} else {
			ws.Close()
			return
		}
	}
	AppendResourceWsToPool(c, serverId, ws)
}

// 实时获取服务器输出以及输入命令
func serversStdListen(c *gin.Context) {
	serverId, ok := c.Params.Get("serverId")
	if serverId == "" || !ok {

	}
	upGrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	// 校验token
	for {
		_, tokenByte, err := ws.ReadMessage()
		if err != nil {
			ws.Close()
			return
		}
		dbtoken := GetFromDatabase(constant.DEFAULT_TOKEN_DB_KEY)
		if string(tokenByte) == dbtoken {
			break
		} else {
			ws.Close()
			return
		}
	}
	AppendStdWsToPool(serverId, ws)
	ListenStdinFromWs(serverId, ws)
}

// 修改用户信息
func updateUserData(c *gin.Context) {
	// 获取用户设置信息
	var reqInfo json_struct.AdminUser
	if err := c.BindJSON(&reqInfo); err != nil {
		WriteLogToDefault(errorFormat(err), constant.LOG_ERROR)
		c.JSON(http.StatusOK, getResponse(constant.HTTP_PARAMS_ERROR, constant.HTTP_PARAMS_ERROR_MESSAGE, ""))
		return
	}
	// 获取数据库信息
	var adminObj json_struct.AdminUser
	adminJson := GetFromDatabase(constant.DEFAULT_ACCOUNT_DB_KEY)
	if err := json.Unmarshal([]byte(adminJson), &adminObj); err != nil {
		WriteLogToDefault(errorFormat(err), constant.LOG_ERROR)
		c.JSON(http.StatusOK, getResponse(constant.HTTP_SYSTEM_ERROR, constant.HTTP_SYSTEM_ERROR_MESSAGE, ""))
		return
	}

	if reqInfo.Account != "" {
		adminObj.Account = reqInfo.Account
	}
	if reqInfo.Avatar != "" {
		adminObj.Avatar = reqInfo.Avatar
	}
	if reqInfo.Password != "" {
		reqInfo.Password = utils.Md5(reqInfo.Password)
	}
	jsonStr, _ := json.Marshal(reqInfo)
	SetFromDatabase(constant.DEFAULT_ACCOUNT_DB_KEY, string(jsonStr))
	c.JSON(http.StatusOK, getResponse(constant.HTTP_OK, "", ""))
}

// 获取配置内容
func getConfig(c *gin.Context) {
	// 获取用户设置信息
	jsonObj := make([]*json_struct.ConfParam, 0)
	conf := GetConfInstance().GetConfigObj()
	for _, v := range conf {
		jsonObj = append(jsonObj, v)
	}
	c.JSON(http.StatusOK, getResponse(constant.HTTP_OK, "", jsonObj))
}

// 更新配置内容
func updateConfig(c *gin.Context) {
	// 获取用户设置信息
	var reqInfo []*json_struct.ConfParam
	if err := c.BindJSON(&reqInfo); err != nil {
		WriteLogToDefault(errorFormat(err), constant.LOG_ERROR)
		c.JSON(http.StatusOK, getResponse(constant.HTTP_PARAMS_ERROR, constant.HTTP_PARAMS_ERROR_MESSAGE, ""))
		return
	}
	confObj := GetConfInstance()
	for _, v := range reqInfo {
		confObj.SetConfig(v.Name, v.ConfVal)
	}
	// 执行配置更改回调
	RunChangeConfCallBacks()
	c.JSON(http.StatusOK, getResponse(constant.HTTP_OK, "", ""))
}

// 设置初始账号密码
func setDefaultAccount() *json_struct.AdminUser {
	pwd := utils.Md5(constant.DEFAULT_PASSWORD)
	adminObj := &json_struct.AdminUser{
		Nickname: constant.DEFAULT_ACCOUNT,
		Account:  constant.DEFAULT_ACCOUNT,
		Password: pwd,
		Roles:    nil,
	}
	adminJson, _ := json.Marshal(adminObj)
	SetFromDatabase(constant.DEFAULT_ACCOUNT_DB_KEY, string(adminJson))
	return adminObj
}

// 返回数据格式化
func getResponse(code int, message string, data interface{}) gin.H {
	responseData := make(gin.H)
	responseData["code"] = code
	responseData["message"] = message
	responseData["data"] = data
	return responseData
}

func errorFormat(err error) string {
	return fmt.Sprintf("web error: %v", err)
}
