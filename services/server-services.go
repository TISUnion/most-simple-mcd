package services

import (
	"context"
	"errors"
	"github.com/TISUnion/most-simple-mcd/constant"
	"github.com/TISUnion/most-simple-mcd/grpc/api"
	"github.com/TISUnion/most-simple-mcd/models"
	"github.com/TISUnion/most-simple-mcd/modules"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

type ServerService struct {
}

func (s *ServerService) GetServerSide(ctx context.Context, req *api.GetServerSideReq) (resp *api.GetServerSideResp, err error) {
	resp = &api.GetServerSideResp{
		ServerSides: modules.GetAllServerSide(),
	}
	return
}

func (s *ServerService) ListenResource(ctx context.Context, req *api.ListenResourceReq) (resp *api.ListenResourceResp, err error) {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		s._listenResource(ginCtx)
	}
	return
}

func (s *ServerService) ServerInteraction(ctx context.Context, req *api.ServerInteractionReq) (resp *api.ServerInteractionResp, err error) {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		s._serverInteraction(ginCtx)
	}
	return
}

func (s *ServerService) List(ctx context.Context, req *api.ListReq) (resp *api.ListResp, err error) {
	list, e := s._list()
	if e != nil {
		err = e
		return
	}
	res := make([]*api.ListResp_Record, 0)
	for _, c := range list {
		res = append(res, &api.ListResp_Record{
			EntryId:        c.EntryId,
			Name:           c.Name,
			CmdStr:         c.CmdStr,
			Port:           c.Port,
			RunPath:        c.RunPath,
			IsMirror:       c.IsMirror,
			IsStartMonitor: c.IsStartMonitor,
			Memory:         c.Memory,
			Version:        c.Version,
			GameType:       c.GameType,
			State:          c.State,
			Ips:            c.Ips,
			Side:           c.Side,
			Comment:        c.Comment,
		})
	}
	return &api.ListResp{
		List: res,
	}, nil
}

func (s *ServerService) GetServerState(ctx context.Context, req *api.GetServerStateReq) (resp *api.GetServerStateResp, err error) {
	ste, e := s._getServerState(req.Id)
	if e != nil {
		err = e
		return
	}
	resp = &api.GetServerStateResp{
		State: ste,
	}
	return
}

func (s *ServerService) Detail(ctx context.Context, req *api.DetailReq) (resp *api.DetailResp, err error) {
	sc, pls, e := s._detail(req.Id)
	if e != nil {
		err = e
		return
	}
	pluginResp := make([]*api.DetailResp_PluginRecord, 0)
	for _, p := range pls {
		pluginResp = append(pluginResp, &api.DetailResp_PluginRecord{
			Name:            p.Name,
			Id:              p.Id,
			IsBan:           p.IsBan,
			CommandName:     p.CommandName,
			Description:     p.Description,
			HelpDescription: p.HelpDescription,
		})
	}

	return &api.DetailResp{
		EntryId:        sc.EntryId,
		Name:           sc.Name,
		CmdStr:         sc.CmdStr,
		Port:           sc.Port,
		RunPath:        sc.RunPath,
		IsMirror:       sc.IsMirror,
		IsStartMonitor: sc.IsStartMonitor,
		Memory:         sc.Memory,
		Version:        sc.Version,
		GameType:       sc.GameType,
		State:          sc.State,
		Ips:            sc.Ips,
		Pluginfo:       pluginResp,
		Side:           sc.Side,
		Comment:        sc.Comment,
	}, nil
}

func (s *ServerService) OperateServer(ctx context.Context, req *api.OperateServerReq) (resp *api.OperateServerResp, err error) {
	s._operateServer(req.ServerId, req.OperateType)
	resp = &api.OperateServerResp{}
	return
}

func (s *ServerService) UpdateServerInfo(ctx context.Context, req *api.UpdateServerInfoReq) (resp *api.UpdateServerInfoResp, err error) {
	reqModel := &models.ServerConf{
		EntryId:        req.EntryId,
		Name:           req.Name,
		CmdStr:         req.CmdStr,
		Port:           req.Port,
		RunPath:        req.RunPath,
		IsMirror:       req.IsMirror,
		IsStartMonitor: req.IsStartMonitor,
		Memory:         req.Memory,
		Version:        req.Version,
		GameType:       req.GameType,
		State:          req.State,
		Ips:            req.Ips,
		Comment:        req.Comment,
	}
	resp = &api.UpdateServerInfoResp{}
	err = s._updateServerInfo(reqModel)
	return
}

