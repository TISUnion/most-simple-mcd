package modules

import (
	"context"
	"fmt"
	"github.com/TISUnion/most-simple-mcd/constant"
	"github.com/TISUnion/most-simple-mcd/interface/server"
	json_struct "github.com/TISUnion/most-simple-mcd/json-struct"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
	wsPool     map[string][]*websocket.Conn
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
	g.wsPool = make(map[string][]*websocket.Conn)
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

func (g *GinServer) appendWsToPool(ctx context.Context, serverId string, ws *websocket.Conn) {
	g.lock.Lock()
	defer g.lock.Unlock()
	mcContainer := GetMinecraftServerContainerInstance()
	mcServ, ok := mcContainer.GetServerById(serverId)
	if !ok {
		return
	}
	if _, ok := g.wsPool[serverId]; !ok {
		childCtx, cancelFunc := context.WithCancel(ctx)
		go g.websocketBroadcast(childCtx, mcServ, cancelFunc)
	}
	g.wsPool[serverId] = append(g.wsPool[serverId], ws)
	return
}

func (g *GinServer) websocketBroadcast(ctx context.Context, serv server.MinecraftServer, cancelFunc context.CancelFunc) {
	serv.StartMonitorServer()
	resouceChan := serv.GetServerMonitor().GetMessageChan()
	serverId := serv.GetServerConf().EntryId
	var resourceMsg *json_struct.MonitorMessage
	for {
		select {
		case resourceMsg = <-resouceChan:
			if len(g.wsPool[serverId]) == 0 {
				cancelFunc()
			}
			for i, ws := range g.wsPool[serverId] {
				if err := ws.WriteJSON(resourceMsg); err != nil {
					// 删除无用ws
					g.wsPool[serverId] = append(g.wsPool[serverId][:i], g.wsPool[serverId][i+1:len(g.wsPool[serverId])]...)
					ws.Close()
				}
			}
		case <-ctx.Done():
			delete(g.wsPool, serverId)
			serv.StopMonitorServer()
			return
		}
	}
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

func AppendWsToPool(ctx context.Context, serverId string, ws *websocket.Conn) {
	ginServerInstance.appendWsToPool(ctx, serverId, ws)
}
