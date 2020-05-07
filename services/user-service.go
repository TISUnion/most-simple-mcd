package services

import (
	"context"
	"github.com/TISUnion/most-simple-mcd/grpc/api"
)

type UserService struct {

}

func (u UserService) Login(ctx context.Context, req *api.LoginReq) (resp *api.LoginResp, err error) {
	panic("implement me")
}

func (u UserService) Logout(ctx context.Context, req *api.LogoutReq) (resp *api.LogoutResp, err error) {
	panic("implement me")
}

func (u UserService) Info(ctx context.Context, req *api.InfoReq) (resp *api.InfoResp, err error) {
	panic("implement me")
}

func (u UserService) Update(ctx context.Context, req *api.UpdateReq) (resp *api.UpdateResp, err error) {
	panic("implement me")
}
