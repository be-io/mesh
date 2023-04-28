# 隐私计算传输协议 -- PTP/1.0

## 说明

本文档为隐私计算传输规范，目前主要面向金融行业，并随时针对各行业针对互联互通提出的讨论和建议做优化。

## 版权

金融行业互联互通由银联发起，各个隐私计算技术服务商各自实现该标准，规范版权由银联及金融行业所有。

## 摘要

隐私计算传输协议是一个应用层协议，是无状态/分布式/通用的一个信息传输协议。

## 内容

* <a href="#1">1 外部通信协议</a>
    * <a href="#1.1">1.1 报文规范</a>
    * <a href="#1.2">1.2 接口规范</a>
* <a href="#2">2 内部通信协议</a>
    * <a href="#2.1">2.1 报文规范</a>
    * <a href="#2.2">2.2 接口规范</a>
* <a href="#3">3 传输SPI接口</a>
    * <a href="#3.1">3.1 管控流SPI</a>
    * <a href="#3.2">3.2 数据流SPI</a>
* <a href="#4">4 标准实现</a>
    * <a href="#4.1">4.1 运行时通信组件</a>
    * <a href="#4.2">4.1 Go</a>
    * <a href="#4.3">4.2 Java</a>
    * <a href="#4.4">4.3 Python</a>
    * <a href="#4.5">4.4 Rust</a>
    * <a href="#4.6">4.5 Typescript</a>

### <a id="1">1 外部通信协议</a>

外部通信协议描述隐私计算互联互通中跨节点间的通信的报文和接口规范，隐私计算互联互通场景下，互联互通的协议不同技术服务提供商会有自身的选择偏好，
本规范不限制各厂商使用的通信协议和业务报文序列化协议，但约定通信的报文规范和在各通信协议下的接口规范。

#### <a id="1.1">1.1 报文规范</a>

互联互通外部传输协议报文包含报头和报文两部分，报头主要提供统一的身份路由，报文Codec SPI默认使用protobuf序列化，和社区及大多数厂商保持一致。

1. 互联互通下如果使用异步传输作为通信协议，Header会复用metadata传输互联互通协议报头，且metadata中会传输异步场景下的消息相关属性

2. 互联互通如果使用HTTP系（HTTP/0.9 HTTP/1.1 HTTP/2 GRPC）作为通信协议，Header会复用HTTP报头传输互联互通协议报头

3. 互联互通如果使用其他协议作为通信协议，Header会复用metadata传输互联互通协议报头

```protobuf
syntax = "proto3";

package io.inc.ptp;

// PTP Private transfer protocol
// 通用报头名称编码，4层无Header以二进制填充到报头，7层以Header传输
enum Header {
  Version = 0;           // 协议版本               对应7层协议头x-ptp-version
  TechProviderCode = 1;  // 厂商编码               对应7层协议头x-ptp-tech-provider-code
  TraceID = 4;           // 链路追踪ID             对应7层协议头x-ptp-trace-id
  Token = 5;             // 认证令牌               对应7层协议头x-ptp-token
  SourceNodeID = 6;      // 发送端节点编号          对应7层协议头x-ptp-source-node-id
  TargetNodeID = 7;      // 接收端节点编号          对应7层协议头x-ptp-target-node-id
  SourceInstID = 8;      // 发送端机构编号          对应7层协议头x-ptp-source-inst-id
  TargetInstID = 9;      // 接收端机构编号          对应7层协议头x-ptp-target-inst-id
  SessionID = 10;        // 通信会话号，全网唯一     对应7层协议头x-ptp-session-id
}

// 通信扩展元数据编码，扩展信息均在metadata扩展
enum Metadata {
  MessageTopic = 0;                    // 消息话题，异步场景
  MessageCode = 1;                     // 消息编码，异步场景
  SourceComponentName = 2;             // 源组件名称
  TargetComponentName = 3;             // 目标组件名称
  TargetMethod = 4;                    // 目标方法
  MessageOffSet = 5;                   // 消息序列号
  InstanceId = 6;                      // 实例ID
  Timestamp = 7;                      // 时间戳
}

// 通信传输层输入报文编码
message Inbound {
  map<string, string>  metadata = 1;   // 报头，可选，预留扩展，Dict，序列化协议由通信层统一实现
  bytes payload = 2;                   // 报文，上层通信内容承载，序列化协议由上层基于SPI可插拔
}

// 通信传输层输出报文编码
message Outbound {
  map<string, string>  metadata = 1;  // 报头，可选，预留扩展，Dict，序列化协议由通信层统一实现
  bytes payload = 2;                  // 报文，上层通信内容承载，序列化协议由上层基于SPI可插拔
  string code = 3;                    // 状态码
  string message = 4;                 // 状态说明
}

// 异步消息
message Message {
  string msgId = 1; //消息ID
  bytes head = 2;   //消息头部
  bytes body = 3;   //消息体
}

// 异步消息元数据
message  TopicInfo {
  string topic = 1;
  string ip = 2;
  int32  port = 3;
  int64  createTimestamp = 4;
  int32  status = 5;
}

```

