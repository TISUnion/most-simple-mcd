// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package wire

import (
	_interface "github.com/TISUnion/most-simple-mcd/interface"
	"github.com/TISUnion/most-simple-mcd/interface/container"
	"github.com/TISUnion/most-simple-mcd/service"
	"github.com/TISUnion/most-simple-mcd/service/containers"
	"github.com/google/wire"
)

func GetConfInstance() _interface.Conf {
	wire.Build(service.InitFlag, service.GetConfObj)
	return &service.Conf{}
}

func GetJobContainerInstance() container.JobContainer {
	wire.Build(containers.GetJobContainerObj)
	return &containers.JobContainer{}
}

func GetLogContainerInstance() container.LogContainer {
	wire.Build(GetConfInstance, GetJobContainerInstance, containers.GetLogContainerObj)
	return &containers.LogContainer{}
}
