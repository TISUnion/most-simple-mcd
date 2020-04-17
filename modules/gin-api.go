package modules

import (
	"encoding/json"
	"fmt"
	"github.com/TISUnion/most-simple-mcd/constant"
	json_struct "github.com/TISUnion/most-simple-mcd/json-struct"
	"github.com/TISUnion/most-simple-mcd/utils"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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
	// 静态后台管理前端文件
	router.Use(static.Serve("/", static.LocalFile(filepath.Join(GetConfVal(constant.WORKSPACE), constant.Web_FILE_DIR_NAME), true)))
	// 静态日志文件
	router.Use(static.Serve("/static/file/logs", static.LocalFile(filepath.Join(GetConfVal(constant.WORKSPACE), constant.LOG_DIR), true)))

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
		v1.POST("/user/account", updateUserData)
		// 获取配置
		v1.GET("/config/list", getConfig)
		// 修改配置
		v1.POST("/config", updateConfig)
		// 获取服务端列表
		v1.GET("/server/list", serversInfoList)
		// 获取服务端运行状态
		v1.GET("/server/ping", getServerState)
		// 获取服务端详情
		v1.GET("/server/detail", serverDetail)
		// 操作服务端
		v1.POST("/server", operateServer)
		// 操作插件
		v1.POST("/plugin", operatePlugin)
		// 获取一个配置
		v1.GET("/config/val", getConfigVal)

		// 修改服务端参数
		v1.POST("/server/info", updateServerInfo)
		// 在指定服务端中运行一条命令
		v1.POST("/server/run/command", runCommand)
		// 获取日志
		v1.GET("/log/download", getLog)
		// 删除临时文件
		v1.POST("/tmp/files", delTmpFlie)
		// 获取上传服务端文件，并注入到容器中
		v1.POST("/upload/server", addUpToContainer)

		// 关闭mcd
		v1.POST("/close", func(c *gin.Context) {
			SendExitSign()
		})
	}
	// websocket实时监听服务端耗费资源
	router.GET("/server/resources/listen", serversResourcesListen)
	// websocket实时获取服务器输出
	router.GET("/server/std/listen", serversStdListen)
}