#### <a id="1.2">1.2 接口规范</a>

##### GRPC协议接口规范

约定GRPC协议要支持互联互通通信传输的标准接口。

```protobuf
service PrivateTransferProtocol {
  // 流式传输
  rpc transport (stream Inbound) returns (stream Outbound);
  // 同步传输
  rpc invoke (Inbound) returns (Outbound);
}
```

##### HTTP协议接口规范

约定HTTP协议要支持互联互通通信传输的标准接口。

```bash
# 同步传输 使用短链接
HTTP/1.1 POST io/inc/ptp/invoke
请求头:
x-ptp-version            协议版本
x-ptp-tech-provider-code 厂商编码
x-ptp-trace-id           链路追踪ID
x-ptp-token              认证令牌
x-ptp-source-node-id     发送端节点编号
x-ptp-target-node-id     接收端节点编号
x-ptp-source-inst-id     发送端机构编号
x-ptp-target-inst-id     接收端机构编号
x-ptp-session-id         通信会话号 全网唯一
请求体:
metadata                 报头，可选，预留扩展，Dict，序列化协议由通信层统一实现
payload                  由业务层序列化为二进制报文，通信协议不约束业务报文序列化协议

响应体:
metadata                 报头，可选，预留扩展，Dict，序列化协议由通信层统一实现
payload                  由业务层序列化为二进制报文，通信协议不约束业务报文序列化协议
code                     状态码
message                  状态说明

// 流式传输 使用长链接实现
HTTP/1.1 POST io/inc/ptp/transport
请求头:
x-ptp-version            协议版本
x-ptp-tech-provider-code 厂商编码
x-ptp-trace-id           链路追踪ID
x-ptp-token              认证令牌
x-ptp-source-node-id     发送端节点编号
x-ptp-target-node-id     接收端节点编号
x-ptp-source-inst-id     发送端机构编号
x-ptp-target-inst-id     接收端机构编号
x-ptp-session-id         通信会话号 全网唯一
请求体:
metadata                 报头，可选，预留扩展，Dict，序列化协议由通信层统一实现
payload                  由业务层序列化为二进制报文，通信协议不约束业务报文序列化协议

响应体:
metadata                 报头，可选，预留扩展，Dict，序列化协议由通信层统一实现
payload                  由业务层序列化为二进制报文，通信协议不约束业务报文序列化协议
code                     状态码
message                  状态说明
```

HTTP协议接口规范

### <a id="2">2 内部通信协议</a>

内部通信协议描述节点内部通信标准，由通信组件自身实现，用于与SDK或者二方应用的通信。内部通信主要通过SOA进行承载。

### <a id="2.1">2.1 报文规范</a>

```python
from typing import Any, Sequence

from mesh.macro import index, serializable


@serializable
class Inbound:

  @index(0)
  def arguments(self) -> Sequence[Any]:
    return []

  @index(1)
  def attachments(self) -> {}:
    return {}

from typing import Any

from mesh.macro import index, serializable


@serializable
class Outbound:

  @index(0)
  def code(self) -> str:
    return ''

  @index(1)
  def message(self) -> str:
    return ''

  @index(2)
  def cause(self) -> str:
    return ''

  @index(3)
  def content(self) -> Any:
    return None
```

### <a id="2.2">2.2 接口规范</a>

节点内部通信默认使用GRPC。

```protobuf
service PrivateTransferProtocol {
  // 流式传输
  rpc transport (stream Inbound) returns (stream Outbound);
  // 同步传输
  rpc invoke (Inbound) returns (Outbound);
}
```

### <a id="3">3 传输SPI接口</a>

### <a id="3.1">3.1 管控流SPI</a>

管控流SPI为管控流提供服务化托管标准接口，默认以第三代service mesh架构实现。

```python

from abc import abstractmethod, ABC
from concurrent.futures import Future

from mesh.kinds import Reference
from mesh.macro import spi
from mesh.mpc.invoker import Execution

@spi(name="grpc")
class Consumer(ABC):
    HTTP = "http"
    GRPC = "grpc"
    TCP = "tcp"
    MQTT = "mqtt"

    """
    Service consumer with any protocol and codec.
    """

    @abstractmethod
    def __enter__(self):
        pass

    @abstractmethod
    def __exit__(self, exc_type, exc_val, exc_tb):
        pass

    @abstractmethod
    def start(self):
        """
        Start the mesh broker.
        :return:
        """
        pass

    @abstractmethod
    def consume(self, address: str, urn: str, execution: Execution[Reference], inbound: bytes) -> Future:
        """
        Consume the input payload.
        :param address: Remote address.
        :param urn: Actual uniform resource domain name.
        :param execution: Service reference.
        :param inbound: Input arguments.
        :return: Output payload
        """
        pass

from abc import abstractmethod
from typing import Any

from mesh.macro import spi


@spi(name="grpc")
class Provider:
  """

  """

  @abstractmethod
  def start(self, address: str, tc: Any):
    """
    Start the mesh broker.
    :param address:
    :param tc:
    :return:
    """
    pass

  @abstractmethod
  def close(self):
    pass

```

