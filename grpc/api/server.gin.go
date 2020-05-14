// Code generated by protoc-gen-lightbrother, DO NOT EDIT.

/*
Package api is a generated gin stub package.
This code was generated with protoc-gen-lightbrother. 

It is generated from these files:
	server.proto
*/
package api

import (
	"github.com/gin-gonic/gin"
	"context"
	"net/http"
)

// to suppressed 'imported but not used warning'

const SERVER_HTTP_METGOD = "GRPC"

var PathMcServerListenResource = "/most.simple.mcd.McServer/listenResource"
var PathMcServerServerInteraction = "/most.simple.mcd.McServer/serverInteraction"
var PathMcServerList = "/most.simple.mcd.McServer/list"
var PathMcServerGetServerState = "/most.simple.mcd.McServer/getServerState"
var PathMcServerDetail = "/most.simple.mcd.McServer/detail"
var PathMcServerOperateServer = "/most.simple.mcd.McServer/operateServer"
var PathMcServerUpdateServerInfo = "/most.simple.mcd.McServer/updateServerInfo"

// McServerGinServer is the server API for McServer service.
type McServerGinServer interface {
	// 监听服务端消耗资源
	ListenResource(ctx context.Context, req *ListenResourceReq) (resp *ListenResourceResp, err error)

	// 服务器交互
	ServerInteraction(ctx context.Context, req *ServerInteractionReq) (resp *ServerInteractionResp, err error)

	// 获取服务端信息列表
	List(ctx context.Context, req *ListReq) (resp *ListResp, err error)

	// 获取服务端信息列表
	GetServerState(ctx context.Context, req *GetServerStateReq) (resp *GetServerStateResp, err error)

	// 获取服务端详情
	Detail(ctx context.Context, req *DetailReq) (resp *DetailResp, err error)

	// 操作服务端
	OperateServer(ctx context.Context, req *OperateServerReq) (resp *OperateServerResp, err error)

	// 修改服务端参数
	UpdateServerInfo(ctx context.Context, req *UpdateServerInfoReq) (resp *UpdateServerInfoResp, err error)
}

var apiMcServerSvc McServerGinServer

func listenResource(c *gin.Context) {
	p := new(ListenResourceReq)
	apiMcServerSvc.ListenResource(c, p)
}

func serverInteraction(c *gin.Context) {
	p := new(ServerInteractionReq)
	apiMcServerSvc.ServerInteraction(c, p)
}

func list(c *gin.Context) {
	p := new(ListReq)
	resp, err := apiMcServerSvc.List(c, p)
	if err != nil {
		c.Set("code", -500)
		c.Set("message", err.Error())
		c.JSON(http.StatusOK, getServerResponse(c, nil))
		return
	}
	c.JSON(http.StatusOK, getServerResponse(c, resp))
}

func getServerState(c *gin.Context) {
	p := new(GetServerStateReq)
	if err := c.BindJSON(p); err != nil {
		c.Set("code", -500)
		c.Set("message", err.Error())
		c.JSON(http.StatusOK, getServerResponse(c, nil))
		return
	}
	resp, err := apiMcServerSvc.GetServerState(c, p)
	if err != nil {
		c.Set("code", -500)
		c.Set("message", err.Error())
		c.JSON(http.StatusOK, getServerResponse(c, nil))
		return
	}
	c.JSON(http.StatusOK, getServerResponse(c, resp))
}

func detail(c *gin.Context) {
	p := new(DetailReq)
	if err := c.BindJSON(p); err != nil {
		c.Set("code", -500)
		c.Set("message", err.Error())
		c.JSON(http.StatusOK, getServerResponse(c, nil))
		return
	}
	resp, err := apiMcServerSvc.Detail(c, p)
	if err != nil {
		c.Set("code", -500)
		c.Set("message", err.Error())
		c.JSON(http.StatusOK, getServerResponse(c, nil))
		return
	}
	c.JSON(http.StatusOK, getServerResponse(c, resp))
}

func operateServer(c *gin.Context) {
	p := new(OperateServerReq)
	if err := c.BindJSON(p); err != nil {
		c.Set("code", -500)
		c.Set("message", err.Error())
		c.JSON(http.StatusOK, getServerResponse(c, nil))
		return
	}
	resp, err := apiMcServerSvc.OperateServer(c, p)
	if err != nil {
		c.Set("code", -500)
		c.Set("message", err.Error())
		c.JSON(http.StatusOK, getServerResponse(c, nil))
		return
	}
	c.JSON(http.StatusOK, getServerResponse(c, resp))
}

func updateServerInfo(c *gin.Context) {
	p := new(UpdateServerInfoReq)
	if err := c.BindJSON(p); err != nil {
		c.Set("code", -500)
		c.Set("message", err.Error())
		c.JSON(http.StatusOK, getServerResponse(c, nil))
		return
	}
	resp, err := apiMcServerSvc.UpdateServerInfo(c, p)
	if err != nil {
		c.Set("code", -500)
		c.Set("message", err.Error())
		c.JSON(http.StatusOK, getServerResponse(c, nil))
		return
	}
	c.JSON(http.StatusOK, getServerResponse(c, resp))
}

func RegisterServerMcServerGinServer(e *gin.Engine, server McServerGinServer) {
	apiMcServerSvc = server
	e.Handle("GET", PathMcServerListenResource, listenResource)
	e.Handle("GET", PathMcServerServerInteraction, serverInteraction)
	e.Handle("POST", PathMcServerList, handleServerAuthMiddleware, list)
	e.Handle("POST", PathMcServerGetServerState, handleServerAuthMiddleware, getServerState)
	e.Handle("POST", PathMcServerDetail, handleServerAuthMiddleware, detail)
	e.Handle("POST", PathMcServerOperateServer, handleServerAuthMiddleware, operateServer)
	e.Handle("POST", PathMcServerUpdateServerInfo, handleServerAuthMiddleware, updateServerInfo)
}

// 返回数据格式化
func getServerResponse(c *gin.Context, data interface{}) gin.H {
	responseData := make(map[string]interface{})
	code, ok := c.Get("code")
	if !ok {
		code = 0
	}
	msg, ok := c.Get("message")
	if !ok {
		msg = ""
	}
	responseData["code"] = code
	responseData["message"] = msg
	responseData["data"] = data
	return responseData
}

var (
	serverAuthMiddleware []gin.HandlerFunc
)

func RegisterServerAuthMiddleware(f gin.HandlerFunc) {
	serverAuthMiddleware = append(serverAuthMiddleware, f)
}

func handleServerAuthMiddleware(c *gin.Context) {
	for _, middleware := range serverAuthMiddleware {
		if c.IsAborted() {
			break
		}
		middleware(c)
	}
}

