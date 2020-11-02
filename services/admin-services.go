package services

import (
	"context"
	"errors"
	"github.com/TISUnion/most-simple-mcd/constant"
	"github.com/TISUnion/most-simple-mcd/grpc/api"
	"github.com/TISUnion/most-simple-mcd/models"
	"github.com/TISUnion/most-simple-mcd/modules"
	"github.com/TISUnion/most-simple-mcd/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type AdminService struct {
}

func (a *AdminService) UpMapToMcServer(ctx context.Context, _ *api.AddUpToContainerReq) (resp *api.AddUpToContainerResp, err error) {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		err = a._upMapToMcServer(ginCtx)
	}
	resp = new(api.AddUpToContainerResp)
	return
}

func (a *AdminService) GetLog(ctx context.Context, _ *api.GetLogReq) (resp *api.GetLogResp, err error) {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		a._getLog(ginCtx)
	}
	return
}

func (a *AdminService) AddUpToContainer(ctx context.Context, _ *api.AddUpToContainerReq) (resp *api.AddUpToContainerResp, err error) {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		err = a._addUpToContainer(ginCtx)
	}
	resp = new(api.AddUpToContainerResp)
	return
}

func (a *AdminService) GetConfig(context.Context, *api.GetConfigReq) (resp *api.GetConfigResp, err error) {
	confs := a._getConfig()
	list := make([]*api.GetConfigResp_Record, 0)
	for _, c := range confs {
		list = append(list, &api.GetConfigResp_Record{
			ConfVal:     c.ConfVal,
			Name:        c.Name,
			Level:       c.Level,
			Description: c.Description,
			IsAlterable: c.IsAlterable,
		})
	}
	resp = &api.GetConfigResp{
		List: list,
	}
	return
}

func (a *AdminService) UpdateConfig(_ context.Context, req *api.UpdateConfigReq) (resp *api.UpdateConfigResp, err error) {
	confs := make([]*models.ConfParam, 0)
	for _, reqC := range req.List {
		confs = append(confs, &models.ConfParam{
			ConfVal:     reqC.ConfVal,
			Name:        reqC.Name,
			Level:       reqC.Level,
			Description: reqC.Description,
			IsAlterable: reqC.IsAlterable,
		})
	}
	err = a._updateConfig(confs)
	resp = new(api.UpdateConfigResp)
	return
}

func (a *AdminService) OperatePlugin(_ context.Context, req *api.OperatePluginReq) (resp *api.OperatePluginResp, err error) {
	resp = new(api.OperatePluginResp)
	err = a._operatePlugin(req.ServerId, req.PluginId, req.OperateType)
	return
}

func (a *AdminService) GetConfigVal(_ context.Context, req *api.GetConfigValReq) (resp *api.GetConfigValResp, err error) {
	conf, e := a._getConfigVal(req.Name)
	err = e
	if e == nil {
		resp = &api.GetConfigValResp{
			ConfVal:     conf.ConfVal,
			Name:        conf.Name,
			Level:       conf.Level,
			Description: conf.Description,
			IsAlterable: conf.IsAlterable,
		}
	}
	return
}

func (a *AdminService) RunCommand(_ context.Context, req *api.RunCommandReq) (resp *api.RunCommandResp, err error) {
	err = a._runCommand(req.Command, req.ServerId, req.Type)
	resp = new(api.RunCommandResp)
	return
}

func (a *AdminService) DelTmpFlie(context.Context, *api.DelTmpFlieReq) (resp *api.DelTmpFlieResp, err error) {
	err = a._delTmpFlie()
	resp = new(api.DelTmpFlieResp)
	return
}

func (a *AdminService) CloseMcd(context.Context, *api.CloseMcdReq) (resp *api.CloseMcdResp, err error) {
	resp = new(api.CloseMcdResp)
	a._closeMcd()
	return
}

