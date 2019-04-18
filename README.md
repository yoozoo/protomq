# protomq

[![Build Status](https://travis-ci.org/yoozoo/protomq.svg?branch=master)](https://travis-ci.org/yoozoo/protomq)

## 初始化

```bash
go get  github.com/yoozoo/protomq
protomq init
```

### php

php项目需要先使用composer安装依赖

```bash
composer require google/protobuf
composer require spiral/roadrunner
```

* https://github.com/protocolbuffers/protobuf/tree/master/php

### go

```bash
go get github.com/spiral/roadrunner
go get -u github.com/golang/protobuf/protoc-gen-go
```

## 使用范例

生成`go 生产者客户端`

```bash
./protomq.exe gen --lang=goproducer ./output_folder ./test.proto
```

生成`go 消费者服务器端`

```bash
./protomq.exe gen --lang=goconsumer ./output_folder ./test.proto
```

生成`php 消费者服务器端`详情在[这里](docs/php_consumer_zh.md)。另外`php 消费者简单类型服务器端`的文档在这里[这里](docs/php_consumer_simpletype_zh.md)。

生成`php 生产者客户端`详情在[这里](docs/php_producer_zh.md)。

## TODO

* `protomq` cli
  * [x] 自动下载`protoc`
  * [x] 内嵌`protomq.proto`
  * [x] 嵌套调用`protoc`
  * [X] protoc文件语法检查、错误提示
    * topic缺失、重复
    * proto namespace检查
    * language namespace检查
  * [X] CI
  * [ ] 消息大小限制、检查
  * [ ] 统计整合
  * [ ] example / guide
* [ ] 集成Prometheus
* kafka
  * 自动控制partition？
* php
  * [X] client / producer
  * [X] handler
  * [X] 支持7.X
  * [X] 能否支持 5.x？()
  * [ ] 使用context传递key？
  * [x] go并发？
  * [x] 控制回收？
* go
  * [x] worker pool
* 测试
  * 大量fetch，但不commit
  * 多个group
  * 乱序commit: https://zhuanlan.zhihu.com/p/27408881
