package _interface

import "github.com/TISUnion/most-simple-mcd/models"

type ServerAdapter interface {
	// 获取版本
	GetVersionRegularExpression(originMsg string) (res *models.ReciveMessage, ok bool)
	// 获取游戏模式
	GetGameTypeRegularExpression(originMsg string) (res *models.ReciveMessage, ok bool)
	// 游戏开始
	GetGameStartRegularExpression(originMsg string) (res *models.ReciveMessage, ok bool)
	// 游戏保存（结束）
	GetGameSaveRegularExpression(originMsg string) (res *models.ReciveMessage, ok bool)
	// 获取玩家发言
	GetMessageRegularExpression(originMsg string) (res *models.ReciveMessage, ok bool)
	// 加入游戏
	GetPlayerJoinRegularExpression(originMsg string) (res *models.ReciveMessage, ok bool)
	// 离开游戏
	GetPlayerLeftRegularExpression(originMsg string) (res *models.ReciveMessage, ok bool)
	// 获得成就
	GetPlayerAdvancementRegularExpression(originMsg string) (res *models.ReciveMessage, ok bool)
	// 玩家死亡
	GetPlayerDeathRegularExpression(originMsg string) (res *models.ReciveMessage, ok bool)
}
