package main

import (
	"fmt"
	_ "fmt"
	"github.com/TISUnion/most-simple-mcd/constant"
	"github.com/TISUnion/most-simple-mcd/modules"
	pack_webfile "github.com/TISUnion/most-simple-mcd/pack-webfile"
	"github.com/TISUnion/most-simple-mcd/plugins"
	"github.com/TISUnion/most-simple-mcd/services"
	"github.com/TISUnion/most-simple-mcd/utils"
	"gopkg.in/ini.v1"
	"runtime"
)

// 初始化
func ModuleInit() {
	// 处理退出逻辑
	go modules.ExitHandle()
	ini.PrettyFormat = false
	// 注册插件到容器中
	plugins.RegisterPlugin()
	// 解压前端静态文件
	_ = pack_webfile.UnCompress()
	// 创建服务端容器对象
	modules.GetMinecraftServerContainerInstance()
	// 注册gin路由
	services.RegisterServices()
	// 开启web服务器
	_ = modules.GetGinServerInstance().Start()
}

func main() {
	// 打印banner
	fmt.Println(modules.DrawBanner("M C D"))

	// 启动各服务
	ModuleInit()

	// 打印信息
	webManageUrl := modules.GetWebManageUrl()
	// 如果是windows就自动打开浏览器
	if runtime.GOOS == constant.OS_WINDOWS {
		utils.OpenBrowser(webManageUrl)
	}
	fmt.Println("启动成功！")
	fmt.Println("本项目开源并且免费，项目仓库：https://github.com/TISUnion/most-simple-mcd")
	fmt.Println("管理后台url: " + webManageUrl)
	fmt.Println("初始账号: " + constant.DEFAULT_ACCOUNT)
	fmt.Println("初始密码: " + constant.DEFAULT_PASSWORD)
	select {}
}
