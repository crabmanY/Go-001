package Week06

import (
	"errors"
	"strconv"
	"sync"
	"time"
)

//滑动窗口计数器实现
type rollingNumber struct {
	time                    func() int64        //获取当前时间
	timeInMilliseconds      int64               //窗口时间宽度
	numberOfBuckets         int64               //bucket数量
	bucketSizeInMillseconds int64               //单位桶时间宽度
	bucketCircularArray     bucketCircularArray //环形通
	TryLockMutex                                //非阻塞锁
}

//构建滑动窗口计数器
func NewRollingNumber(time func() int64, timeInMilliseconds int64, numberOfBuckets int64) (*rollingNumber, error) {
	if timeInMilliseconds%numberOfBuckets != 0 {
		return nil, errors.New("The timeInMilliseconds must divide equally into numberOfBuckets" +
			". For example 1000/10 is ok, 1000/11 is not.")
	}

	return &rollingNumber{
		time:                    time,
		timeInMilliseconds:      timeInMilliseconds,
		numberOfBuckets:         numberOfBuckets,
		bucketSizeInMillseconds: timeInMilliseconds / numberOfBuckets,
		bucketCircularArray:     NewBucketCircularArray(numberOfBuckets),
	}, nil

}

//对于给定的事件类型，将当前buckets中计数+1
func (rollingNumber *rollingNumber) increment(event RollingNumberEvent) {
	adder, _ := rollingNumber.getCurrentBuckets().getLongAdder(event)
	adder.increment()
}

//获取当前bucket
func (rollingNumber *rollingNumber) getCurrentBuckets() *bucket {
	currentTime := rollingNumber.time()
	get := rollingNumber.bucketCircularArray.listStates.get()
	currentBucket := get.getTail()

	//如果在“时间窗口”内，返回当前窗口
	//只要不在窗户后面，就用最新的
	if currentBucket != nil && currentTime < currentBucket.windowStart+rollingNumber.bucketSizeInMillseconds {
		return currentBucket
	}

	//如果上面没有找到当前的bucket，那就必须创建一个
	if rollingNumber.TryLock() {
		defer rollingNumber.Unlock()
		get := rollingNumber.bucketCircularArray.listStates.get()
		if get.getTail() == nil {
			bucket := NewBucket(currentTime)
			rollingNumber.bucketCircularArray.addLast(bucket)
			return bucket
		} else {
			//进入循环创建尽可能多的bucket，赶上当前窗口
			for i := 0; i < int(rollingNumber.numberOfBuckets); i++ {
				get := rollingNumber.bucketCircularArray.listStates.get()
				lastBucket := get.getTail()
				//我们在“时间窗口”内，返回当前窗口
				if currentTime < lastBucket.windowStart+rollingNumber.bucketSizeInMillseconds {
					return lastBucket
				} else if currentTime-(lastBucket.windowStart+rollingNumber.bucketSizeInMillseconds) > rollingNumber.timeInMilliseconds {
					//经过的时间比整个滚动计数器要长，所以要清除它，从头开始
					rollingNumber.reset()
					//递归调用getCurrentBucket，它将创建一个新的bucket
					return rollingNumber.getCurrentBuckets()
				} else {
					//我们已经过了窗口，所以我们需要创建一个新的bucket
					//创建一个新的bucket并将其添加为新的“last”
					rollingNumber.bucketCircularArray.addLast(NewBucket(lastBucket.windowStart + rollingNumber.bucketSizeInMillseconds))
				}

			}
			//遍历创建完成之后返回最新的bucket
			get := rollingNumber.bucketCircularArray.listStates.get()
			lastBucket := get.getTail()
			return lastBucket
		}

	} else {
		//没有获取到锁
		get := rollingNumber.bucketCircularArray.listStates.get()
		lastBucket := get.getTail()
		if lastBucket != nil {
			//返回最新的
			return lastBucket
		} else {
			//多线程的情况下休眠5秒的，等待其他的线程创建bucket
			time.Sleep(5)
			return rollingNumber.getCurrentBuckets()
		}

	}

}

//Bucket 定义
type bucket struct {
	windowStart           int64                 //标识是那一秒的桶的数据
	adderForCounterType   []*longAdder          //简单自增统计
	updaterForCounterType []longMaxUpdaterAdder //最大并发统计
}