### <a id="3.2">3.2 数据流SPI</a>

数据流SPI为数据流传输提供SPI标准接口，默认提供两个版本实现，一个基于GRPC客户端全双工模式，一个基于分布式分组高速管道实现。

```python
from abc import abstractmethod, ABC
from typing import Generic, Dict

from mesh.macro import spi, mpi, T


@spi("mesh")
class Transport(ABC, Generic[T]):
    """
    Private compute data channel in async and blocking mode.
    """

    MESH = "mesh"
    GRPC = "grpc"

    @abstractmethod
    @mpi("mesh.chan.open")
    def open(self, session_id: str, metadata: Dict[str, str]) -> "Session":
        """
        Open a channel session.
        :param session_id:  node id or inst id
        :param metadata channel metadata
        :return:
        """
        pass

    @abstractmethod
    @mpi("mesh.chan.close")
    def close(self, timeout: int):
        """
        Close the channel.
        :return:
        """
        pass


@spi("mesh")
class Session(ABC, Generic[T]):
    """
    Remote queue in async and blocking mode.
    """

    @abstractmethod
    @mpi("mesh.chan.peek")
    def peek(self, topic: str = "") -> bytes:
        """
        Retrieves, but does not remove, the head of this queue,
        or returns None if this queue is empty.
        :param topic: message topic
        :return: the head of this queue, or None if this queue is empty
        :return:
        """
        pass

    @abstractmethod
    @mpi(name="mesh.chan.pop", timeout=120 * 1000)
    def pop(self, timeout: int, topic: str = "") -> bytes:
        """
        Retrieves and removes the head of this queue,
        or returns None if this queue is empty.
        :param timeout: timeout in mills.
        :param topic: message topic
        :return: the head of this queue, or None if this queue is empty
        """
        pass

    @abstractmethod
    @mpi("mesh.chan.push")
    def push(self, payload: bytes, metadata: Dict[str, str], topic: str = ""):
        """
        Inserts the specified element into this queue if it is possible to do
        so immediately without violating capacity restrictions.
        When using a capacity-restricted queue, this method is generally
        preferable to add, which can fail to insert an element only
        by throwing an exception.
        :param payload: message payload
        :param metadata: Message metadata
        :param topic: message topic
        :return:
        """
        pass

    @abstractmethod
    @mpi("mesh.chan.release")
    def release(self, timeout: int, topic: str):
        """
        Close the channel session.
        :param timeout:
        :param topic: message topic
        :return:
        """
        pass
```

### <a id="4">4 标准实现库</a>

隐私计算互联互通传输协议标准实现为开源产品mesh，该库主要由运行时组件和多语言SDK组成，运行时提供通信的管控和传输以及适配能力，
同时提供规模化部署以及多活等高可用能力，运行时提供方式有OCI容器以及多架构可执行文件等多种方式，支持当下的所有国际和国产化CPU和OS。

### <a id="4.1">4.1 运行时</a>

运行时是互联互通通信的运行时组件，该组件提供跨节点间分布式高可用的通信管控和传输，[项目地址](https://github/incio/imesh)

```bash
docker pull imesh
docker run -d -e MID=xxx -p 7304:7304 imesh
```

### <a id="4.2">4.1 Go</a>

Go客户端提供Golang的通信的SPI集成包，通过该包可以启用通信SPI实现，进行跨节点间的数据传输和调用

```bash
go get github.com/incio/imesh/golang/client
```

### <a id="4.3">4.2 Java</a>

Java客户端提供Java的通信的SPI集成包，通过该包可以启用通信SPI实现，进行跨节点间的数据传输和调用

```xml

<dependency>
    <groupId>io.inc.mesh</groupId>
    <artifactId>imesh</artifactId>
    <version>0.0.1</version>
</dependency>
```

### <a id="4.4">4.3 Python</a>

Python客户端提供Python的通信的SPI集成包，通过该包可以启用通信SPI实现，进行跨节点间的数据传输和调用

```bash
poetry add imesh
```

### <a id="4.5">4.4 Rust</a>

Rust客户端提供Rust的通信的SPI集成包，通过该包可以启用通信SPI实现，进行跨节点间的数据传输和调用

```bash
cargo add imesh
```

### <a id="4.6">4.5 Typescript</a>

Rust客户端提供Typescript的通信的SPI集成包，通过该包可以启用通信SPI实现，进行跨节点间的数据传输和调用

```bash
yarn add imesh
```
