// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package wire

import (
	"github.com/TISUnion/most-simple-mcd/interface"
	"github.com/TISUnion/most-simple-mcd/service"
)

// Injectors from wire.go:

func GetConfIntance() _interface.Conf {
	terminalType := service.InitFlag()
	conf := service.GetConfObj(terminalType)
	return conf
}
