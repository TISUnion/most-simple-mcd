// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package wire

import (
	_interface "github.com/TISUnion/most-simple-mcd/interface"
	"github.com/TISUnion/most-simple-mcd/interface/container"
	"github.com/TISUnion/most-simple-mcd/service"
	"github.com/google/wire"
)

func GetConfIntance() _interface.Conf {
	wire.Build(service.InitFlag, service.GetConfObj)
	return &service.Conf{}
}

func GetJobContainerIntance() container.JobContainer {
	wire.Build(service.GetJobContainerObj)
	return &service.JobContainer{}
}
