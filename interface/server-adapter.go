package _interface

type ServerAdapter interface {
	// 获取版本
	GetVersionRegularExpression() string
	// 获取游戏模式
	GetGameTypeRegularExpression() string
	// 游戏开始
	GetGameStartRegularExpression() string
	// 游戏保存（结束）
	GetGameSaveRegularExpression() string
	// 获取玩家发言
	GetMessageRegularExpression() string
	// 加入游戏
	GetPlayerJoinRegularExpression() string
	// 离开游戏
	GetPlayerLeftRegularExpression() string
	// 获得成就
	GetPlayerAdvancementRegularExpression() string
	// 玩家死亡
	GetPlayerDeathRegularExpression() (message []string)
}
