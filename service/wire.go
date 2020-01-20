// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package service

import (
	_interface "github.com/TISUnion/most-simple-mcd/interface"
	"github.com/TISUnion/most-simple-mcd/interface/container"
	"github.com/google/wire"
)

func GetConfInstance() _interface.Conf {
	wire.Build(InitFlag, GetConfObj)
	return &Conf{}
}

func GetJobContainerInstance() container.JobContainer {
	wire.Build(GetJobContainerObj)
	return &JobContainer{}
}

func GetLogContainerInstance() container.LogContainer {
	wire.Build(GetLogContainerObj)
	return &LogContainer{}
}
