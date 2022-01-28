
## 1. 总结几种 socket 粘包的解包方式：fix length/delimiter based/length field based frame decoder。尝试举例其应用。

### 出现粘包的原因

TCP是面向连接的传输协议，TCP传输的数据是以流的形式，而流数据是没有明确的开始结尾边界，所以TCP也没办法判断哪一段流属于一条消息。  
粘包产生的主要原因：
* 发送方每次发送的数据 < socket缓冲区大小  
* 接收方读取socket缓冲区数据不够及时

### 粘包的解决方案

##### Fix Length

> 发送端和接收端规定固定大小的缓冲区，当字符长度不够时使用空字符弥补。 适用于发送接收种类少的特定信息。

示例代码： 
[fix_length](./fix_length)

运行效果： 

```
Received message from client, the message is: Hello, send fix length of package!
Received message from client, the message is: Hello, send fix length of package!
Received message from client, the message is: Hello, send fix length of package!
Received message from client, the message is: Hello, send fix length of package!
Received message from client, the message is: Hello, send fix length of package!
Received message from client, the message is: Hello, send fix length of package!
Received message from client, the message is: Hello, send fix length of package!
Received message from client, the message is: Hello, send fix length of package!
Received message from client, the message is: Hello, send fix length of package!
Received message from client, the message is: Hello, send fix length of package!
```


##### Delimiter Based

> 使用某几个特殊字符组合作为分隔符，每次读取的时候读到分隔符就停止，然后解析包。比较适合传递长度不固定、文本格式的数据，
> 分隔符可以用文本中不会出现的字符。

示例代码：
[delimiter](./delimiter)

运行效果： 

```
Received message from client, the message is: Hello, send msg split by delimiter!

Received message from client, the message is: Hello, send msg split by delimiter!

Received message from client, the message is: Hello, send msg split by delimiter!

Received message from client, the message is: Hello, send msg split by delimiter!

Received message from client, the message is: Hello, send msg split by delimiter!

Received message from client, the message is: Hello, send msg split by delimiter!

Received message from client, the message is: Hello, send msg split by delimiter!

Received message from client, the message is: Hello, send msg split by delimiter!

Received message from client, the message is: Hello, send msg split by delimiter!

Received message from client, the message is: Hello, send msg split by delimiter!
```

##### Length field based frame decoder

> 给数据包增加一个消息头，里面包含了包的长度信息。读取的时候先读取消息头（可以是固定长度，比如4个字节），
> 然后再按这个长度读取消息体。这种方式可以传递几乎所有信息。

示例代码： 
[length_field](./length_field)

运行效果： 

```
Received message from client, the message is: send message encoded by length field based frame!
Received message from client, the message is: send message encoded by length field based frame!
Received message from client, the message is: send message encoded by length field based frame!
Received message from client, the message is: send message encoded by length field based frame!
Received message from client, the message is: send message encoded by length field based frame!
Received message from client, the message is: send message encoded by length field based frame!
Received message from client, the message is: send message encoded by length field based frame!
Received message from client, the message is: send message encoded by length field based frame!
Received message from client, the message is: send message encoded by length field based frame!
Received message from client, the message is: send message encoded by length field based frame!
```



## 2. 实现一个从 socket connection 中解码出 goim 协议的解码器

##### Goim协议设计

主要包/帧方式： 

* Package Length: 包长度（4 bytes）
* Header Length: 头长度（2 bytes）
* Protocol Version: 协议版本（2 bytes）
* Operation: 操作码（4 bytes）
* Sequence: 请求序号ID （4 bytes）
* Body: 包内容（PackageLen-HeaderLen)

Operation: 

* Auth
* HeartBeat
* Message

Sequence: 

* 按请求、响应对应递增ID

编码实现： 
[goim_decoder](./goim_decoder)

运行结果：
```
The goim package decoded as:
packageLen: 36, headerLen: 16, version 1, operation: 2, sequenceId: 10000, body: Test goim Decoder...
The goim package decoded as:
packageLen: 36, headerLen: 16, version 1, operation: 2, sequenceId: 10000, body: Test goim Decoder...
The goim package decoded as:
packageLen: 36, headerLen: 16, version 1, operation: 2, sequenceId: 10000, body: Test goim Decoder...
```