//初始化Bucket
func NewBucket(startTime int64) *bucket {
	adderForCounterType := make([]*longAdder, 10)
	updaterForCounterType := make([]longMaxUpdaterAdder, 10)

	//普通自增类型，预分配内存，将不同的事件分发到不同的index
	counterValues := GetEventValues(Counter)
	for _, val := range counterValues {
		adderForCounterType[val.num] = NewLongAdder()
	}

	//计算最大并发，预分配内存
	maxUpdaterValues := GetEventValues(MaxUpdater)
	for _, val := range maxUpdaterValues {
		updaterForCounterType[val.num] = NewLongMaxUpdaterAdder()
	}

	return &bucket{
		windowStart:           startTime,
		adderForCounterType:   adderForCounterType,
		updaterForCounterType: updaterForCounterType,
	}
}

//获取某个类型下的计数
func (bucket *bucket) get(event RollingNumberEvent) (int64, error) {
	if event.dataType == Counter {
		return bucket.adderForCounterType[event.num].counter, nil
	}

	if event.dataType == MaxUpdater {
		return bucket.updaterForCounterType[event.num].max, nil
	}

	return -1, errors.New("unknown type of event " + strconv.Itoa(int(event.dataType)))
}

//获取某个类型下的long adder
func (bucket *bucket) getLongAdder(event RollingNumberEvent) (*longAdder, error) {
	if event.dataType == MaxUpdater {
		return nil, errors.New("type is not a Counter")
	}

	return bucket.adderForCounterType[event.num], nil
}

//获取某个类型下的 maxUpdater adder
func (bucket *bucket) getMaxUpdaterAdder(event RollingNumberEvent) (longMaxUpdaterAdder, error) {
	if event.dataType == MaxUpdater {
		return longMaxUpdaterAdder{}, errors.New("type is not a Counter")
	}

	return bucket.updaterForCounterType[event.num], nil
}

//包装单个桶
type atomicBucket struct {
	buckets []bucket
	sync.RWMutex
}

//安全替换桶
func (atomicBucket *atomicBucket) set(i int, bucket bucket) {
	atomicBucket.Lock()
	defer atomicBucket.Unlock()
	atomicBucket.buckets[i] = bucket
}

//安全替换桶
func (atomicBucket *atomicBucket) get(i int) *bucket {
	atomicBucket.Lock()
	defer atomicBucket.Unlock()
	return &atomicBucket.buckets[i]
}

//环形桶的状态记录不可变，保证并发安全
type atomicListState struct {
	listStates listState
	sync.RWMutex
}

//安全设置数据
func (atomicState *atomicListState) set(listState listState) {
	atomicState.Lock()
	defer atomicState.Unlock()
	atomicState.listStates = listState
}

//安全取出
func (atomicState *atomicListState) get() listState {
	atomicState.RLock()
	defer atomicState.RUnlock()
	return atomicState.listStates
}

//环形桶
type listState struct {
	atomicBuckets atomicBucket //通过锁包装的环形桶
	size          int          //大小
	tail          int          //尾部引用
	head          int          //头部引用
	dataLength    int          //数据长度
}

//初始化ListState
func NewListState(atomicBuckets atomicBucket, dataLength int, head int, tail int) *listState {
	var size int
	if head == 0 && tail == 0 {
		size = 0
	} else {
		size = (tail + dataLength - head) % dataLength
	}
	return &listState{
		atomicBuckets: atomicBuckets,
		size:          size,
		tail:          tail,
		head:          head,
	}
}

//获取尾部的bucket
func (listState *listState) getTail() *bucket {
	if listState.size == 0 {
		return nil
	} else {
		return listState.atomicBuckets.get((listState.size - 1 + listState.head) % listState.dataLength)
	}

}

//环形桶包装
type bucketCircularArray struct {
	listStates atomicListState
	dataLength int   //桶的长度
	numBuckets int64 //桶的数量
}

//添加bucket TODO
func (bucketarrary *bucketCircularArray) addLast(bucket *bucket) {

}

//重置所有桶的计数 //TODO
func (rollingNumber *rollingNumber) reset() {

}

//初始环形桶
func NewBucketCircularArray(size int64) bucketCircularArray {

	state := atomicListState{
		listStates: listState{
			atomicBuckets: atomicBucket{
				buckets: make([]bucket, size+1),
			},
			tail: 0,
			head: 0,
		},
	}

	return bucketCircularArray{
		listStates: state,
		dataLength: len(state.listStates.atomicBuckets.buckets),
		numBuckets: size,
	}
}

func GetCurrentMillions() int64 {
	return time.Now().Unix()
}
