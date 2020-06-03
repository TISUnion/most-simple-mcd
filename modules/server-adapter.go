package modules

import "github.com/TISUnion/most-simple-mcd/constant"

type ServerAdapter struct {
	side string
}

func (sa *ServerAdapter) GetVersionRegularExpression() string {
	switch sa.side {
	case constant.VANILLA:
		return constant.VANILLA_VERSION
	default:
		return ""
	}
}

func (sa *ServerAdapter) GetGameTypeRegularExpression() string {
	switch sa.side {
	case constant.VANILLA:
		return constant.VANILLA_GAME_TYPE
	default:
		return ""
	}
}

func (sa *ServerAdapter) GetGameStartRegularExpression() string {
	switch sa.side {
	case constant.VANILLA:
		return constant.VANILLA_GAME_START
	default:
		return ""
	}
}

func (sa *ServerAdapter) GetGameSaveRegularExpression() string {
	switch sa.side {
	case constant.VANILLA:
		return constant.VANILLA_GAME_SAVE
	default:
		return ""
	}
}