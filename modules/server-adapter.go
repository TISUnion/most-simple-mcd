package modules

import "github.com/TISUnion/most-simple-mcd/constant"

type ServerAdapter struct {
	side string
}
//================================启动
// 获取版本
func (sa *ServerAdapter) GetVersionRegularExpression() string {
	switch sa.side {
	case constant.VANILLA:
		return constant.VANILLA_VERSION
	default:
		return ""
	}
}

// 获取游戏模式
func (sa *ServerAdapter) GetGameTypeRegularExpression() string {
	switch sa.side {
	case constant.VANILLA:
		return constant.VANILLA_GAME_TYPE
	default:
		return ""
	}
}

// 游戏开始
func (sa *ServerAdapter) GetGameStartRegularExpression() string {
	switch sa.side {
	case constant.VANILLA:
		return constant.VANILLA_GAME_START
	default:
		return ""
	}
}

// 游戏保存（结束）
func (sa *ServerAdapter) GetGameSaveRegularExpression() string {
	switch sa.side {
	case constant.VANILLA:
		return constant.VANILLA_GAME_SAVE
	default:
		return ""
	}
}

//=================================游戏事件
// 获取玩家发言 [09:00:00] [Server thread/INFO]: <Steve> Hello
func (sa *ServerAdapter) GetMessageRegularExpression() string {
	switch sa.side {
	case constant.VANILLA:
		return constant.VANILLA_MESSAGE
	default:
		return ""
	}
}

// 加入游戏 Steve joined the game
func (sa *ServerAdapter) GetPlayerJoinRegularExpression() string {
	switch sa.side {
	case constant.VANILLA:
		return constant.VANILLA_PLAYER_JOIN
	default:
		return ""
	}
}

// 离开游戏 Steve left the game
func (sa *ServerAdapter) GetPlayerLeftRegularExpression() string {
	switch sa.side {
	case constant.VANILLA:
		return constant.VANILLA_PLAYER_LEFT
	default:
		return ""
	}
}

// 获得成就 Steve has made the advancement [Stone Age]
func (sa *ServerAdapter) GetPlayerAdvancementRegularExpression() string {
	switch sa.side {
	case constant.VANILLA:
		return constant.VANILLA_PLAYER_ADVANCEMENT
	default:
		return ""
	}
}

// 玩家死亡
func (sa *ServerAdapter) GetPlayerDeathRegularExpression() (message []string) {
	switch sa.side {
	case constant.VANILLA:
		return constant.DeathMessage
	default:
		return
	}
}

// 支持的所有服务端类型
func GetAllServerSide() []string {
	return []string{constant.VANILLA}
}
