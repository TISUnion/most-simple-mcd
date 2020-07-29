package modules

import "github.com/TISUnion/most-simple-mcd/constant"

type ServerAdapter struct {
	side string
}

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

// 获取玩家发言
func (sa *ServerAdapter) GetMessageRegularExpression() string {
	switch sa.side {
	case constant.VANILLA:
		return constant.VANILLA_MESSAGE
	default:
		return ""
	}
}

// 支持的所有服务端类型
func GetAllServerSide() []string {
	return []string{constant.VANILLA}
}
