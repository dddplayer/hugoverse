package valueobject

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
