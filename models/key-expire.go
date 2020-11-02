package models

type ValueWithTTL struct {
	Value      string `json:"value"`
	NanoExpire int64  `json:"nano_expire"` // 过期时间戳(毫秒)
}
