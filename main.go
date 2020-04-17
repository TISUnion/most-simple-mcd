package main

import (
	"fmt"
	_ "fmt"
	"github.com/TISUnion/most-simple-mcd/constant"
	"github.com/TISUnion/most-simple-mcd/modules"
	"github.com/TISUnion/most-simple-mcd/utils"
	"runtime"
)


func main() {
	// 打印banner
	fmt.Println(modules.DrawBanner("M C D"))

	// 启动各服务
	modules.ModuleInit()

	// 打印信息
	webManageUrl := modules.GetWebManageUrl()
	// 如果是windows就自动打开浏览器
	if runtime.GOOS == constant.OS_WINDOWS {
		utils.OpenBrowser(webManageUrl)
	}
	fmt.Println("启动成功！")
	fmt.Println("管理后台url: " + webManageUrl)
	fmt.Println("初始账号: " + constant.DEFAULT_ACCOUNT)
	fmt.Println("初始密码: " + constant.DEFAULT_PASSWORD)
	select {}
}
