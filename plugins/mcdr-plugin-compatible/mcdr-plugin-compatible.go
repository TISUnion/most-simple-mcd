package mcdr_plugin_compatible

import (
	"fmt"
	"github.com/TISUnion/most-simple-mcd/constant"
	"github.com/TISUnion/most-simple-mcd/modules"
	uuid "github.com/satori/go.uuid"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	pluginDir = "mcdr-plugins"
	pythonExt = ".py"
)

var mpc *McdrPluginCompatible

type McdrPluginCompatible struct {
	pluginRootPath string
	plugins        map[string]*McdrPlugin
}

func (mpc *McdrPluginCompatible) ChangeConfCallBack() {
}

func (mpc *McdrPluginCompatible) DestructCallBack() {
	PyVmEnd()
}

func (mpc *McdrPluginCompatible) InitCallBack() {
	if !PyVmStart() {
		modules.WriteLogToDefault("python虚拟机开启失败")
		modules.SendExitSign()
	}
	mpc.plugins = make(map[string]*McdrPlugin)
	mpc.ScanPlugin()
}

func (mpc *McdrPluginCompatible) ScanPlugin() {
	err := filepath.Walk(mpc.pluginRootPath, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}
		name := info.Name()
		if filepath.Ext(name) == pythonExt {
			// 兼容windows
			path = strings.ReplaceAll(path, "\\", "/")

			re := regexp.MustCompile(fmt.Sprintf("(.*%s/)|(.py)", pluginDir))
			pluginName := re.ReplaceAllString(path, "")
			packageName := strings.ReplaceAll(pluginName, "/", ".")
			pluginName = strings.ReplaceAll(pluginName, "/", "_")
			// 已经加载过的，则不重复加载
			if _, ok := mpc.plugins[pluginName]; ok {
				return nil
			}
			packageName = fmt.Sprint(pluginDir, ".", packageName)
			pluginInfo := fmt.Sprint("mcdr插件: ", pluginName)
			plp := &McdrPlugin{
				pluginPath:        path,
				packageName:       packageName,
				pluginName:        pluginName,
				id:                uuid.NewV4().String(),
				pluginDescription: pluginInfo,
				helpDescription:   pluginInfo,
			}
			modules.RegisterCallBack(plp)
			mpc.plugins[pluginName] = plp
		}
		return nil
	})
	if err != nil {
		modules.WriteLogToDefault("加载插件失败！")
	}
}

func StartLoadMcdrPlugin() {
	mpc = &McdrPluginCompatible {
		pluginRootPath:filepath.Join(modules.GetConfVal(constant.WORKSPACE), pluginDir),
	}
	modules.RegisterCallBack(mpc)
}

func GetMcdrPluginCompatibleObj() *McdrPluginCompatible{
	if mpc != nil {
		return mpc
	}
	StartLoadMcdrPlugin()
	return mpc
}