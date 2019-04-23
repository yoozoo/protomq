# PHP消费者

## 原理
go程序通过roadrunner管理着一个php的进程池。
当go程序从kafka获取到消息的时候，把消息发送给进程池中的一个php进程来处理。

## 使用
### 依赖
运行以下composer命令
```bat
composer require yoozoo/protomq dev-master
```

### 用法
我们要创建Consumer类的实例，注册一个消息handler，然后运行起来。下面是一个例子。
```php
<?php
include "vendor/autoload.php";

use Yoozoo\Protomq;

$consumer = new Protomq\SimpleConsumer();
$consumer->register_handler(function ($topic, $msg) {
    // 业务逻辑在这里
    return;
});
$consumer->run();
```

## 运行

### 运行go服务
```bat
protomq consumerd --brokers=localhost:9092,localhost:9093 --group=php --workers=5 demo_topic ./consumer.php
```
* `brokers`是kafka集群的地址，用逗号隔开，默认值为localhost:9092。
* `workers`是php worker的数量，默认是5.
* `group`是consumer group的标识，默认是php.
* 第一个参数是本程序想要监听的topic，用逗号隔开。
* 第二个参数是要运行的php脚本的地址。

## 完整示例
假设我们需要消费demo1,demo2这两个topic.

运行如下命令
```bat
    composer require yoozoo/protomq
```
创建demo.php文件，内容如下
```php
<?php
include "vendor/autoload.php";

use Yoozoo\Protomq;

$consumer = new Protomq\SimpleConsumer();
$consumer->register_handler(function ($topic, $msg) {
    switch ($topic) {
        case "demo1":
            // 业务逻辑1
            break;
        case "demo2":
            // 业务逻辑2
            break;
    }
    return;
});
$consumer->run();

```

运行如下命令启动go程序
```bat
protomq consumerd --brokers=localhost:9092 --topics=demo1,demo2 --workers=5 ./demo.php
```
## 设定offset
consumerd 命令支持设定partition/offset, 格式为partition0:offset0,partition1:offset1,partition2:offset2

### 示例
```
protomq consumerd --brokers=localhost:9092 --topics=demo1,demo2 --workers=5 --group=demo --offset=0:2,1:3 ./demo.php
```
上述命令会在运行时将demo这个consumer group在partition 0上的offset设为2，在partition 1上的offset设为3.

注意，offset设置只有在当前consumer对应到相应partition时才会生效，如果同时运行多个protomq consumerd命令，建议为每个命令都设置完整的offset.

注意，如果不设定group的话，会使用默认group 'php'