// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package wire

import (
	_interface "github.com/TISUnion/most-simple-mcd/interface"
	"github.com/TISUnion/most-simple-mcd/interface/container"
	"github.com/TISUnion/most-simple-mcd/service"
	container2 "github.com/TISUnion/most-simple-mcd/service/container"
	"github.com/google/wire"
)

func GetConfInstance() _interface.Conf {
	wire.Build(service.InitFlag, service.GetConfObj)
	return &service.Conf{}
}

func GetJobContainerInstance() container.JobContainer {
	wire.Build(container2.GetJobContainerObj)
	return &container2.JobContainer{}
}

func GetLogContainerInstance() container.LogContainer {
	wire.Build(GetConfInstance, GetJobContainerInstance, container2.GetLogContainerObj)
	return &container2.LogContainer{}
}
