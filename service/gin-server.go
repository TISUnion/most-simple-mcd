package service

import (
	"fmt"
	"github.com/TISUnion/most-simple-mcd/constant"
	"github.com/TISUnion/most-simple-mcd/interface/server"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var ginServerInstance *GinServer

type GinServer struct {
	router     *gin.Engine
	httpServer *http.Server
	port       int
	lock       *sync.Mutex
}

func (g *GinServer) GetRouter() *gin.Engine {
	return g.router
}

func (g *GinServer) ChangeConfCallBack() {

}

func (g *GinServer) DestructCallBack() {

}

func (g *GinServer) InitCallBack() {
	RegisterRouter()
}

func (g *GinServer) Start() error {
	g.lock.Lock()
	defer g.lock.Unlock()
	go g.httpServer.ListenAndServe()
	return nil
}

func (g *GinServer) Stop() error {
	g.lock.Lock()
	defer g.lock.Unlock()
	_ = g.httpServer.Close()
	g.httpServer = getHttpServerObj(g.port, g.router)
	return nil
}

func (g *GinServer) Restart() error {
	g.lock.Lock()
	defer g.lock.Unlock()
	_ = g.httpServer.Close()
	g.httpServer = getHttpServerObj(g.port, g.router)
	go g.httpServer.ListenAndServe()
	return nil
}

func GetGinServerInstance() server.GinServer {
	if ginServerInstance != nil {
		return ginServerInstance
	}
	portStr := GetConfVal(constant.MANAGE_HTTP_SERVER_PORT)
	port, _ := strconv.Atoi(portStr)

	// 添加日志
	logger := GetLogContainerInstance().AddLog("gin-server")
	gin.DefaultWriter = logger
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	httpServer := getHttpServerObj(port, router)

	ginServerInstance = &GinServer{
		router:     router,
		httpServer: httpServer,
		port:       port,
		lock:       &sync.Mutex{},
	}
	RegisterCallBack(ginServerInstance)
	return ginServerInstance
}

func GetGinServerInstanceRouter() *gin.Engine {
	return GetGinServerInstance().GetRouter()
}

func getHttpServerObj(port int, router *gin.Engine) *http.Server {
	return &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           router,
		ReadTimeout:       time.Second * 8,
		ReadHeaderTimeout: time.Second * 8,
		WriteTimeout:      time.Second * 8,
	}
}
