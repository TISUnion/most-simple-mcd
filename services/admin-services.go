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

func (a *AdminService) GetLog(ctx context.Context, req *api.GetLogReq) (resp *api.GetLogResp, err error) {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		a._getLog(ginCtx)
	}
	return
}

func (a *AdminService) AddUpToContainer(ctx context.Context, req *api.AddUpToContainerReq) (resp *api.AddUpToContainerResp, err error) {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		err = a._addUpToContainer(ginCtx)
	}
	resp = new(api.AddUpToContainerResp)
	return
}

func (a *AdminService) GetConfig(ctx context.Context, req *api.GetConfigReq) (resp *api.GetConfigResp, err error) {
	confs := a._getConfig()
	list := make([]*api.GetConfigResp_Record, 0)
	for _, c := range confs {
		list = append(list, &api.GetConfigResp_Record{
			ConfVal:              c.ConfVal,
			Name:                 c.Name,
			Level:                c.Level,
			Description:          c.Description,
			IsAlterable:          c.IsAlterable,
			XXX_NoUnkeyedLiteral: struct{}{},
			XXX_unrecognized:     nil,
			XXX_sizecache:        0,
		})
	}
	resp = &api.GetConfigResp{
		List: list,
	}
	return
}

func (a *AdminService) UpdateConfig(ctx context.Context, req *api.UpdateConfigReq) (resp *api.UpdateConfigResp, err error) {
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

func (a *AdminService) OperatePlugin(ctx context.Context, req *api.OperatePluginReq) (resp *api.OperatePluginResp, err error) {
	resp = new(api.OperatePluginResp)
	err = a._operatePlugin(req.ServerId, req.PluginId, req.OperateType)
	return
}

func (a *AdminService) GetConfigVal(ctx context.Context, req *api.GetConfigValReq) (resp *api.GetConfigValResp, err error) {
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

func (a *AdminService) RunCommand(ctx context.Context, req *api.RunCommandReq) (resp *api.RunCommandResp, err error) {
	err = a._runCommand(req.Command, req.ServerId, req.Type)
	resp = new(api.RunCommandResp)
	return
}

func (a *AdminService) DelTmpFlie(ctx context.Context, req *api.DelTmpFlieReq) (resp *api.DelTmpFlieResp, err error) {
	err = a._delTmpFlie()
	resp = new(api.DelTmpFlieResp)
	return
}

func (a *AdminService) CloseMcd(ctx context.Context, req *api.CloseMcdReq) (resp *api.CloseMcdResp, err error) {
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
	header, err := c.FormFile(constant.UPLOAD_FILE_NAME)
	if err != nil {
		modules.WriteLogToDefault(constant.PARSE_FILE_ERROR+err.Error(), constant.LOG_ERROR)
		return errors.New(constant.HTTP_PARAMS_ERROR_MESSAGE)
	}

	dst := filepath.Join(modules.GetConfVal(constant.TMP_PATH), header.Filename)
	if err := c.SaveUploadedFile(header, dst); err != nil {
		modules.WriteLogToDefault(constant.COPY_FILE_ERROR+err.Error(), constant.LOG_ERROR)
		return errors.New(constant.HTTP_PARAMS_ERROR_MESSAGE)
	}
	ext := filepath.Ext(header.Filename)
	if ext != constant.JAR_SUF {
		return errors.New(constant.HTTP_PARAMS_ERROR_MESSAGE)
	}
	ctr := modules.GetMinecraftServerContainerInstance()
	port, _ := strconv.ParseInt(c.DefaultPostForm(constant.UPLOAD_PORT_TEXT, strconv.Itoa(constant.MC_DEFAULT_PORT)), 10, 64)
	name := c.DefaultPostForm(constant.UPLOAD_NAME_TEXT, header.Filename)
	memory, _ := strconv.ParseInt(c.DefaultPostForm(constant.UPLOAD_MEMORY_TEXT, strconv.Itoa(constant.MC_DEFAULT_MEMORY)), 10, 64)
	if memory == 0 {
		memory = constant.MC_DEFAULT_MEMORY
	}
	mcCfg := ctr.HandleMcFile(dst, name, port, memory)
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
	utils.CreatDir(modules.GetConfVal(constant.TMP_PATH))
	if err != nil {
		modules.WriteLogToDefault(err.Error(), constant.LOG_ERROR)
		return errors.New(constant.HTTP_PARAMS_ERROR_MESSAGE)
	}
	return nil
}

func (a *AdminService) _closeMcd() {
	modules.SendExitSign()
}
