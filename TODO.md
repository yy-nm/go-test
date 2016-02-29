#### Task

- 2.15 ~ 2.19
  + 熟悉并练习使用 go 编程, 主要侧重 go 并发编程
  + 主要的参考文档
    * [go doc](https://golang.org/doc/)
    * [Go Concurrency Patterns](https://www.youtube.com/watch?v=f6kdp27TYZs)
      * 主要介绍 goroutine 和 channel, 介绍 goroutine 是个很廉价的东西, 可以轻松达到 10k
    * [Advanced Go Concurrency Patterns](https://www.youtube.com/watch?v=QDDwwePbDtw)
      * 除了介绍 panic 和 go 工具, 其他均是废话
  + 开始编写服务器架构


- 2.22 ~ 2.26
  + 2.24: 初步实现网络部分, 当时分层和模块有几个部分分隔不明确/不清晰, 主要的隔离依靠配置接口, 当前的配置实现是以解析 json 数据
  + 确定新的结构之后, 可以看一下 [Effective Go](https://golang.org/doc/effective_go.html) 充分了解 go
  + 完成 io, stream, connMgr 部分编写


- 2.29 ~ 3.4
  + 2.29: 根据今天进行流量压测发现几个问题:
    * linux goroutine 启动数量比较小(维持在 1K, 主要是 goroutine 涉及到 io 操作所以是阻塞且 linux 有对线程的限制
    , 当前还没找到方法解除), windows 上可以跑到 5k ~ 6k, 不过可能爆内存
    * 不同的 goroutine 对于内存似乎是独立的, 所以当开启多个客户端进行流量攻击时, 内存增加快,
    基本上测试下来, 客户端流量发送总量在 5g 左右, ，每米流量在 87m, go 服务器上占用的内存在 1.1g 左右,
    之前数据填太多直接爆了内存, 爆内存还有一个原因是 window defender 处于开启状态
    单个客户端进行流量攻击时, 流量总量在 4.85g, go 服务器占用的内存在 230m 左右, 持续时间 60s, 每秒流量在 85m 左右
    这块可以**通过实现一个带限制功能的 connMgr 接口即可实现**
    * 鉴于之前的关系, 同时需要设计连接数量限制的功能