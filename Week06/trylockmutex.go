package Week06

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

const (
	mutexLocked      = 1 << iota // 加锁标识位置
	mutexWoken                   // 唤醒标识位置
	mutexStarving                // 锁饥饿标识位置
	mutexWaiterShift = iota      // 标识waiter的起始bit位置
)

//尝试获取锁，获取不到不阻塞，直接退出
type TryLockMutex struct {
	sync.Mutex
}

//尝试获取锁
func (tm *TryLockMutex) TryLock() bool {
	//如果是第一次请求加锁成功
	if atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&tm.Mutex)), 0, mutexLocked) {
		return true
	}

	//如果处于加锁,唤醒,饥饿状态,这次请求不参与竞争
	old := atomic.LoadInt32((*int32)(unsafe.Pointer(&tm.Mutex)))
	if old&(mutexLocked|mutexStarving|mutexWoken) != 0 {
		return false
	}

	//在竞争状态下请求锁
	new := old | mutexLocked
	return atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&tm.Mutex)), old, new)
}

//统计等待者的数量
func (tm *TryLockMutex) Count() int {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(&tm.Mutex)))
	state = state >> mutexWaiterShift
	state = state + (state & mutexLocked)
	return int(state)
}

// 锁是否被持有
func (tm *TryLockMutex) IsLocked() bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(&tm.Mutex)))
	return state&mutexLocked == mutexLocked
}

// 是否有等待者被唤醒
func (tm *TryLockMutex) IsWoken() bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(&tm.Mutex)))
	return state&mutexWoken == mutexWoken
}

// 锁是否处于饥饿状态
func (tm *TryLockMutex) IsStarving() bool {
	state := atomic.LoadInt32((*int32)(unsafe.Pointer(&tm.Mutex)))
	return state&mutexStarving == mutexStarving
}
