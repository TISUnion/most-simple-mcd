package server

type HttpServer interface {
	BasicServer
	AddRoute(string, func(map[string]string)) error
}
