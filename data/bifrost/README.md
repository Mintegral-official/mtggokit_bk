# Bifrost

Bifrost取自北欧神话中是连结阿斯加德（Asgard）和 米德加尔特（中庭/Midgard）的巨大彩虹桥，意即“摇晃的天国道路”。
在这里寓意“数据的传输道路”，它是一个数据加载工具，可以将本地文件、远程文件、mongo等多种数据源的数据加载至本机内存中，
同时支持多种内存模型，满足用户的不同需求

# 支持语言
* go

# 支持平台
* Linux

# 特性
* 多种数据源
* 



# UML设计

![bifrost_arch_uml](pic/bifront_arch.png)

设计上主要分为三个组件
1. Bifrost: 用户接口，
   1. 注册、管理、更新Streamer
   2. 提供数据查询接口
2. Streamer
   1. 数据源的抽象
   2. 负责数据的更新，解析
3. Container
   1. 数据的容器
   2. 负责更新、维护内存中的数据

## Streamer与Container的关系

Streamer代表数据源，Container则代表数据的内存组织方式。

通常情况下Streamer跟Container是一对一的关系，特殊情况下也会出现一对多或者多对一的关系

- 一对一：
- 一对多：DirStreamer
- 多对一：索引，基准增量来自不同的数据源

## Streamer数据更新

### 数据更新模式

1. static 不更新
2. dynamic 动态全量更新
3. increase  全量更新一次，之后动态更定增量
4. dynInc 定时全量更细，动态增量更新

### 数据更新方式

1. sync 同步更新
2. async  异步更新

### Streamer之间的关系

Streamer与Streamer之间可以存在依赖关系，数据更新时需要先更新被依赖的Streamer

备选方案

	1. 多个Streamer可以形成DAG，Bifrost按DAG更新streamer
 	2. 直接依赖

### Streamer更新流程

Streamer根据数据源的不同会分为主动更新和被动更新

主动更新会实时监控数据源变化，发现变化会自动触发更新，不需要Bifrost调度模块管理

被动更新会定期的触发更新，由Bifrost调度模块统一管理

对于同一个Streamer，若上次更新未完成，可根据用户配置决定是否终止上次更新

数据更新失败用户可自定义错误处理函数

Bifrost调度模块采用邮件队列

## FileStreamer

FileStreamer代表本地文件文件

1. 全量更新
2. 增量更新
3. 绑定Parser
4. 错误处理
   1. 回调方式（）
   2. 打日志
   3. 出现错误是否终止本次更新（用户可配置）

## DirStreamer



## MongoStreamer





# Example

