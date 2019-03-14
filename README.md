# protomq

[![Build Status](https://travis-ci.org/yoozoo/protomq.svg?branch=master)](https://travis-ci.org/yoozoo/protomq)

## 初始化

```bash
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
./protomq.exe gen --lang=go ./output_folder ./test.proto
```

生成`php 消费者服务器端`

```bash
./protomq.exe gen --lang=php ./output_folder ./test.proto
```

## TODO

* `protomq` cli
  * [x] 自动下载`protoc`
  * [x] 内嵌`protomq.proto`
  * [x] 嵌套调用`protoc`
  * [ ] protoc文件语法检查、错误提示
    * topic缺失、重复
    * proto namespace检查
    * language namespace检查
  * [ ] CI
  * 消息大小限制、检查
  * 统计整合
  * [ ] example / guide
* [ ] 集成Prometheus
* kafka
  * 自动控制partition？
* php
  * [ ] client / producer
  * [ ] handler
  * 支持7.X
  * 能否支持 5.x？
  * 使用context传递key？
  * go并发？
  * 控制回收？
* go
  * worker pool
* 测试
  * 大量fetch，但不commit
  * 多个group
  * 乱序commit: https://zhuanlan.zhihu.com/p/27408881
