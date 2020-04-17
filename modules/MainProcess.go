package modules

import (
	"fmt"
	"github.com/TISUnion/most-simple-mcd/constant"
	pack_webfile "github.com/TISUnion/most-simple-mcd/pack-webfile"
	"github.com/TISUnion/most-simple-mcd/plugins"
	"gopkg.in/ini.v1"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var ExitChan chan os.Signal

// 优雅的处理退出
func exitHandle() {
	// 退出信号管道
	ExitChan = make(chan os.Signal)
	signal.Notify(ExitChan, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGTERM)
	<-ExitChan
	RunDestructCallBacks()
	os.Exit(1)
}

// 发送退出信号
func SendExitSign() {
	ExitChan <- syscall.SIGQUIT
}

// 获取管理后台url
func GetWebManageUrl() string {
	portStr := strconv.Itoa(constant.DEFAULT_MANAGE_HTTP_SERVER_PORT)
	url := constant.LOACALHOST_URL
	if GetConfVal(constant.MANAGE_HTTP_SERVER_PORT) != portStr {
		url = fmt.Sprintf("%s:%s", url, portStr)
	}
	return url
}

// 初始化
func ModuleInit() {
	go exitHandle()
	ini.PrettyFormat = false
	plugins.RegisterPlugin()
	// 解压前端静态文件
	_ = pack_webfile.UnCompress()
	GetMinecraftServerContainerInstance()
	_ = GetGinServerInstance().Start()
}