func addUpToContainer(c *gin.Context) {
	header, err := c.FormFile(constant.UPLOAD_FILE_NAME)
	if err != nil {
		c.JSON(http.StatusOK, getResponse(constant.HTTP_PARAMS_ERROR, constant.HTTP_PARAMS_ERROR_MESSAGE, ""))
		WriteLogToDefault(constant.PARSE_FILE_ERROR+err.Error(), constant.LOG_ERROR)
		return
	}

	dst := filepath.Join(GetConfVal(constant.TMP_PATH), header.Filename)
	if err := c.SaveUploadedFile(header, dst); err != nil {
		c.JSON(http.StatusOK, getResponse(constant.HTTP_PARAMS_ERROR, constant.HTTP_PARAMS_ERROR_MESSAGE, ""))
		WriteLogToDefault(constant.COPY_FILE_ERROR+err.Error(), constant.LOG_ERROR)
		return
	}
	ext := filepath.Ext(header.Filename)
	if ext != constant.JAR_SUF {
		c.JSON(http.StatusOK, getResponse(constant.HTTP_PARAMS_ERROR, constant.HTTP_PARAMS_ERROR_MESSAGE, ""))
		return
	}
	ctr := GetMinecraftServerContainerInstance()
	port, _ := strconv.Atoi(c.DefaultPostForm(constant.UPLOAD_PORT_TEXT, strconv.Itoa(constant.MC_DEFAULT_PORT)))
	name := c.DefaultPostForm(constant.UPLOAD_NAME_TEXT, header.Filename)
	memory, _ := strconv.Atoi(c.DefaultPostForm(constant.UPLOAD_MEMORY_TEXT, strconv.Itoa(constant.MC_DEFAULT_MEMORY)))
	if memory == 0 {
		memory = constant.MC_DEFAULT_MEMORY
	}
	mcCfg := ctr.HandleMcFile(dst, name, port, memory)
	ctr.AddServer(mcCfg, true)
	c.JSON(http.StatusOK, getResponse(constant.HTTP_OK, "", ""))
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
	serverId := c.Query(constant.QUERY_ID)
	if serverId == "" {

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
	serverId := c.Query(constant.QUERY_ID)
	if serverId == "" {

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

// 获取服务端列表
func serversInfoList(c *gin.Context) {
	ctr := GetMinecraftServerContainerInstance()
	confs := ctr.GetAllServerConf()
	c.JSON(http.StatusOK, getResponse(constant.HTTP_OK, "", confs))
}

// 获取服务端运行状态
func getServerState(c *gin.Context) {
	serverId := c.Query(constant.QUERY_ID)
	if serverId == "" {
		c.JSON(http.StatusOK, getResponse(constant.HTTP_PARAMS_ERROR, constant.HTTP_PARAMS_ERROR_MESSAGE, ""))
		return
	}
	ctr := GetMinecraftServerContainerInstance()
	serv, err := ctr.GetServerById(serverId)
	if err != nil {
		c.JSON(http.StatusOK, getResponse(constant.HTTP_PARAMS_ERROR, constant.HTTP_PARAMS_ERROR_MESSAGE, ""))
		return
	}
	servCfg := serv.GetServerConf()
	c.JSON(http.StatusOK, getResponse(constant.HTTP_OK, "", &json_struct.ServerRunState{State: servCfg.State}))
}

// 获取服务端详情
func serverDetail(c *gin.Context) {
	serverId := c.Query(constant.QUERY_ID)
	if serverId == "" {
		c.JSON(http.StatusOK, getResponse(constant.HTTP_PARAMS_ERROR, constant.HTTP_PARAMS_ERROR_MESSAGE, ""))
		return
	}
	ctr := GetMinecraftServerContainerInstance()
	serv, err := ctr.GetServerById(serverId)
	if err != nil {
		c.JSON(http.StatusOK, getResponse(constant.HTTP_PARAMS_ERROR, constant.HTTP_PARAMS_ERROR_MESSAGE, ""))
		return
	}
	info := &json_struct.ServerDetail{
		ServInfo: serv.GetServerConf(),
		PlgnInfo: serv.GetPluginsInfo(),
	}
	c.JSON(http.StatusOK, getResponse(constant.HTTP_OK, "", info))
}

// 修改服务端信息
func updateServerInfo(c *gin.Context) {
	var reqInfo json_struct.ServerConf
	if err := c.BindJSON(&reqInfo); err != nil {
		WriteLogToDefault(errorFormat(err), constant.LOG_ERROR)
		c.JSON(http.StatusOK, getResponse(constant.HTTP_PARAMS_ERROR, constant.HTTP_PARAMS_ERROR_MESSAGE, ""))
		return
	}
	ctr := GetMinecraftServerContainerInstance()
	serv, err := ctr.GetServerById(reqInfo.EntryId)
	if err != nil {
		c.JSON(http.StatusOK, getResponse(constant.HTTP_PARAMS_ERROR, constant.HTTP_PARAMS_ERROR_MESSAGE, ""))
		return
	}
	servConf := serv.GetServerConf()
	if reqInfo.Version != "" {
		servConf.Version = reqInfo.Version
	}
	if reqInfo.Name != "" {
		servConf.Name = reqInfo.Name
	}
	if reqInfo.Port != 0 {
		servConf.Port = reqInfo.Port
	}
	if reqInfo.Memory != 0 {
		servConf.Memory = reqInfo.Memory
	}
	if len(reqInfo.CmdStr) > 2 {
		servConf.CmdStr = reqInfo.CmdStr
	}
	if reqInfo.GameType != "" {
		servConf.GameType = reqInfo.GameType
	}
	serv.SetServerConf(servConf)
	ctr.SaveToDb()
	c.JSON(http.StatusOK, getResponse(constant.HTTP_OK, "", ""))
}

// 向服务端执行一条命令
func runCommand(c *gin.Context) {
	var reqInfo json_struct.SingleCommand
	if err := c.BindJSON(&reqInfo); err != nil {
		WriteLogToDefault(errorFormat(err), constant.LOG_ERROR)
		c.JSON(http.StatusOK, getResponse(constant.HTTP_PARAMS_ERROR, constant.HTTP_PARAMS_ERROR_MESSAGE, ""))
		return
	}
	RunOneCommand(reqInfo.ServerId, reqInfo.Command, reqInfo.Type)
	c.JSON(http.StatusOK, getResponse(constant.HTTP_OK, "", ""))
}

// 获取日志
func getLog(c *gin.Context) {
	//type 1. 根据id获取服务端日志  2. gin服务器日志    3. 默认全局日志
	type_ := c.Query(constant.QUERY_TYPE)
	//id 如果是根据id获取服务端日志
	id := c.Query(constant.QUERY_ID)
	if type_ == "" {
		c.JSON(http.StatusOK, getResponse(constant.HTTP_PARAMS_ERROR, constant.HTTP_PARAMS_ERROR_MESSAGE, ""))
	}
	originFilePath := ""
	logName := time.Now().Format(constant.LOG_TIME_FORMAT) + ".zip"
	filePath := filepath.Join(GetConfVal(constant.TMP_PATH), logName)
	switch type_ {
	case constant.LOG_TYPE_SERVER:
		ctr := GetMinecraftServerContainerInstance()
		serv, err := ctr.GetServerById(id)
		if err != nil {
			c.JSON(http.StatusOK, getResponse(constant.HTTP_PARAMS_ERROR, constant.HTTP_PARAMS_ERROR_MESSAGE, ""))
			return
		}
		runPath := filepath.Dir(serv.GetServerConf().RunPath)
		originFilePath = filepath.Join(runPath, constant.LOG_DIR)
		// 压缩
		_ = utils.CompressFile(originFilePath, filePath)

	case constant.LOG_TYPE_GIN:
		originFilePath = filepath.Join(GetConfVal(constant.WORKSPACE), constant.LOG_DIR, constant.GIN_LOG_NAME)
		// 压缩
		_ = utils.CompressFile(originFilePath, filePath)
	case constant.LOG_TYPE_DEFAULT:
		originFilePath = filepath.Join(GetConfVal(constant.WORKSPACE), constant.LOG_DIR, constant.DEFAULT_LOG_NAME)
		// 压缩
		_ = utils.CompressFile(originFilePath, filePath)
	}
	c.FileAttachment(filePath, logName)
}

// 删除临时文件
func delTmpFlie(c *gin.Context) {
	err := os.RemoveAll(GetConfVal(constant.TMP_PATH))
	// 创建tmp目录
	utils.CreatDir(GetConfVal(constant.TMP_PATH))
	if err != nil {
		WriteLogToDefault(err.Error(), constant.LOG_ERROR)
		c.JSON(http.StatusOK, getResponse(constant.HTTP_SYSTEM_ERROR, constant.HTTP_SYSTEM_ERROR_MESSAGE, ""))
		return
	}
	c.JSON(http.StatusOK, getResponse(constant.HTTP_OK, "", ""))
}

// 服务端操作
func operateServer(c *gin.Context) {
	ops := &json_struct.OperateServer{}
	if err := c.BindJSON(ops); err != nil {
		WriteLogToDefault(errorFormat(err), constant.LOG_ERROR)
		c.JSON(http.StatusOK, getResponse(constant.HTTP_PARAMS_ERROR, constant.HTTP_PARAMS_ERROR_MESSAGE, ""))
		return
	}
	opType := ops.OperateType
	ctr := GetMinecraftServerContainerInstance()
	var err error
	for _, s := range ops.ServerId {
		switch opType {
		case constant.MC_SERVER_START:
			err = ctr.StartById(s)
		case constant.MC_SERVER_STOP:
			err = ctr.StopById(s)
		case constant.MC_SERVER_RESTART:
			err = ctr.RestartById(s)
		}
		if err != nil {
			WriteLogToDefault(errorFormat(err), constant.LOG_ERROR)
		}
	}
	c.JSON(http.StatusOK, getResponse(constant.HTTP_OK, "", ""))
}

// 服务端插件操作
func operatePlugin(c *gin.Context) {
	opp := &json_struct.OperatePlugin{}
	if err := c.BindJSON(opp); err != nil {
		WriteLogToDefault(errorFormat(err), constant.LOG_ERROR)
		c.JSON(http.StatusOK, getResponse(constant.HTTP_PARAMS_ERROR, constant.HTTP_PARAMS_ERROR_MESSAGE, ""))
		return
	}
	opType := opp.OperateType
	ctr := GetMinecraftServerContainerInstance()
	serv, err := ctr.GetServerById(opp.ServerId)
	if err != nil {
		WriteLogToDefault(errorFormat(err), constant.LOG_ERROR)
		c.JSON(http.StatusOK, getResponse(constant.HTTP_PARAMS_ERROR, constant.HTTP_PARAMS_ERROR_MESSAGE, ""))
		return
	}
	for _, plId := range opp.PluginId {
		switch opType {
		case constant.MC_PLUGIN_START:
			serv.UnbanPlugin(plId)
		case constant.MC_PLUGIN_STOP:
			serv.BanPlugin(plId)
		}
	}
	c.JSON(http.StatusOK, getResponse(constant.HTTP_OK, "", ""))
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
	if reqInfo.Nickname != "" {
		adminObj.Nickname = reqInfo.Nickname
	}
	if reqInfo.Password != "" {
		adminObj.Password = utils.Md5(reqInfo.Password)
	}
	jsonStr, _ := json.Marshal(adminObj)
	SetFromDatabase(constant.DEFAULT_ACCOUNT_DB_KEY, string(jsonStr))
	// 清空token
	SetFromDatabase(constant.DEFAULT_TOKEN_DB_KEY, "")
	c.JSON(http.StatusOK, getResponse(constant.HTTP_OK, "", ""))
}

// 获取一个配置
func getConfigVal(c *gin.Context) {
	configName := c.Query(constant.QUERY_NAME)
	conf := GetConfInstance().GetConfigObj()
	valconf, ok := conf[configName]
	if !ok {
		c.JSON(http.StatusOK, getResponse(constant.HTTP_PARAMS_ERROR, constant.HTTP_PARAMS_ERROR_MESSAGE, ""))
	}
	c.JSON(http.StatusOK, getResponse(constant.HTTP_OK, "", valconf))
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
		confObj.SetConfParam(v.Name, v.ConfVal, v.Description, constant.CONF_SYSTEM_LEVEL)
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
