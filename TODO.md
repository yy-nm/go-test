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
