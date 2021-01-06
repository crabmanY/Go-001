package Week06

//计数器事件类型
type RollingNumberEvent struct {
	num      int64
	dataType int8
}

//事件类型
const (
	Counter    = iota //普通自增
	MaxUpdater        //最大并发值
)

//计数器事件类型
var (
	SUCCESS              = RollingNumberEvent{num: 1, dataType: 1} //成功
	FAILURE              = RollingNumberEvent{num: 2, dataType: 1} //失败
	TIMEOUT              = RollingNumberEvent{num: 3, dataType: 1} //超时
	SHORT_CIRCUITED      = RollingNumberEvent{num: 4, dataType: 1} //短路
	THREAD_POOL_REJECTED = RollingNumberEvent{num: 5, dataType: 1} //线程池拒绝
	SEMAPHORE_REJECTED   = RollingNumberEvent{num: 6, dataType: 1} //信号量拒绝
	BAD_REQUEST          = RollingNumberEvent{num: 7, dataType: 1} //请求失败
	THREAD_MAX_ACTIVE    = RollingNumberEvent{num: 8, dataType: 2} //达到最大并发
)

//根据类型获取对应的事件类型
func GetEventValues(dataType int) []RollingNumberEvent {
	switch dataType {
	case 1:
		return []RollingNumberEvent{SUCCESS, FAILURE, TIMEOUT, SHORT_CIRCUITED,
			THREAD_POOL_REJECTED, SEMAPHORE_REJECTED, BAD_REQUEST}
	case 2:
		return []RollingNumberEvent{THREAD_MAX_ACTIVE}
	default:
		return nil
	}

}
