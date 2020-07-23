package mcdr_plugin_compatible

import (
	"github.com/TISUnion/most-simple-mcd/constant"
	"github.com/TISUnion/most-simple-mcd/modules"
	"os"
	"path/filepath"
)

type McdrPluginCompatible struct {
	pluginRootPath string
	loadPluginPath []string
	loadPlugin map[string]string
}

func (mpc *McdrPluginCompatible) ChangeConfCallBack() {
}

func (mpc *McdrPluginCompatible) DestructCallBack() {
}

func (mpc *McdrPluginCompatible) InitCallBack() {
	mpc.pluginRootPath = modules.GetConfVal(constant.WORKSPACE)
}

// TODO
func(mpc *McdrPluginCompatible) scanPlugin(path string) {
	_ = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			mpc.scanPlugin(info.Name())
		} else {

		}
		return nil
	})
}