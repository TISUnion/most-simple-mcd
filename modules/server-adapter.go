package modules

import (
	"github.com/TISUnion/most-simple-mcd/constant"
	"github.com/TISUnion/most-simple-mcd/models"
	"github.com/TISUnion/most-simple-mcd/utils"
	"regexp"
	"strconv"
	"strings"
)

type ServerAdapter struct {
	side string
}

// 获取版本
func (sa *ServerAdapter) GetVersionRegularExpression(originMsg string) (res *models.ReciveMessage, ok bool) {
	res = new(models.ReciveMessage)
	res.Source = constant.MC_STDOUT
	var (
		re    *regexp.Regexp
		match []string
	)
	match = make([]string, 0)
	switch sa.side {
	case constant.VANILLA_SERVER:
		re = regexp.MustCompile(constant.VANILLA_VERSION)
		match = re.FindStringSubmatch(originMsg)
	}
	if len(match) > 1 {
		res.Content = match[1]
		ok = true
	}
	return
}

// 获取游戏模式
func (sa *ServerAdapter) GetGameTypeRegularExpression(originMsg string) (res *models.ReciveMessage, ok bool) {
	res = new(models.ReciveMessage)
	res.OriginData = originMsg
	res.Source = constant.MC_STDOUT
	var (
		re    *regexp.Regexp
		match []string
	)
	match = make([]string, 0)
	switch sa.side {
	case constant.VANILLA_SERVER:
		re = regexp.MustCompile(constant.VANILLA_GAME_TYPE)
		match = re.FindStringSubmatch(originMsg)

	}
	if len(match) > 1 {
		res.Content = match[1]
		ok = true
	}
	return
}

// 游戏开始
func (sa *ServerAdapter) GetGameStartRegularExpression(originMsg string) (res *models.ReciveMessage, ok bool) {
	res = new(models.ReciveMessage)
	res.OriginData = originMsg
	res.Source = constant.MC_STDOUT
	var (
		re    *regexp.Regexp
		match []string
	)
	match = make([]string, 0)
	switch sa.side {
	case constant.VANILLA_SERVER:
		re = regexp.MustCompile(constant.VANILLA_GAME_START)
		match = re.FindStringSubmatch(originMsg)
	}
	if len(match) > 0 {
		res.Content = match[0]
		ok = true
	}
	return
}

// 游戏保存（结束）
func (sa *ServerAdapter) GetGameSaveRegularExpression(originMsg string) (res *models.ReciveMessage, ok bool) {
	res = new(models.ReciveMessage)
	res.OriginData = originMsg
	res.Source = constant.MC_STDOUT
	var (
		re    *regexp.Regexp
		match []string
	)
	match = make([]string, 0)
	switch sa.side {
	case constant.VANILLA_SERVER:
		re = regexp.MustCompile(constant.VANILLA_GAME_SAVE)
		match = re.FindStringSubmatch(originMsg)
	}
	if len(match) > 0 {
		res.Content = match[0]
		ok = true
	}
	return
}

//=================================游戏事件
// 解析服务器消息
func (sa *ServerAdapter) GetMessageRegularExpression(originMsg string) (res *models.ReciveMessage, ok bool) {
	res = new(models.ReciveMessage)
	res.OriginData = originMsg
	res.Source = constant.MC_STDOUT
	var (
		re    *regexp.Regexp
		match []string
	)
	match = make([]string, 0)
	switch sa.side {
	case constant.VANILLA_SERVER, constant.BUKKIT14:
		re = regexp.MustCompile(constant.VANILLA_MESSAGE)
		match = re.FindStringSubmatch(originMsg)
	case constant.BUKKIT_SERVER:
		re = regexp.MustCompile(constant.BUKKIT_MESSAGE)
		match = re.FindStringSubmatch(originMsg)
	}
	if len(match) > 3 {
		// 非发言参数
		res.Content = match[3]
		res.LoggingLevel = match[2]
		res.Time = match[1]
		timeElm := strings.Split(match[1], ":")
		res.Hour, _ = strconv.Atoi(timeElm[0])
		res.Minute, _ = strconv.Atoi(timeElm[1])
		res.Second, _ = strconv.Atoi(timeElm[2])
		// 玩家发言参数
		re = regexp.MustCompile(constant.PLAYER_MESSAGE)
		match = re.FindStringSubmatch(match[3])
		if len(match) > 2 {
			res.Player = match[1]
			res.Speak = match[2]
			res.Command, res.Params = utils.ParsePluginCommand(match[2])
			res.IsPlayer = true
			res.IsUser = true
		}
		ok = true
	}
	return
}

