# protomq

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


## TODO

* `protomq` cli
  * 自动下载`protoc`
  * 内嵌`protomq.proto`
  * 嵌套调用`protoc`
  * protoc文件语法检查、错误提示
    * topic缺失、重复
    * proto namespace检查
    * language namespace检查
  * 消息大小限制、检查
  * 统计整合
* kafka
  * 自动控制partition？
* php
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
