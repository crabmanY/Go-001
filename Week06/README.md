学习笔记

#### 作业 

参考hystrix实现滑动窗口计数器

#### 目录结构
```
├── README.md
├── event_adder.go    -----原子计数实现
├── rolloing_event.go -----计数事件类型定义
├── sliding_window.go -----滑动窗口实现
└── trylockmutex.go   -----trylock实现，非阻塞获取锁，快速失败
```

#### 代码说明
`rollingNumber`滑动窗口计数器的具体结构
`rollingNumber.increment()` 计数的具体实现
其余数据接口详见代码注释

#### 思路
`hystrix`的滑动窗口是由环形通实现的，通过`ListState` 保留各个时间窗口的状态，
通过当前环形桶的时间去增加各个类型的统计计数，这里自己实现了`tryLock`，在hystrix的源码中
是通过`trylock`的方式进行获取当前`bucket`。

#### 感悟
`hystrix`的代码还是非常强的，很地方设计的很巧妙。