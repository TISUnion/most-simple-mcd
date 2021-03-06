package modules

import (
	"context"
	"fmt"
	"github.com/TISUnion/most-simple-mcd/constant"
	"github.com/TISUnion/most-simple-mcd/interface/server"
	"github.com/TISUnion/most-simple-mcd/models"
	"github.com/gin-contrib/gzip"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

var ginServerInstance *GinServer

type GinServer struct {
	router              *gin.Engine
	httpServer          *http.Server
	port                int
	lock                *sync.Mutex
	resourceWsPool      map[string][]*websocket.Conn
	stdoutWsPool        map[string][]*websocket.Conn
	stdoutChans         map[string]chan *models.ReciveMessage
	lockeResourceWsPool *sync.Mutex
	lockeStdoutWsPool   *sync.Mutex
}

func (g *GinServer) GetRouter() *gin.Engine {
	return g.router
}

func (g *GinServer) ChangeConfCallBack() {

}

func (g *GinServer) DestructCallBack() {

}

func (g *GinServer) InitCallBack() {
	g.resourceWsPool = make(map[string][]*websocket.Conn)
	g.stdoutWsPool = make(map[string][]*websocket.Conn)
	g.stdoutChans = make(map[string]chan *models.ReciveMessage)

	//  启用gzip压缩
	g.router.Use(gzip.Gzip(gzip.DefaultCompression))

	// 静态后台管理前端文件
	g.router.Use(static.Serve("/", static.LocalFile(filepath.Join(GetConfVal(constant.WORKSPACE), constant.Web_FILE_DIR_NAME), true)))

	// 静态日志文件
	g.router.Use(static.Serve("/static/file/logs", static.LocalFile(filepath.Join(GetConfVal(constant.WORKSPACE), constant.LOG_DIR), true)))
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

func (g *GinServer) appendResourceWsToPool(ctx context.Context, serverId string, ws *websocket.Conn) {
	g.lock.Lock()
	defer g.lock.Unlock()
	mcContainer := GetMinecraftServerContainerInstance()
	mcServ, err := mcContainer.GetServerById(serverId)
	if err != nil {
		ws.Close()
		return
	}
	g.lockeResourceWsPool.Lock()
	defer g.lockeResourceWsPool.Unlock()
	if _, ok := g.resourceWsPool[serverId]; !ok {
		childCtx, cancelFunc := context.WithCancel(ctx)
		go g.resourceWebsocketBroadcast(childCtx, mcServ, cancelFunc)
	}
	g.resourceWsPool[serverId] = append(g.resourceWsPool[serverId], ws)
	return
}

func (g *GinServer) resourceWebsocketBroadcast(ctx context.Context, serv server.MinecraftServer, cancelFunc context.CancelFunc) {
	serv.StartMonitorServer()
	var resouceChan chan *models.MonitorMessage
	monitorServer := serv.GetServerMonitor()
	if serv.GetServerMonitor() == nil {
		cancelFunc()
	} else {
		resouceChan = monitorServer.GetMessageChan()
	}

	serverId := serv.GetServerConf().EntryId
	var resourceMsg *models.MonitorMessage
	for {
		select {
		case resourceMsg = <-resouceChan:
			if len(g.resourceWsPool[serverId]) == 0 {
				cancelFunc()
			}
			g.lockeResourceWsPool.Lock()
			for i, ws := range g.resourceWsPool[serverId] {
				if err := ws.WriteJSON(resourceMsg); err != nil {
					// 删除无用ws
					g.resourceWsPool[serverId] = append(g.resourceWsPool[serverId][:i], g.resourceWsPool[serverId][i+1:len(g.resourceWsPool[serverId])]...)
					ws.Close()
					if i < len(g.resourceWsPool[serverId]) {
						i--
					}
				}
			}
			g.lockeResourceWsPool.Unlock()
		case <-ctx.Done():
			g.lockeResourceWsPool.Lock()
			g.resourceWsPool = make(map[string][]*websocket.Conn)
			if monitorServer != nil {
				serv.StopMonitorServer()
			}
			g.lockeResourceWsPool.Unlock()
			return
		}
	}
}

func (g *GinServer) appendStdWsToPool(serverId string, ws *websocket.Conn) {
	g.lock.Lock()
	defer g.lock.Unlock()
	mcContainer := GetMinecraftServerContainerInstance()
	mcServ, err := mcContainer.GetServerById(serverId)
	if err != nil {
		ws.Close()
		return
	}
	g.lockeStdoutWsPool.Lock()
	defer g.lockeStdoutWsPool.Unlock()
	g.stdoutWsPool[serverId] = append(g.stdoutWsPool[serverId], ws)
	if _, ok := g.stdoutChans[serverId]; !ok {
		g.stdoutChans[serverId] = make(chan *models.ReciveMessage, 10)
		mcServ.RegisterSubscribeMessageChan(g.stdoutChans[serverId])
		go g.stdoutWebsocketBroadcast(serverId)
	}
}

func (g *GinServer) stdoutWebsocketBroadcast(serverId string) {
	for {
		msg := <-g.stdoutChans[serverId]
		g.lockeStdoutWsPool.Lock()
		for i, stdoutWs := range g.stdoutWsPool[serverId] {
			if err := stdoutWs.WriteJSON(msg); err != nil {
				// 删除无用ws
				g.stdoutWsPool[serverId] = append(g.stdoutWsPool[serverId][:i], g.stdoutWsPool[serverId][i+1:len(g.stdoutWsPool[serverId])]...)
				stdoutWs.Close()
				if i < len(g.stdoutWsPool[serverId]) {
					i--
				}
			}
		}
		g.lockeStdoutWsPool.Unlock()
	}
}

func (g *GinServer) listenStdinFromWs(serverId string, ws *websocket.Conn) {
	commandReq := &models.Command{}
	for {
		err := ws.ReadJSON(commandReq)
		if err != nil {
			return
		}
		if ok := RunOneCommand(serverId, commandReq.Command, commandReq.Type); !ok {
			return
		}
	}
}

func RunOneCommand(serverId, command string, commandType int64) bool {
	mcContainer := GetMinecraftServerContainerInstance()
	mcServ, err := mcContainer.GetServerById(serverId)
	if err != nil {
		return false
	}
	// TODO 分发给各插件
	if commandType == constant.PLUGIN_COMMAND_TYPE || commandType == constant.ALL_COMMAND_TYPE {

	}

	// 运行服务端命令
	if commandType == constant.SERVER_COMMAND_TYPE || commandType == constant.ALL_COMMAND_TYPE {
		// 如果是关闭服务器，走容器关闭
		if command == "stop" || command == "/stop" {
			_ = mcContainer.StopById(serverId)
		} else {
			_ = mcServ.Command(command)
		}
	}
	WriteLogToDefault(fmt.Sprintf("web后台运行命令：%s", command))
	if command == "stop" || command == "/stop" {
		return false
	}
	return true
}

func GetGinServerInstance() server.GinServer {
	if ginServerInstance != nil {
		return ginServerInstance
	}
	portStr := GetConfVal(constant.MANAGE_HTTP_SERVER_PORT)
	port, _ := strconv.Atoi(portStr)

	// 添加日志
	logger := AddLog(constant.GIN_LOG_NAME)
	gin.DefaultWriter = logger
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	httpServer := getHttpServerObj(port, router)

	ginServerInstance = &GinServer{
		router:              router,
		httpServer:          httpServer,
		port:                port,
		lock:                &sync.Mutex{},
		lockeResourceWsPool: &sync.Mutex{},
		lockeStdoutWsPool:   &sync.Mutex{},
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

func AppendResourceWsToPool(ctx context.Context, serverId string, ws *websocket.Conn) {
	ginServerInstance.appendResourceWsToPool(ctx, serverId, ws)
}

func AppendStdWsToPool(serverId string, ws *websocket.Conn) {
	ginServerInstance.appendStdWsToPool(serverId, ws)
}

func ListenStdinFromWs(serverId string, ws *websocket.Conn) {
	go ginServerInstance.listenStdinFromWs(serverId, ws)
}
