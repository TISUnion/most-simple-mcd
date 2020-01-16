// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package wire

import (
	"github.com/TISUnion/most-simple-mcd/service"
	"github.com/google/wire"
)

func GetConfIntance() *service.Conf {
	wire.Build(service.InitFlag, service.GetConfObj)
	return &service.Conf{}
}