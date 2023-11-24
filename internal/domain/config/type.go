package config

// Provider 定义提供方需要具备的能力
// 通过Key查询值
// 设置键值对
// 设置默认参数

type Provider interface {
	Get(key string) any
	Set(key string, value any)
	SetDefaults(params Params)
	GetString(key string) string
	IsSet(key string) bool
}

type CompositeConfig interface {
	Base() Provider
	Layer() Provider
}

// Params 参数格式定义
// 关键字为字符类型
// 值为通用类型any
type Params map[string]any

// Set 根据新传入参数，对应层级进行重写
// pp为新传入参数
// p为当前参数
// 将pp的值按层级结构写入p
// 递归完成
func (p Params) Set(pp Params) {
	for k, v := range pp {
		vv, found := p[k]
		if !found {
			p[k] = v
		} else {
			switch vvv := vv.(type) {
			case Params:
				if pv, ok := v.(Params); ok {
					vvv.Set(pv)
				} else {
					p[k] = v
				}
			default:
				p[k] = v
			}
		}
	}
}