func (a *AdminService) _getLog(ginCtx *gin.Context) {
	//type 1. 根据id获取服务端日志  2. gin服务器日志    3. 默认全局日志
	type_ := ginCtx.Query(constant.QUERY_TYPE)
	//id 如果是根据id获取服务端日志
	id := ginCtx.Query(constant.QUERY_ID)
	if type_ == "" {
		ginCtx.JSON(http.StatusOK, getResponse(constant.HTTP_PARAMS_ERROR, constant.HTTP_PARAMS_ERROR_MESSAGE, ""))
	}
	originFilePath := ""
	logName := time.Now().Format(constant.LOG_TIME_FORMAT) + ".zip"
	filePath := filepath.Join(modules.GetConfVal(constant.TMP_PATH), logName)
	switch type_ {
	case constant.LOG_TYPE_SERVER:
		ctr := modules.GetMinecraftServerContainerInstance()
		serv, err := ctr.GetServerById(id)
		if err != nil {
			ginCtx.JSON(http.StatusOK, getResponse(constant.HTTP_PARAMS_ERROR, constant.HTTP_PARAMS_ERROR_MESSAGE, ""))
			return
		}
		runPath := filepath.Dir(serv.GetServerConf().RunPath)
		originFilePath = filepath.Join(runPath, constant.LOG_DIR)
	case constant.LOG_TYPE_GIN:
		originFilePath = filepath.Join(modules.GetConfVal(constant.WORKSPACE), constant.LOG_DIR, constant.GIN_LOG_NAME)
	case constant.LOG_TYPE_DEFAULT:
		originFilePath = filepath.Join(modules.GetConfVal(constant.WORKSPACE), constant.LOG_DIR, constant.DEFAULT_LOG_NAME)
	}
	// 压缩
	_ = utils.CompressFile(originFilePath, filePath)
	ginCtx.FileAttachment(filePath, logName)
}

func (a *AdminService) _addUpToContainer(c *gin.Context) error {
	filename, dst, err := a.getUploadFile(c)
	if err != nil {
		return err
	}
	ext := filepath.Ext(filename)
	if ext != constant.JAR_SUF {
		return errors.New(constant.HTTP_PARAMS_ERROR_MESSAGE)
	}
	ctr := modules.GetMinecraftServerContainerInstance()
	port, _ := strconv.ParseInt(c.DefaultPostForm(constant.UPLOAD_PORT_TEXT, strconv.Itoa(constant.MC_DEFAULT_PORT)), 10, 64)
	name := c.DefaultPostForm(constant.UPLOAD_NAME_TEXT, filename)
	memory, _ := strconv.ParseInt(c.DefaultPostForm(constant.UPLOAD_MEMORY_TEXT, strconv.Itoa(constant.MC_DEFAULT_MEMORY)), 10, 64)
	if memory == 0 {
		memory = constant.MC_DEFAULT_MEMORY
	}
	side := c.DefaultPostForm(constant.UPLOAD_SIDE_TEXT, constant.VANILLA_SERVER)
	comment := c.DefaultPostForm(constant.UPLOAD_COMMENT_TEXT, "")
	mcCfg := ctr.HandleMcFile(dst, name, port, memory, side, comment)
	ctr.AddServer(mcCfg, true)
	return nil
}

func (a *AdminService) _getConfig() []*models.ConfParam {
	// 获取用户设置信息
	jsonObj := make([]*models.ConfParam, 0)
	conf := modules.GetConfInstance().GetConfigObj()
	for _, v := range conf {
		jsonObj = append(jsonObj, v)
	}
	return jsonObj
}

func (a *AdminService) _updateConfig(updateConf []*models.ConfParam) error {
	confObj := modules.GetConfInstance()
	for _, v := range updateConf {
		confObj.SetConfParam(v.Name, v.ConfVal, v.Description, constant.CONF_SYSTEM_LEVEL)
	}
	// 执行配置更改回调
	modules.RunChangeConfCallBacks()

	return nil
}

func (a *AdminService) _operatePlugin(serverId string, pluginId []string, opType int64) error {
	ctr := modules.GetMinecraftServerContainerInstance()
	serv, err := ctr.GetServerById(serverId)
	if err != nil {
		modules.WriteLogToDefault(errorFormat(err), constant.LOG_ERROR)
		return err
	}
	for _, plId := range pluginId {
		switch opType {
		case constant.MC_PLUGIN_START:
			serv.UnbanPlugin(plId)
		case constant.MC_PLUGIN_STOP:
			serv.BanPlugin(plId)
		}
	}
	return nil
}

func (a *AdminService) _getConfigVal(configName string) (*models.ConfParam, error) {
	conf := modules.GetConfInstance().GetConfigObj()
	valconf, ok := conf[configName]
	if !ok {
		return nil, errors.New(constant.HTTP_PARAMS_ERROR_MESSAGE)
	}
	return valconf, nil
}