func (s *ServerService) _listenResource(ginCtx *gin.Context) {
	serverId := ginCtx.Query(constant.QUERY_ID)
	if serverId == "" {

	}
	upGrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	ws, err := upGrader.Upgrade(ginCtx.Writer, ginCtx.Request, nil)
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
		dbtoken := modules.GetFromDatabase(constant.DEFAULT_TOKEN_DB_KEY)
		if string(tokenByte) == dbtoken {
			break
		} else {
			ws.Close()
			return
		}
	}
	modules.AppendResourceWsToPool(ginCtx, serverId, ws)
}

func (s *ServerService) _serverInteraction(c *gin.Context) {
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
		dbtoken := modules.GetFromDatabase(constant.DEFAULT_TOKEN_DB_KEY)
		if string(tokenByte) == dbtoken {
			break
		} else {
			ws.Close()
			return
		}
	}
	modules.AppendStdWsToPool(serverId, ws)
	modules.ListenStdinFromWs(serverId, ws)
}

func (s *ServerService) _list() ([]*models.ServerConf, error) {
	ctr := modules.GetMinecraftServerContainerInstance()
	return ctr.GetAllServerConf(), nil
}

func (s *ServerService) _getServerState(serverId string) (state int64, err error) {
	if serverId == "" {
		err = errors.New(constant.HTTP_PARAMS_ERROR_MESSAGE)
		return
	}
	ctr := modules.GetMinecraftServerContainerInstance()
	serv, e := ctr.GetServerById(serverId)
	if e != nil {
		err = errors.New(constant.HTTP_PARAMS_ERROR_MESSAGE)
		return
	}
	servCfg := serv.GetServerConf()
	return servCfg.State, nil
}

func (s *ServerService) _detail(serverId string) (*models.ServerConf, []*models.PluginInfo, error) {
	if serverId == "" {
		err := errors.New(constant.HTTP_PARAMS_ERROR_MESSAGE)
		return nil, nil, err
	}
	ctr := modules.GetMinecraftServerContainerInstance()
	serv, err := ctr.GetServerById(serverId)
	if err != nil {
		err := errors.New(constant.HTTP_PARAMS_ERROR_MESSAGE)
		return nil, nil, err
	}
	return serv.GetServerConf(), serv.GetPluginsInfo(), nil
}

func (s *ServerService) _operateServer(serverIds []string, opType int64) error {
	ctr := modules.GetMinecraftServerContainerInstance()
	var err error
	for _, s := range serverIds {
		switch opType {
		case constant.MC_SERVER_START:
			err = ctr.StartById(s)
		case constant.MC_SERVER_STOP:
			err = ctr.StopById(s)
		case constant.MC_SERVER_RESTART:
			err = ctr.RestartById(s)
		}
		if err != nil {
			modules.WriteLogToDefault(errorFormat(err), constant.LOG_ERROR)
		}
	}
	return nil
}

func (s *ServerService) _updateServerInfo(reqModel *models.ServerConf) (err error) {
	ctr := modules.GetMinecraftServerContainerInstance()
	serv, e := ctr.GetServerById(reqModel.EntryId)
	if e != nil {
		err = e
		return
	}
	servConf := serv.GetServerConf()
	if reqModel.Version != "" {
		servConf.Version = reqModel.Version
	}
	if reqModel.Name != "" {
		servConf.Name = reqModel.Name
	}
	if reqModel.Port != 0 {
		servConf.Port = reqModel.Port
	}
	if reqModel.Memory != 0 {
		servConf.Memory = reqModel.Memory
	}
	if len(reqModel.CmdStr) > 2 {
		servConf.CmdStr = reqModel.CmdStr
	}
	if reqModel.GameType != "" {
		servConf.GameType = reqModel.GameType
	}
	if reqModel.Side != "" {
		servConf.Side = reqModel.Side
	}
	if reqModel.Comment != "" {
		servConf.Comment = reqModel.Comment
	}
	serv.SetServerConf(servConf)
	ctr.SaveToDb()
	return nil
}
