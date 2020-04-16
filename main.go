package main

import (
	_ "fmt"
	"github.com/TISUnion/most-simple-mcd/modules"
	"github.com/TISUnion/most-simple-mcd/plugins"
	"gopkg.in/ini.v1"
	"os"
	"os/signal"
	"syscall"
)

// 优雅的处理退出
func exitHandle() {
	// 退出信号管道
	exitChan := make(chan os.Signal)
	signal.Notify(exitChan, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGTERM)
	<-exitChan
	modules.RunDestructCallBacks()
	os.Exit(1)
}

func moduleInit() {
	go exitHandle()
	ini.PrettyFormat = false
}
func main() {
	moduleInit()
	plugins.RegisterPlugin()
	modules.GetMinecraftServerContainerInstance()
	_ = modules.GetGinServerInstance().Start()
	select {}
}
