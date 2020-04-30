package server

import json_struct "github.com/TISUnion/most-simple-mcd/models"

type MonitorServer interface {
	BasicServer

	GetMessageChan() chan *json_struct.MonitorMessage
}
