package server

import "github.com/gin-gonic/gin"

type GinServer interface {
	BasicServer
	GetRouter() *gin.Engine
}
