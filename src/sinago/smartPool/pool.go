package smartPool

import (
	"time"
)

type SmartPool struct {

	strategies map[int]int // 每 X 秒 X 次 修改 触发

	pool       map[int] []string // 每 X 秒 的 数据

	onSatisfy  func(contents []string)

}

func NewSmartPool(strategies map[int]int) *SmartPool {
	var mp        = new(SmartPool)
	mp.strategies = make(map[int]int)
	mp.pool       = make(map[int] []string)
	mp.strategies = strategies
	for k := range strategies {
		mp.pool[k] = []string{}
	}
	return mp
}

func (mp *SmartPool) SetSatisfy(onSatisfy func(contents []string)) {
	mp.onSatisfy = onSatisfy
}

func (mp *SmartPool) Add(content string) {
	for k := range mp.pool {
		mp.pool[k] = append(mp.pool[k], content)
	}
}

func (mp *SmartPool) Start() {
	// 按照 map 进行 开 协程 如果 时间之内
	for k, v := range mp.strategies {
		go func(k int, v int) {
			for{
				time.Sleep(time.Duration(k) * time.Second)
				if len(mp.pool[k]) >= v {
					if mp.onSatisfy != nil {
						mp.onSatisfy(mp.pool[k])
					}
					// 完成后 清空
					for k := range mp.pool {
						mp.pool[k] = []string{}
					}
				}
			}
		}(k,v)
	}
}


