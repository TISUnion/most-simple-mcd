package services

import (
	"context"
	"github.com/TISUnion/most-simple-mcd/grpc/api"
)

type AdminService struct {
	
}

func (a *AdminService) GetConfig(ctx context.Context, req *api.GetConfigReq) (resp *api.GetConfigResp, err error) {
	panic("implement me")
}

func (a *AdminService) UpdateConfig(ctx context.Context, req *api.UpdateConfigReq) (resp *api.UpdateConfigResp, err error) {
	panic("implement me")
}

func (a *AdminService) OperatePlugin(ctx context.Context, req *api.OperatePluginReq) (resp *api.OperatePluginResp, err error) {
	panic("implement me")
}

func (a *AdminService) GetConfigVal(ctx context.Context, req *api.GetConfigValReq) (resp *api.GetConfigValResp, err error) {
	panic("implement me")
}

func (a *AdminService) RunCommand(ctx context.Context, req *api.RunCommandReq) (resp *api.RunCommandResp, err error) {
	panic("implement me")
}

func (a *AdminService) GetLog(ctx context.Context, req *api.GetLogReq) (resp *api.GetLogResp, err error) {
	panic("implement me")
}

func (a *AdminService) DelTmpFlie(ctx context.Context, req *api.DelTmpFlieReq) (resp *api.DelTmpFlieResp, err error) {
	panic("implement me")
}

func (a *AdminService) AddUpToContainer(ctx context.Context, req *api.AddUpToContainerReq) (resp *api.AddUpToContainerResp, err error) {
	panic("implement me")
}

func (a *AdminService) CloseMcd(ctx context.Context, req *api.CloseMcdReq) (resp *api.CloseMcdResp, err error) {
	panic("implement me")
}
 