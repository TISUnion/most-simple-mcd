package modules

import "github.com/TISUnion/most-simple-mcd/constant"

type ServerAdapter struct {
	side string
}

// 获取版本
func (sa *ServerAdapter) GetVersionRegularExpression() string {
	switch sa.side {
	case constant.VANILLA_SERVER:
		return constant.VANILLA_VERSION
	default:
		return ""
	}
}

// 获取游戏模式
func (sa *ServerAdapter) GetGameTypeRegularExpression() string {
	switch sa.side {
	case constant.VANILLA_SERVER:
		return constant.VANILLA_GAME_TYPE
	default:
		return ""
	}
}

// 游戏开始
func (sa *ServerAdapter) GetGameStartRegularExpression() string {
	switch sa.side {
	case constant.VANILLA_SERVER:
		return constant.VANILLA_GAME_START
	default:
		return ""
	}
}

// 游戏保存（结束）
func (sa *ServerAdapter) GetGameSaveRegularExpression() string {
	switch sa.side {
	case constant.VANILLA_SERVER:
		return constant.VANILLA_GAME_SAVE
	default:
		return ""
	}
}

//=================================游戏事件
// 获取玩家发言
func (sa *ServerAdapter) GetMessageRegularExpression() string {
	switch sa.side {
	case constant.VANILLA_SERVER, constant.BUKKIT14:
		return constant.VANILLA_MESSAGE
	case constant.BUKKIT_SERVER:
		return constant.BUKKIT_MESSAGE
	default:
		return ""
	}
}

// 加入游戏
func (sa *ServerAdapter) GetPlayerJoinRegularExpression() string {
	switch sa.side {
	case constant.VANILLA_SERVER:
		return constant.VANILLA_PLAYER_JOIN
	case constant.BUKKIT_SERVER, constant.BUKKIT14:
		return constant.BUKKIT_JOIN
	default:
		return ""
	}
}

// 离开游戏
func (sa *ServerAdapter) GetPlayerLeftRegularExpression() string {
	switch sa.side {
	case constant.VANILLA_SERVER, constant.BUKKIT_SERVER, constant.BUKKIT14:
		return constant.VANILLA_PLAYER_LEFT
	default:
		return ""
	}
}

// 获得成就
func (sa *ServerAdapter) GetPlayerAdvancementRegularExpression() string {
	switch sa.side {
	case constant.VANILLA_SERVER, constant.BUKKIT14:
		return constant.VANILLA_PLAYER_ADVANCEMENT
	case constant.BUKKIT_SERVER:
		return constant.VANILLA_PLAYER_ADVANCEMENT
	default:
		return ""
	}
}

// 玩家死亡
func (sa *ServerAdapter) GetPlayerDeathRegularExpression() (message []string) {
	switch sa.side {
	case constant.VANILLA_SERVER, constant.BUKKIT_SERVER, constant.BUKKIT14:
		return constant.DeathMessage
	default:
		return
	}
}

// 支持的所有服务端类型
func GetAllServerSide() []string {
	return []string{constant.VANILLA_SERVER}
}
