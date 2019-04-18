# PHP消费者

## 原理
go程序通过roadrunner管理着一个php的进程池。
当go程序从kafka获取到消息的时候，把消息发送给进程池中的一个php进程来处理。

## 代码生成

### 命令
```bat
protomq gen --lang=phpconsumer ./output_folder ./test.proto
```
该命令会生成如下结构的文件, `PackageName`为proto文件中所记录的包名，同时也会成为生成的代码的`Namespace`.
```
.
+-- output_folder
|    +-- GPBMetadata
|    |    +-- Test.php
|    +-- PackageName
|    |    +-- Message1.php
|    |    +-- Message1Consumer.php
|    |    +-- Message2.php
|    |    +-- Message2Consumer.php
```

## 集成

### GPBMetadata目录
* 每次运行上述命令生成php代码，都会生成`GPBMetadata`文件夹。建议同一个PHP项目下只有一个GPBMetadata文件夹，把每次生成的GPBMetadata文件都放在这一目录下，它们共享同一个Namespace.
* 所有protomq生成的项目都需要用`protoc`额外生成[protomq.proto](../protomq.proto)文件的GPBMetadata. 你也可以直接把这个[已生成的文件](Protomq.php)放在对应的目录(GPBMetadata)下。


### 添加依赖和设置自动加载
*   运行以下composer命令
    ```bat
    composer require google/protobuf
    composer require spiral/roadrunner
    ```

*   编辑composer.json文件，添加如下autoload项
    ```json
    ...,
    "autoload":{
        ...,
        "psr-4": {
            ...,
            "PackageName\\":"path/to/PackageName",
            "GPBMetadata\\":"path/to/GPBMetadata"
        }
    }
    ```
    之后运行命令
    ```bash
    composer dump-autoload
    ```

### 使用生成的代码
以上文生成的代码目录为例：
* Message1.php和Message2.php是两个Message类。
* Message1Consumer.php和Message2Consumer.php是两个Consumer类。

我们要创建Consumer类的实例，注册一个消息handler，然后运行起来。下面是一个例子。
```php
<?php
include "vendor/autoload.php";

$consumer = new \PackageName\Message1Consumer();
$consumer->register_handler(function(\PackageName\Message1 $msg){
    return;
});
$consumer->run();
```


protomq也支持简单类型的消息传递。当proto文件中定义的Message只有一个子项，并且是简单类型时，handler函数会直接获得string类型的数据。
假如proto中定义的message如下
```protobuf
message DemoQueue {
    option (protomq.topic) = "demo_topic";
    string msg = 1;
}
```
PHP中会是这样的
```php
<?php
include "vendor/autoload.php";

$consumer = new \PackageName\Message1Consumer();
$consumer->register_handler(function($msg){
    // $msg 是 string
    return;
});
$consumer->run();
```

## 运行

### 运行go服务
```bat
protomq consumerd --brokers=localhost:9092,localhost:9093 --topics=demo_topic --group=php --workers=5 ./consumer.php
```
* `brokers`是kafka集群的地址，用逗号隔开，默认值为localhost:9092。
* `topics`是本程序想要监听的topic，用逗号隔开，必填。
* `workers`是php worker的数量，默认是5.
* `group`是consumer group的标识，默认是php.
* 最后的参数是要运行的php脚本的地址。

## 完整示例
假设我们有proto文件demo.proto,内容如下
```protobuf
syntax = "proto3";

import "protomq.proto";

package demopackage;

message SimpleQueue {
    option (protomq.topic) = "simple";

    string data = 1;
}

message Log {
    string msg = 1;
    int32 version = 2;
}

message TypedQueue {
    option (protomq.topic) = "typed";

    Log data = 1;
}
```
运行
```bat
protomq gen --lang=phpconsumer ./demo ./demo.proto
```
之后生成如下目录
```
.
+-- demo.proto
+-- demo
|    +-- GPBMetadata
|    |    +-- Demo.php
|    +-- Demopackage
|    |    +-- SimpleQueue.php
|    |    +-- SimpleQueueConsumer.php
|    |    +-- Log.php
|    |    +-- TypedQueue.php
|    |    +-- TypedQueueConsumer.php
```
将Protomq.php放入GPBMetadata文件夹下，之后运行在根目录运行
```bat
    composer require google/protobuf
    composer require spiral/roadrunner
```
编辑composer.json文件，加入如下内容
```json
"autoload":{
        ...,
        "psr-4": {
            ...,
            "Demopackage\\":"demo/Demopackage",
            "GPBMetadata\\":"demo/GPBMetadata"
        }
    }
```
之后运行
```bat
composer dump-autoload
```
之后在根目录创建demo1.php和demo2.php

此时项目目录如下
```
.
+-- demo.proto
+-- demo
|    +-- GPBMetadata
|    |    +-- Demo.php
|    |    +-- Protomq.php
|    +-- Demopackage
|    |    +-- SimpleQueue.php
|    |    +-- SimpleQueueConsumer.php
|    |    +-- Log.php
|    |    +-- TypedQueue.php
|    |    +-- TypedQueueConsumer.php
+-- vendor
|    +-- ...
|    +-- autoload.php
+-- composer.json
+-- composer.lock
+-- demo1.php
+-- demo2.php
```
demo1.php内容如下
```php
<?php
// 复杂类型的示例
include "vendor/autoload.php";

$consumer = new \Demopackage\TypedQueueConsumer();
$consumer->register_handler(function(\Demopackage\TypedQueue $msg){
    return;
});
$consumer->run();

```
demo2.php内容如下
```php
<?php
// 简单类型的示例
include "vendor/autoload.php";

$consumer = new \Demopackage\SimpleQueueConsumer();
$consumer->register_handler(function($msg){
    // $msg 是 string
    return;
});
$consumer->run();

```

运行如下命令启动go程序
```bat
protomq consumerd --brokers=localhost:9092 --topics=typed --workers=1 ./demo1.php
```
```bat
protomq consumerd --brokers=localhost:9092 --topics=simple --workers=1 ./demo2.php
```
