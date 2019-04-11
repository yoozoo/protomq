# PHP生产者

## 原理
PHP通过生成的代码和向go程序发送topic和数据，go收到请求后将数据直接推送到kafka对应的topic.

## 代码生成

### 命令
```bat
protomq gen --lang=phpproducer ./output_folder ./test.proto
```
该命令会生成如下结构的文件, `PackageName`为proto文件中所记录的包名，同时也会成为生成的代码的`Namespace`.
```
.
+-- output_folder
|    +-- GPBMetadata
|    |    +-- Test.php
|    +-- PackageName
|    |    +-- Message1.php
|    |    +-- Message1Producer.php
|    |    +-- Message2.php
|    |    +-- Message2Producer.php
```

## 集成

### GPBMetadata目录
* 每次运行上述命令生成php代码，都会生成`GPBMetadata`文件夹。建议同一个PHP项目下只有一个GPBMetadata文件夹，把每次生成的GPBMetadata文件都放在这一目录下，它们共享同一个Namespace.
* 所有protomq生成的项目都需要用`protoc`额外生成[protomq.proto](../protomq.proto)文件的GPBMetadata. 你也可以直接把这个[已生成的文件](Protomq.php)放在对应的目录(GPBMetadata)下。


### 添加依赖和设置自动加载
*   运行以下composer命令
    ```bat
    composer require google/protobuf
    composer require spiral/goridge
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
* Message1Producer.php和Message2Producer.php是两个Producer类。

我们要创建Message类的实例，并使用Producer类的实例把Message发送出去。下面是一个例子。
```php
<?php
include "vendor/autoload.php";

// 构造函数的参数是go程序运行的地址和端口，默认是 127.0.0.1:8080
$sender = new \PackageName\Message1Producer("127.0.0.1","8080");

$data = new \PackageName\Message1([
    // Message2 是 Message1 的child
    "data1" => new \PackageName\Message2([
        "msg1" => "hello",
        "msg2" => "world",
    ]),
    "data2" => "!!",
]);
try {
    $sender->send($data);
} catch (Exception $e) {
    echo $e;
}
```


protomq也支持简单类型的消息传递。当proto文件中定义的Message只有一个子项，并且是简单类型时，可以不创建Message类，而是直接传递字符串。
假如proto中定义的message如下
```protobuf
message DemoQueue {
    option (protomq.topic) = "demo_topic";
    string msg = 1;
}
```
PHP中可以直接这样使用
```php
<?php
include "vendor/autoload.php";

$sender = new \PackageName\DemoQueueProducer();

try {
    $sender->send("this is message.");
} catch (Exception $e) {
    echo $e;
}
```

## 运行

### 运行go服务
```bat
protomq producerd --brokers=localhost:9092,localhost:9093 --port=8080
```
`brokers`是kafka集群的地址，用逗号隔开，默认值为localhost:9092。
`port`是本程序想要监听的端口，影响PHP端创建Producer类时的参数，默认值为8080。

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
protomq gen --lang=phpproducer ./demo ./demo.proto
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
|    |    +-- SimpleQueueProducer.php
|    |    +-- Log.php
|    |    +-- TypedQueue.php
|    |    +-- TypedQueueProducer.php
```
将Protomq.php放入GPBMetadata文件夹下，之后运行在根目录运行
```bat
    composer require google/protobuf
    composer require spiral/goridge
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
之后在根目录创建demo.php

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
|    |    +-- SimpleQueueProducer.php
|    |    +-- Log.php
|    |    +-- TypedQueue.php
|    |    +-- TypedQueueProducer.php
+-- vendor
|    +-- ...
|    +-- autoload.php
+-- composer.json
+-- composer.lock
+-- demo.php
```
demo.php内容如下
```php
<?php
include "vendor/autoload.php";

// 发送简单字符串
$simpleSender = new \Demopackage\SimpleQueueProducer("127.0.0.1","8080");
try {
    $simpleSender->send("this is a string");
} catch (Exception $e) {
    echo $e;
}

// 发送复杂类型
$typedSender = new \Demopackage\TypedQueueProducer("127.0.0.1","8080");
$data = new \Demopackage\TypedQueue([
    "data" => new \Demopackage\Log([
        "msg" => "hello",
        "version" => 3,
    ]),
]);
try {
    $typedSender->send($data);
} catch (Exception $e) {
    echo $e;
}
```

运行如下命令启动go程序
```bat
protomq producerd
```

之后在项目根目录运行如下命令会向kafka的simple和typed两个topic各发送一条消息。
```bat
php ./demo.php
```
