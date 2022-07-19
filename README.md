# goffee
A lightweight Go framework/一个轻量级的Go语言开发框架

## 前言
It's still being built.
It will integrate HTTP, Cache, RPC, ORM frameworks in one.

还在建造中，它将会集合Web,RPC,ORM框架于一身


## 日志

### HTTP

2022.7.9 实现了路由映射表，提供了用户注册静态路由的方法，包装了启动服务的函数

2022.7.10 提供了对 Method 和 Path 这两个常用属性的直接访问。
提供了访问Query和PostForm参数的方法。
提供了快速构造String/Data/JSON/HTML响应的方法。
将路由相关功能单独抽离了出来

2022.7.11 通过前缀树实现了动态路由的功能，并实现了分组路由

2022.7.12 实现了中间件功能并提供了日志中间件

2022.7.13 支持模板渲染

2022.7.15 实现错误恢复功能

### Cache

2022.7.17 完成LRU缓存淘汰算法

2022.7.18 使用互斥锁完成封装，支持并发

2022.7.19 完成缓存的HTTP服务端