func (a *AdminService) _runCommand(command, serverId string, type_ int64) error {
	if ok := modules.RunOneCommand(serverId, command, type_); !ok {
		return errors.New(constant.HTTP_PARAMS_ERROR_MESSAGE)
	}
	return nil
}

func (a *AdminService) _delTmpFlie() error {
	err := os.RemoveAll(modules.GetConfVal(constant.TMP_PATH))
	// 创建tmp目录
	_ = utils.CreatDir(modules.GetConfVal(constant.TMP_PATH))
	if err != nil {
		modules.WriteLogToDefault(err.Error(), constant.LOG_ERROR)
		return errors.New(constant.HTTP_PARAMS_ERROR_MESSAGE)
	}
	return nil
}

func (a *AdminService) _closeMcd() {
	modules.SendExitSign()
}

func (a *AdminService) _upMapToMcServer(c *gin.Context) error {
	filename, dst, err := a.getUploadFile(c)
	if err != nil {
		return err
	}
	ext := filepath.Ext(filename)
	if ext != constant.ZIP_SUF {
		return errors.New(constant.HTTP_PARAMS_ERROR_MESSAGE)
	}
	ctr := modules.GetMinecraftServerContainerInstance()
	id := c.DefaultPostForm(constant.UPLOAD_ID_TEXT, "")
	if id == "" {
		return errors.New(constant.HTTP_PARAMS_ERROR_MESSAGE)
	}
	srv, err := ctr.GetServerById(id)
	if err != nil {
		return errors.New(constant.HTTP_PARAMS_ERROR_MESSAGE)
	}
	// 锁住服务端，防止后续操作时文件被锁住
	srv.LockWithMessage("导入地图数据中")
	// 保证导入时服务端为关闭状态
	if srv.GetServerConf().State != constant.MC_SERVER_STOP {
		srv.Unlock()
		return errors.New(constant.HTTP_PARAMS_ERROR_MESSAGE)
	}
	go func() {
		defer srv.Unlock()
		uncompressFilepath := filepath.Join(constant.TMP_PATH, utils.GetRandomString(10))
		err = utils.UnCompressDir(dst, uncompressFilepath)
		if err != nil {
			modules.WriteLogToDefault()
			return
		}
		serverJarArr, err := filepath.Glob(uncompressFilepath + "/*.jar")
		if err != nil || len(serverJarArr) > 1 {
			return
		}
		// 如果上传文件提供了服务端，则使用上传服务端
		if len(serverJarArr) == 1 {
			serverJar := serverJarArr[0]
			err := os.Rename(serverJar, filepath.Join(uncompressFilepath, srv.GetServerEntryId()+".jar"))
			if err != nil {
				return
			}
		}
		// 删除日志
		err = os.Remove(filepath.Join(uncompressFilepath, constant.LOG_DIR))
		if err != nil {
			return
		}
		serverBasePath := filepath.Join(modules.GetConfVal(constant.WORKSPACE), constant.MC_SERVER_DIR, srv.GetServerEntryId())

		// 循环转移
		err = filepath.Walk(uncompressFilepath, func(path string, info os.FileInfo, err error) error {
			if utils.IsFile(path) {
				filePath := filepath.Join(serverBasePath, info.Name())
				_ = os.Rename(path, filePath)
			}
			return nil
		})
	}()
	return nil
}

func (a *AdminService) getUploadFile(c *gin.Context) (filename, dst string, err error) {
	header, err := c.FormFile(constant.UPLOAD_FILE_NAME)
	if err != nil {
		modules.WriteLogToDefault(constant.PARSE_FILE_ERROR+err.Error(), constant.LOG_ERROR)
		err = errors.New(constant.HTTP_PARAMS_ERROR_MESSAGE)
		return
	}
	dst = filepath.Join(modules.GetConfVal(constant.TMP_PATH), header.Filename)
	if e := c.SaveUploadedFile(header, dst); e != nil {
		modules.WriteLogToDefault(constant.COPY_FILE_ERROR+e.Error(), constant.LOG_ERROR)
		err = errors.New(constant.HTTP_PARAMS_ERROR_MESSAGE)
		return
	}
	filename = header.Filename
	return
}
