学习笔记
#### 作业问题
基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够 一个退出，全部注销退出。

#### 作业思路
使用`errgroup`以及`context`解决级联取消的问题,使用`channel`解决`signal`信号的监听,以及各个`gorountine`之间的通讯的问题
,以达到优雅退出的目的.

#### 监听到signal信号的运行结果日志
```
received os signal, ready cancel other running server
context canceled
端口:8088 shutdown succrss
端口:8089 shutdown succrss
http: Server closed
server  graceful shutdown completed, success total is 2,fail total is 0

```