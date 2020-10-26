package server

import "github.com/TISUnion/most-simple-mcd/models"

type MonitorServer interface {
	BasicServer

	GetMessageChan() chan *models.MonitorMessage
}
