# goffee
A lightweight Go framework/一个轻量级的Go语言开发框架

## 前言
It's still being built.
It will integrate Web, RPC, ORM frameworks in one.

还在建造中，它将会集合Web,RPC,ORM框架于一身


## 日志
Day1 实现了路由映射表，提供了用户注册静态路由的方法，包装了启动服务的函数

Day2 提供了对 Method 和 Path 这两个常用属性的直接访问。
提供了访问Query和PostForm参数的方法。
提供了快速构造String/Data/JSON/HTML响应的方法。
将路由相关功能单独抽离了出来

Day3 通过前缀树实现了动态路由的功能，并实现了分组路由

Day4 实现了中间件功能并提供了日志中间件