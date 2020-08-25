package utils

import (
	"fmt"
	"github.com/TISUnion/most-simple-mcd/constant"
	"strconv"
	"strings"
)

// 解析mc玩家发言插件命令
func ParsePluginCommand(msg string) (command string, params []string) {
	msgSub := strings.Fields(msg)
	command = msgSub[0]
	if len(msgSub) > 1 {
		params = msgSub[1:]
	}
	return command, params
}

// 比较mc版本 1表示大于，0表示等于，-1表示小于
func CompareMcVersion(mainVersion, compareVersion string) int {
	if mainVersion == compareVersion {
		return constant.COMPARE_EQ
	}
	aMainVersionStr := strings.Split(mainVersion, ".")
	aCompareVersionStr := strings.Split(compareVersion, ".")
	laMainVersionStr := len(aMainVersionStr)
	laCompareVersionStr := len(aCompareVersionStr)
	shortLen := 0
	if laMainVersionStr > laCompareVersionStr {
		shortLen = laCompareVersionStr
	} else {
		shortLen = laMainVersionStr
	}

	// 比对每个次版本
	for i := 0; i < shortLen; i++ {
		mainSubVersion, _ := strconv.ParseInt(aMainVersionStr[i], 10, 64)
		compareSubVersion, _ := strconv.ParseInt(aCompareVersionStr[i], 10, 64)
		if mainSubVersion > compareSubVersion {
			return constant.COMPARE_GT
		} else if mainSubVersion < compareSubVersion {
			return constant.COMPARE_LT
		}
	}
	// 比对最小版本
	if laMainVersionStr > laCompareVersionStr {
		return constant.COMPARE_GT
	} else if laMainVersionStr < laCompareVersionStr {
		return constant.COMPARE_LT
	}

	return constant.COMPARE_EQ
}

// 获取cmd命令数组
func GetCommandArr(memory int64, runPath string) []string {
	return []string{
		"java",
		"-jar",
		fmt.Sprintf("-Xmx%dM", memory),
		fmt.Sprintf("-Xms%dM", memory),
		runPath,
		"nogui",
	}
}
