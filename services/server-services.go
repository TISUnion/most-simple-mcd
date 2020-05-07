package services

import (
	"context"
	"github.com/TISUnion/most-simple-mcd/grpc/api"
	"github.com/TISUnion/most-simple-mcd/models"
	"github.com/TISUnion/most-simple-mcd/modules"
)

type ServerService struct {
}

func (s *ServerService) ListenResource(ctx context.Context, req *api.ListenResourceReq) (resp *api.ListenResourceResp, err error) {
	panic("implement me")
}

func (s *ServerService) ServerInteraction(ctx context.Context, req *api.ServerInteractionReq) (resp *api.ServerInteractionResp, err error) {
	panic("implement me")
}

func (s *ServerService) List(ctx context.Context, req *api.ListReq) (resp *api.ListResp, err error) {
	return &api.ListResp{
		List: models.ServerConfObjs2ServerConfPbs(s._list()),
	}, nil
}

func (s *ServerService) GetServerState(ctx context.Context, req *api.GetServerStateReq) (resp *api.GetServerStateReq, err error) {
	panic("implement me")
}

func (s *ServerService) Detail(ctx context.Context, req *api.DetailReq) (resp *api.DetailResp, err error) {
	panic("implement me")
}

func (s *ServerService) OperateServer(ctx context.Context, req *api.OperateServerReq) (resp *api.OperateServerResp, err error) {
	panic("implement me")
}

func (s *ServerService) UpdateServerInfo(ctx context.Context, req *api.UpdateServerInfoReq) (resp *api.UpdateServerInfoResp, err error) {
	panic("implement me")
}

func (s *ServerService) _list() []*models.ServerConf {
	ctr := modules.GetMinecraftServerContainerInstance()
	return ctr.GetAllServerConf()
}
