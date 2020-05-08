package services

import (
	"github.com/TISUnion/most-simple-mcd/grpc/api"
	"github.com/TISUnion/most-simple-mcd/modules"
)

func RegisterServices() {
	ginEngine := modules.GetGinServerInstanceRouter()
	api.RegisterServerMcServerGinServer(ginEngine, &ServerService{})
}
