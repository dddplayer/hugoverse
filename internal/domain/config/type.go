package config

import "github.com/dddplayer/hugoverse/internal/domain/config/valueobject"

// Provider 定义提供方需要具备的能力
// 通过Key查询值
// 设置键值对
// 设置默认参数

type Provider interface {
	Get(key string) any
	Set(key string, value any)
	SetDefaults(params valueobject.Params)
	GetString(key string) string
}
