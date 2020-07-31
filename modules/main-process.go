package modules

import (
	"fmt"
	"github.com/TISUnion/most-simple-mcd/constant"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var ExitChan chan os.Signal

// 优雅的处理退出
func ExitHandle() {
	// 退出信号管道
	ExitChan = make(chan os.Signal)
	signal.Notify(ExitChan, os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGTERM)
	<-ExitChan
	RunDestructCallBacks()
	time.Sleep(time.Second) // wait write log
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
