package Week06

import "sync/atomic"

type longAdder struct {
	counter int64
}

//初始化 LongAdder
func NewLongAdder() *longAdder {
	return &longAdder{
		counter: 0,
	}
}

//TODO 有点难度，短时间无法完成，后续完善
type longMaxUpdaterAdder struct {
	max int64
}

//初始化 LongAdder
func NewLongMaxUpdaterAdder() longMaxUpdaterAdder {
	return longMaxUpdaterAdder{
		max: 0,
	}
}

//原子自增
func (adder *longAdder) increment() {
	atomic.AddInt64(&adder.counter, 1)
}

//原子递减
func (adder *longAdder) decrement() {
	atomic.AddInt64(&adder.counter, -1)
}

//重置
func (adder *longAdder) reset() {
	atomic.StoreInt64(&adder.counter, 0)
}
