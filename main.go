package main

import (
	"fmt"
	_ "fmt"
	"github.com/TISUnion/most-simple-mcd/constant"
	"github.com/TISUnion/most-simple-mcd/modules"
	pack_webfile "github.com/TISUnion/most-simple-mcd/pack-webfile"
	"github.com/TISUnion/most-simple-mcd/plugins"
	"github.com/TISUnion/most-simple-mcd/utils"
	"gopkg.in/ini.v1"
	"os"
	"os/signal"
	"runtime"
	"strconv"
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

// 获取管理后台url
func getWebManageUrl() string {
	portStr := strconv.Itoa(constant.DEFAULT_MANAGE_HTTP_SERVER_PORT)
	url := constant.LOACALHOST_URL
	if modules.GetConfVal(constant.MANAGE_HTTP_SERVER_PORT) != portStr {
		url = fmt.Sprintf("%s:%s", url, portStr)
	}
	return url
}

// 初始化
func moduleInit() {
	fmt.Println(modules.DrawBanner("M C D"))
	go exitHandle()
	ini.PrettyFormat = false
	plugins.RegisterPlugin()
	// 解压前端静态文件
	pack_webfile.UnCompress()
}
func main() {
	moduleInit()
	modules.GetMinecraftServerContainerInstance()
	_ = modules.GetGinServerInstance().Start()

	// 如果是windows就自动打开浏览器
	//if runtime.GOOS == constant.OS_WINDOWS {
	webManageUrl := getWebManageUrl()
	if runtime.GOOS == constant.OS_DARWIN {
		utils.OpenBrowser(webManageUrl)
	}
	fmt.Println("启动成功！")
	fmt.Println("管理后台url: " + webManageUrl)
	select {}
}