// 加入游戏
func (sa *ServerAdapter) GetPlayerJoinRegularExpression(originMsg string) (res *models.ReciveMessage, ok bool) {
	res = new(models.ReciveMessage)
	res.OriginData = originMsg
	res.Source = constant.MC_STDOUT
	var (
		re    *regexp.Regexp
		match []string
	)
	match = make([]string, 0)
	switch sa.side {
	case constant.VANILLA_SERVER:
		re = regexp.MustCompile(constant.VANILLA_PLAYER_JOIN)
		match = re.FindStringSubmatch(originMsg)

	case constant.BUKKIT_SERVER, constant.BUKKIT14:
		re = regexp.MustCompile(constant.BUKKIT_JOIN)
		match = re.FindStringSubmatch(originMsg)
	}
	if len(match) > 1 {
		res.Player = match[1]
		ok = true
	}
	return
}

// 离开游戏
func (sa *ServerAdapter) GetPlayerLeftRegularExpression(originMsg string) (res *models.ReciveMessage, ok bool) {
	res = new(models.ReciveMessage)
	res.OriginData = originMsg
	res.Source = constant.MC_STDOUT
	var (
		re    *regexp.Regexp
		match []string
	)
	match = make([]string, 0)
	switch sa.side {
	case constant.VANILLA_SERVER, constant.BUKKIT_SERVER, constant.BUKKIT14:
		re = regexp.MustCompile(constant.VANILLA_PLAYER_LEFT)
		match = re.FindStringSubmatch(originMsg)
	}
	if len(match) > 0 {
		res.Player = match[1]
		ok = true
	}
	return
}

// 获得成就
func (sa *ServerAdapter) GetPlayerAdvancementRegularExpression(originMsg string) (res *models.ReciveMessage, ok bool) {
	res = new(models.ReciveMessage)
	res.OriginData = originMsg
	res.Source = constant.MC_STDOUT
	var (
		re    *regexp.Regexp
		match []string
	)
	match = make([]string, 0)
	switch sa.side {
	case constant.VANILLA_SERVER, constant.BUKKIT14, constant.BUKKIT_SERVER:
		re = regexp.MustCompile(constant.VANILLA_PLAYER_ADVANCEMENT)
		match = re.FindStringSubmatch(originMsg)
	}
	if len(match) > 2 {
		res.Player = match[1]
		res.Content = match[2]
		ok = true
	}
	return
}

// 玩家死亡
func (sa *ServerAdapter) GetPlayerDeathRegularExpression(originMsg string) (res *models.ReciveMessage, ok bool) {
	res = new(models.ReciveMessage)
	res.OriginData = originMsg
	res.Source = constant.MC_STDOUT
	var (
		re    *regexp.Regexp
		match []string
	)
	match = make([]string, 0)
	switch sa.side {
	case constant.VANILLA_SERVER, constant.BUKKIT_SERVER, constant.BUKKIT14:
		for _, deathMsg := range constant.DeathMessage {
			re = regexp.MustCompile(deathMsg)
			match = re.FindStringSubmatch(originMsg)
			if len(match) > 1 {
				res.Player = match[1]
				ok = true
				break
			}
		}
	}
	return
}

// 支持的所有服务端类型
func GetAllServerSide() []string {
	return []string{constant.VANILLA_SERVER}
}
