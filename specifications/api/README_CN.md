# 隐私计算传输标准使用文档

查看英文版 [PTP Standard Document](README.md)

查看规范 [PTP/1.0](Specifications_CN.md)

## 说明

本文档为隐私计算传输标准使用文档，作为隐私计算技术服务商快速集成互联互通的操作手册。

## 版权

本标准由北京金融科技产业联盟（Beijing FinTech Industry Alliance）牵头编写，帮助各个隐私计算技术服务商实现快速集成和落地，版权由北京金融科技产业联盟所有。

**开源协议**

| 协议名        | 地址                                              |
|------------|-------------------------------------------------|
| Apache 2.0 | http://www.apache.org/licenses/LICENSE-2.0.html |

**联系方式**

| 联系人 | 邮箱                       |
|-----|--------------------------|
| 王超  | congying.wang@trustbe.cn |
| 曾成  | coyzeng@gmail.com        |

## 摘要

隐私计算传输标准实现为Mesh。Mesh分为客户端和服务端两部分，客户端为各语言的集成SDK，服务端为通信模块底座。

Mesh支持两种方式集成

1. 直接使用MESH的客户端SDK
2. 直接调用Mesh的原生API。

Mesh支持多种方式运行(支持ARM/AMD架构)

1. 作为容器运行在各个容器运行时（Docker、Containerd等）
2. 作为容器运行在各类PaaS及云平台
3. 作为应用运行在各类裸机环境（Windows、Linux、Darwin）

![logo](https://github.com/be-io/mesh/blob/master/specifications/api/pipline.png?raw=true)

## 内容

* <a href="#1">1 安装使用</a>
    * <a href="#1.1">1.1 Docker部署</a>
    * <a href="#1.2">1.2 K8S部署</a>
    * <a href="#1.3">1.3 裸系统部署</a>
    * <a href="#1.4">1.4 客户端使用</a>
* <a href="#2">2 环境初始化</a>
    * <a href="#2.1">2.1 设置节点信息</a>
    * <a href="#2.2">2.2 节点组网申请</a>
    * <a href="#2.3">2.3 节点组网授权</a>
* <a href="#3">3 数据面通信</a>
    * <a href="#3.1">3.1 开启一个传输通道</a>
    * <a href="#3.2">3.2 发送报文</a>
    * <a href="#3.3">3.3 接受报文</a>
    * <a href="#3.4">3.4 释放一个传输通道</a>
* <a href="#4">4 管理面通信</a>
    * <a href="#4.1">4.1 动态路由注册</a>
* <a href="#5">5 附录</a>

### <a id="1">1 安装部署</a>

Mesh支持以容器或可执行文件两种形式进行部署。Mesh通过一个端口路由所有通信协议，因此部署Mesh需要映射该端口，端口默认为7304。

若想修改Mesh监听的默认端口，可以通过环境变量`PLUGIN_PROXY_TRANSPORT_Y`进行修改，如`PLUGIN_PROXY_TRANSPORT_Y=17304`。

#### <a id="1.1">1.1 Docker部署</a>

使用Docker运行Mesh容器命令如下，该命令会启动容器，并映射7304端口提供服务。

```shell
docker pull bfia/mesh
docker run -d -p 7304:7304 bfia/mesh
```

#### <a id="1.2">1.2 K8S部署</a>

使用K8S运行Mesh容器命令及配置如下，该命令会创建Pod，并映射7304端口提供服务。

K8S部署描述文件，保存为`mesh.yaml`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: mesh
    group: bfia
  name: mesh
spec:
  progressDeadlineSeconds: 600
  replicas: 1
  revisionHistoryLimit: 10
  selector:
    matchLabels:
      app: mesh
      group: bfia
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: mesh
        group: bfia
    spec:
      containers:
        - name: main
          image: bfia/mesh
          imagePullPolicy: Always
          env:
            - name: MESH_HOME
              value: /mesh
          livenessProbe:
            failureThreshold: 24
            httpGet:
              path: /stats
              port: 7304
              scheme: HTTP
            initialDelaySeconds: 30
            periodSeconds: 5
            successThreshold: 1
            timeoutSeconds: 5
          readinessProbe:
            failureThreshold: 24
            httpGet:
              path: /stats
              port: 7304
              scheme: HTTP
            initialDelaySeconds: 30
            periodSeconds: 5
            successThreshold: 1
            timeoutSeconds: 5
          terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
      dnsPolicy: ClusterFirst
      restartPolicy: Always
      schedulerName: default-scheduler
      securityContext: { }
      terminationGracePeriodSeconds: 30
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: mesh
    group: bfia
  name: mesh
spec:
  ipFamilies:
    - IPv4
  ipFamilyPolicy: SingleStack
  ports:
    - name: "7304"
      port: 7304
      protocol: TCP
      targetPort: 7304
  selector:
    app: mesh
    group: bfia
  sessionAffinity: None
  type: ClusterIP
```

K8S部署命令

```shell
kubectl apply -f mesh.yaml -n namespace
```

#### <a id="1.3">1.3 裸系统部署</a>

使用裸系统部署Mesh命令如下，该命令会启动容器，并映射7304端口提供服务。

```shell
curl -Ol https://github.com/be-io/mesh/releases/download/v0.0.1/mesh-0.0.1-linux-amd64.tar.gz
tar -zxvf mesh-0.0.1-linux-amd64.tar.gz
cd mesh-0.0.1
mesh start --server
```

#### <a id="1.4">1.4 客户端使用</a>

Mesh提供多语言SDK进行快速开发集成，SDK需要配置一个环境变量或启动参数`MESH_ADDRESS`，
该地址为需要链接的Mesh服务端地址，如：192.168.1.100:7304。

Go SDK提供Golang语言通信SPI集成包，实现Go语言节点间跨域通信管控和传输。

```bash
go get github.com/be-io/mesh/client/golang
```

Java SDK提供Golang语言通信SPI集成包，实现Java语言节点间跨域通信管控和传输。

```xml

<dependency>
    <groupId>io.be.mesh</groupId>
    <artifactId>mesh</artifactId>
    <version>0.0.1</version>
</dependency>
```

Python SDK提供Golang语言通信SPI集成包，实现Python语言节点间跨域通信管控和传输。

```bash
poetry add imesh
```

Rust SDK提供Golang语言通信SPI集成包，实现Rust语言节点间跨域通信管控和传输。

```bash
cargo add imesh
```

Typescript SDK提供Golang语言通信SPI集成包，实现Typescript语言节点间跨域通信管控和传输。

```bash
yarn add imesh
```

### <a id="2">2 环境初始化</a>

自2023年开始，Mesh自身不内置节点信息，节点信息由管理面调用Mesh接口进行初始化，管理面调用Mesh接口设置的数据均会被持久化到`MESH_HOME`
目录下。

`MESH_HOME`默认在`~/mesh`，`MESH_HOME`在容器部署模式下，建议挂载到持久化存储券。

### <a id="2.1">2.1 设置节点信息</a>

设置节点信息由管理面调用Mesh接口进行初始化，初始化只需要一次，数据持久化在`MESH_HOME`挂载目录对应的存储介质中。

SDK调用设置节点信息接口，如Java语言

```java

import io.be.mesh.client.mpc.ServiceLoader;
import io.be.mesh.client.prisim.KMS;
import io.be.mesh.client.struct.Environ;

public class Example {

    public static void main(String[] args) {
        System.setProperty("MESH_ADDRESS", "127.0.0.1:7304");
        Environ environ = new Environ();
        environ.setVersion("1.0.0");
        environ.setNodeId("LX0000010000280");
        environ.setInstId("JG0100002800000000");
        environ.setInstName("互联互通节点X");
        // 使用上传自有根证书设置，不设置Mesh会自动签发节点独立的根证书
        environ, setRootCrt("-----BEGIN CERTIFICATE-----\nMIIEZTCCA02gAwIBAgIRAKkfFud6oGcbEmGOZ2i917wwDQYJKoZIhvcNAQELBQAw\ngZIxCzAJBgNVBAYTAkNOMQswCQYDVQQIEwJaSjELMAkGA1UEBxMCSFoxJDAiBgNV\nBAoTG0xYMDAwMDAxMDAwMDI3MC50cnVzdGJlLm5ldDEkMCIGA1UECxMbTFgwMDAw\nMDEwMDAwMjcwLnRydXN0YmUubmV0MR0wGwYDVQQDDBTok53osaHpm4blm6IyN+WF\nrOWPuDAgFw0yMjA5MjEwNjU1MTRaGA8yMTIyMDkyMTA2NTUxNFowgZIxCzAJBgNV\nBAYTAkNOMQswCQYDVQQIEwJaSjELMAkGA1UEBxMCSFoxJDAiBgNVBAoTG0xYMDAw\nMDAxMDAwMDI3MC50cnVzdGJlLm5ldDEkMCIGA1UECxMbTFgwMDAwMDEwMDAwMjcw\nLnRydXN0YmUubmV0MR0wGwYDVQQDDBTok53osaHpm4blm6IyN+WFrOWPuDCCASIw\nDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAOy3fx+vyvaRBTg9hXSXW/XEz5A9\nazavO4KNMPB/Rp6UtQ9+ilireBQXqCM1Zzo5/cZ+FoRmQp+Rl6b1HYAnRaPSTNFM\n4JDTl3AczsLwJcvFlqobqap9/3k5ZEfEEKJp0LFcl3GR2Ral7Uhsmgzm87PSZVRX\nFYPgCgfkuWGbeJ9FwCB5ZDovPH77pyJeW0AzPu3JKLk3JnCtUvmXK4HzoxChK4k0\nqJF1DEVSPxG3JLOKQKQqe/fqSTkRfgrU7D7U/TEaEl5lLFPGi6ZT+3AbOUAasWWO\nMTYMGMG6qC2wlQ7WdJJ4KiWhEfA7aGQfzoW85PZyN1QzRyzW1vyXY5MKNXsCAwEA\nAaOBsTCBrjAOBgNVHQ8BAf8EBAMCArwwJwYDVR0lBCAwHgYIKwYBBQUHAwEGCCsG\nAQUFBwMCBggrBgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBSrXjQL\nGNwhXBQxqx5Sh+LQwXK3mzBDBgNVHREEPDA6ghtMWDAwMDAwMTAwMDAyNzAudHJ1\nc3RiZS5uZXSBG0xYMDAwMDAxMDAwMDI3MEB0cnVzdGJlLm5ldDANBgkqhkiG9w0B\nAQsFAAOCAQEAcRhabDb2O55lQfeMkHozqFz3V4fQwQKQzYW6gwf7FEUxPMaxmFrk\nPKK2rE5Li49mUxk/norXMVmEpBQdYgsu2PnY5J39RivFku0nObzIMFkyDPer05tp\nx7sKXS6sIy9gfaN6uZkzHZH+0MD2NlOj19SI5BNJxiwHzB9BZwGwbKtq3DVZzWRq\nVyZnAcXVTq1Pk9gypatvBgD7r6edXSBXEz2d8HMaidJfChweZPVp98uY8s9EzHnJ\nowsia5sPaVNqVFE72IgQ5TUrRQsQIBGOF81i37Bpsc9g/pf2s6eNGW4CV9a8yfN1\n+BLRh/2eCYXUNQOm5IovVXqBJ4MlhchAOQ==\n-----END CERTIFICATE-----\n");
        environ.setRootKey("-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA7Ld/H6/K9pEFOD2FdJdb9cTPkD1rNq87go0w8H9GnpS1D36K\nWKt4FBeoIzVnOjn9xn4WhGZCn5GXpvUdgCdFo9JM0UzgkNOXcBzOwvAly8WWqhup\nqn3/eTlkR8QQomnQsVyXcZHZFqXtSGyaDObzs9JlVFcVg+AKB+S5YZt4n0XAIHlk\nOi88fvunIl5bQDM+7ckouTcmcK1S+ZcrgfOjEKEriTSokXUMRVI/Ebcks4pApCp7\n9+pJORF+CtTsPtT9MRoSXmUsU8aLplP7cBs5QBqxZY4xNgwYwbqoLbCVDtZ0kngq\nJaER8DtoZB/Ohbzk9nI3VDNHLNbW/Jdjkwo1ewIDAQABAoIBAQDAhzIu3HTQe/zp\n1CfSPzT9PLixETM9Q+K7+Qgf4vTWEA7/biUpnzTH6sHG+S1fT0FXir/XqbBwRiM5\nGM2IqOhcKLRv2v4e7OmTtup35IhpJui2rE8fquD5gLNOJ2p8HmItjyhhp4UQhZ3r\nNOFKsyDtVacypK2MF9EwwFgCykeeCbVwBj8PxmTunDuVlD/hSTUa9JK/lh3ssGfN\novPPiuL2TTPiZUy415y4Kd4dVvuQW0QIlm4yV+qMynSNgU7p8f6jyKY5R7FAG3MF\nBEO07LZSYzTOhfhd2zz6lQpi5t2srScIUs+6hQSyZlFqILwkKeIgKTYbvb1q+o+s\nk5XijuwBAoGBAP2JjcENSCUj5c1bRiOvczS/XkDrE5YIOwJMFLWjuwuMNZN++PV3\ngZ58cM4leJ2TmpVidfu9Rkcp6q95M5fZa7Ms7qtS8sFu2JZhWBpydToRrsFAApsG\nBU09cfoq7DhlqJyE+p/X0lmmKyTw2HNeWUGU5AjSPeVa/4da9CSjDn97AoGBAO8E\nHevtVODXCJC+VJXOjlvCff4zGQi3pFo5s9hCXUDk8KNmfieJwcLXlqz3MG/kbmHi\n0ywToeBK5wA/npImAFKXGR5U6ivOeTBML8dM0wG+2gVCwazxm3qO1d2jyD6d/o2H\n/pH3rmuj+Hi9QSZPVVp+nrJsrG2HtMkwGFWP+kIBAoGANtZsqafUxeu4xa0LQ6as\nNWl62nG9/8Jx+PI5vHvYdgvyfp+E+5rIl131DDGAoByP3+W2/ScYL0Y6s490gFCP\ngeajDL1ZMktmX0hYxQeioVe3w6azqZIozWcP4vsrspsSWCBPEQmePrO5Ozk4p+Nt\nTMkGdX3700LWaBFdIxt9hEcCgYAf7CvW49bPRMkHE/SWIYVP6hULy2VPjb9ssYI8\novhzf2BIYpr8yuBPFp4wMb+NYjP/7NyJaYHYRAjANr8GA/9NCJM5QtwXx7bV5YcI\nFlGkTQovY7AcWhSK9OLJfGN1QYLLAlvUwQDRrY+1CInYBQaAVKL7b5pD8rkJmdvW\nKamiAQKBgQCq3cjGNqvMBVWeq0qJRsouSXRaHFRTmS9hgrRD/GZ+GQJNuhsDwfV9\nhb74cz7EqaFhQVhLHZ2QP/AauqtZ+p9wIM5X9xJDCaEKGVR+95vvxilBIOqua5IK\nphkbWGfS/Wp0fWGDaBoc05xuMm3Zv7GdiNeKr6soT1dOOCA88L3eTQ==\n-----END RSA PRIVATE KEY-----\n");
        KMS kms = ServiceLoader.load(KMS.class).getDefault();
        kms.reset(environ);
    }

}
```

原生接口（HTTP或GRPC协议）调用设置节点信息接口

```shell
curl --location 'http://10.99.28.33:7304/v1/interconn/node/refresh' --header 'Content-Type: application/json' --data '{
    "version": "1.0.0",
    "node_id": "LX0000010000280",
    "inst_id": "JG0100002800000000",
    "inst_name": "互联互通节点X",
    "root_crt": "-----BEGIN CERTIFICATE-----\nMIIEZTCCA02gAwIBAgIRAKkfFud6oGcbEmGOZ2i917wwDQYJKoZIhvcNAQELBQAw\ngZIxCzAJBgNVBAYTAkNOMQswCQYDVQQIEwJaSjELMAkGA1UEBxMCSFoxJDAiBgNV\nBAoTG0xYMDAwMDAxMDAwMDI3MC50cnVzdGJlLm5ldDEkMCIGA1UECxMbTFgwMDAw\nMDEwMDAwMjcwLnRydXN0YmUubmV0MR0wGwYDVQQDDBTok53osaHpm4blm6IyN+WF\nrOWPuDAgFw0yMjA5MjEwNjU1MTRaGA8yMTIyMDkyMTA2NTUxNFowgZIxCzAJBgNV\nBAYTAkNOMQswCQYDVQQIEwJaSjELMAkGA1UEBxMCSFoxJDAiBgNVBAoTG0xYMDAw\nMDAxMDAwMDI3MC50cnVzdGJlLm5ldDEkMCIGA1UECxMbTFgwMDAwMDEwMDAwMjcw\nLnRydXN0YmUubmV0MR0wGwYDVQQDDBTok53osaHpm4blm6IyN+WFrOWPuDCCASIw\nDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAOy3fx+vyvaRBTg9hXSXW/XEz5A9\nazavO4KNMPB/Rp6UtQ9+ilireBQXqCM1Zzo5/cZ+FoRmQp+Rl6b1HYAnRaPSTNFM\n4JDTl3AczsLwJcvFlqobqap9/3k5ZEfEEKJp0LFcl3GR2Ral7Uhsmgzm87PSZVRX\nFYPgCgfkuWGbeJ9FwCB5ZDovPH77pyJeW0AzPu3JKLk3JnCtUvmXK4HzoxChK4k0\nqJF1DEVSPxG3JLOKQKQqe/fqSTkRfgrU7D7U/TEaEl5lLFPGi6ZT+3AbOUAasWWO\nMTYMGMG6qC2wlQ7WdJJ4KiWhEfA7aGQfzoW85PZyN1QzRyzW1vyXY5MKNXsCAwEA\nAaOBsTCBrjAOBgNVHQ8BAf8EBAMCArwwJwYDVR0lBCAwHgYIKwYBBQUHAwEGCCsG\nAQUFBwMCBggrBgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBSrXjQL\nGNwhXBQxqx5Sh+LQwXK3mzBDBgNVHREEPDA6ghtMWDAwMDAwMTAwMDAyNzAudHJ1\nc3RiZS5uZXSBG0xYMDAwMDAxMDAwMDI3MEB0cnVzdGJlLm5ldDANBgkqhkiG9w0B\nAQsFAAOCAQEAcRhabDb2O55lQfeMkHozqFz3V4fQwQKQzYW6gwf7FEUxPMaxmFrk\nPKK2rE5Li49mUxk/norXMVmEpBQdYgsu2PnY5J39RivFku0nObzIMFkyDPer05tp\nx7sKXS6sIy9gfaN6uZkzHZH+0MD2NlOj19SI5BNJxiwHzB9BZwGwbKtq3DVZzWRq\nVyZnAcXVTq1Pk9gypatvBgD7r6edXSBXEz2d8HMaidJfChweZPVp98uY8s9EzHnJ\nowsia5sPaVNqVFE72IgQ5TUrRQsQIBGOF81i37Bpsc9g/pf2s6eNGW4CV9a8yfN1\n+BLRh/2eCYXUNQOm5IovVXqBJ4MlhchAOQ==\n-----END CERTIFICATE-----\n",
    "root_key": "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA7Ld/H6/K9pEFOD2FdJdb9cTPkD1rNq87go0w8H9GnpS1D36K\nWKt4FBeoIzVnOjn9xn4WhGZCn5GXpvUdgCdFo9JM0UzgkNOXcBzOwvAly8WWqhup\nqn3/eTlkR8QQomnQsVyXcZHZFqXtSGyaDObzs9JlVFcVg+AKB+S5YZt4n0XAIHlk\nOi88fvunIl5bQDM+7ckouTcmcK1S+ZcrgfOjEKEriTSokXUMRVI/Ebcks4pApCp7\n9+pJORF+CtTsPtT9MRoSXmUsU8aLplP7cBs5QBqxZY4xNgwYwbqoLbCVDtZ0kngq\nJaER8DtoZB/Ohbzk9nI3VDNHLNbW/Jdjkwo1ewIDAQABAoIBAQDAhzIu3HTQe/zp\n1CfSPzT9PLixETM9Q+K7+Qgf4vTWEA7/biUpnzTH6sHG+S1fT0FXir/XqbBwRiM5\nGM2IqOhcKLRv2v4e7OmTtup35IhpJui2rE8fquD5gLNOJ2p8HmItjyhhp4UQhZ3r\nNOFKsyDtVacypK2MF9EwwFgCykeeCbVwBj8PxmTunDuVlD/hSTUa9JK/lh3ssGfN\novPPiuL2TTPiZUy415y4Kd4dVvuQW0QIlm4yV+qMynSNgU7p8f6jyKY5R7FAG3MF\nBEO07LZSYzTOhfhd2zz6lQpi5t2srScIUs+6hQSyZlFqILwkKeIgKTYbvb1q+o+s\nk5XijuwBAoGBAP2JjcENSCUj5c1bRiOvczS/XkDrE5YIOwJMFLWjuwuMNZN++PV3\ngZ58cM4leJ2TmpVidfu9Rkcp6q95M5fZa7Ms7qtS8sFu2JZhWBpydToRrsFAApsG\nBU09cfoq7DhlqJyE+p/X0lmmKyTw2HNeWUGU5AjSPeVa/4da9CSjDn97AoGBAO8E\nHevtVODXCJC+VJXOjlvCff4zGQi3pFo5s9hCXUDk8KNmfieJwcLXlqz3MG/kbmHi\n0ywToeBK5wA/npImAFKXGR5U6ivOeTBML8dM0wG+2gVCwazxm3qO1d2jyD6d/o2H\n/pH3rmuj+Hi9QSZPVVp+nrJsrG2HtMkwGFWP+kIBAoGANtZsqafUxeu4xa0LQ6as\nNWl62nG9/8Jx+PI5vHvYdgvyfp+E+5rIl131DDGAoByP3+W2/ScYL0Y6s490gFCP\ngeajDL1ZMktmX0hYxQeioVe3w6azqZIozWcP4vsrspsSWCBPEQmePrO5Ozk4p+Nt\nTMkGdX3700LWaBFdIxt9hEcCgYAf7CvW49bPRMkHE/SWIYVP6hULy2VPjb9ssYI8\novhzf2BIYpr8yuBPFp4wMb+NYjP/7NyJaYHYRAjANr8GA/9NCJM5QtwXx7bV5YcI\nFlGkTQovY7AcWhSK9OLJfGN1QYLLAlvUwQDRrY+1CInYBQaAVKL7b5pD8rkJmdvW\nKamiAQKBgQCq3cjGNqvMBVWeq0qJRsouSXRaHFRTmS9hgrRD/GZ+GQJNuhsDwfV9\nhb74cz7EqaFhQVhLHZ2QP/AauqtZ+p9wIM5X9xJDCaEKGVR+95vvxilBIOqua5IK\nphkbWGfS/Wp0fWGDaBoc05xuMm3Zv7GdiNeKr6soT1dOOCA88L3eTQ==\n-----END RSA PRIVATE KEY-----\n"
}'
```

### <a id="2.2">2.2 节点组网申请</a>

节点组网申请由管理面调用Mesh接口进行组网申请，数据持久化在`MESH_HOME`挂载目录对应的存储介质中。

SDK调用节点组网申请接口，如Java语言

```java
import io.be.mesh.client.mpc.ServiceLoader;
import io.be.mesh.client.prisim.Network;
import io.be.mesh.client.struct.Route;

public class Example {

    public static void main(String[] args) {
        System.setProperty("MESH_ADDRESS", "127.0.0.1:7304");
        Route route = new Route();
        route.setNodeId("LX0000010000280");
        route.setInstId("JG0100002800000000");
        route.setAddress("192.168.100.63:27304");
        // 如需上传已有的双向认证证书设置，不设置Mesh会自动签署点对点双向认证证书
        route, setGuestRoot("-----BEGIN CERTIFICATE-----\nMIIEZTCCA02gAwIBAgIRAKkfFud6oGcbEmGOZ2i917wwDQYJKoZIhvcNAQELBQAw\ngZIxCzAJBgNVBAYTAkNOMQswCQYDVQQIEwJaSjELMAkGA1UEBxMCSFoxJDAiBgNV\nBAoTG0xYMDAwMDAxMDAwMDI3MC50cnVzdGJlLm5ldDEkMCIGA1UECxMbTFgwMDAw\nMDEwMDAwMjcwLnRydXN0YmUubmV0MR0wGwYDVQQDDBTok53osaHpm4blm6IyN+WF\nrOWPuDAgFw0yMjA5MjEwNjU1MTRaGA8yMTIyMDkyMTA2NTUxNFowgZIxCzAJBgNV\nBAYTAkNOMQswCQYDVQQIEwJaSjELMAkGA1UEBxMCSFoxJDAiBgNVBAoTG0xYMDAw\nMDAxMDAwMDI3MC50cnVzdGJlLm5ldDEkMCIGA1UECxMbTFgwMDAwMDEwMDAwMjcw\nLnRydXN0YmUubmV0MR0wGwYDVQQDDBTok53osaHpm4blm6IyN+WFrOWPuDCCASIw\nDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAOy3fx+vyvaRBTg9hXSXW/XEz5A9\nazavO4KNMPB/Rp6UtQ9+ilireBQXqCM1Zzo5/cZ+FoRmQp+Rl6b1HYAnRaPSTNFM\n4JDTl3AczsLwJcvFlqobqap9/3k5ZEfEEKJp0LFcl3GR2Ral7Uhsmgzm87PSZVRX\nFYPgCgfkuWGbeJ9FwCB5ZDovPH77pyJeW0AzPu3JKLk3JnCtUvmXK4HzoxChK4k0\nqJF1DEVSPxG3JLOKQKQqe/fqSTkRfgrU7D7U/TEaEl5lLFPGi6ZT+3AbOUAasWWO\nMTYMGMG6qC2wlQ7WdJJ4KiWhEfA7aGQfzoW85PZyN1QzRyzW1vyXY5MKNXsCAwEA\nAaOBsTCBrjAOBgNVHQ8BAf8EBAMCArwwJwYDVR0lBCAwHgYIKwYBBQUHAwEGCCsG\nAQUFBwMCBggrBgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBSrXjQL\nGNwhXBQxqx5Sh+LQwXK3mzBDBgNVHREEPDA6ghtMWDAwMDAwMTAwMDAyNzAudHJ1\nc3RiZS5uZXSBG0xYMDAwMDAxMDAwMDI3MEB0cnVzdGJlLm5ldDANBgkqhkiG9w0B\nAQsFAAOCAQEAcRhabDb2O55lQfeMkHozqFz3V4fQwQKQzYW6gwf7FEUxPMaxmFrk\nPKK2rE5Li49mUxk/norXMVmEpBQdYgsu2PnY5J39RivFku0nObzIMFkyDPer05tp\nx7sKXS6sIy9gfaN6uZkzHZH+0MD2NlOj19SI5BNJxiwHzB9BZwGwbKtq3DVZzWRq\nVyZnAcXVTq1Pk9gypatvBgD7r6edXSBXEz2d8HMaidJfChweZPVp98uY8s9EzHnJ\nowsia5sPaVNqVFE72IgQ5TUrRQsQIBGOF81i37Bpsc9g/pf2s6eNGW4CV9a8yfN1\n+BLRh/2eCYXUNQOm5IovVXqBJ4MlhchAOQ==\n-----END CERTIFICATE-----\n");
        route.setGuestCrt("-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA7Ld/H6/K9pEFOD2FdJdb9cTPkD1rNq87go0w8H9GnpS1D36K\nWKt4FBeoIzVnOjn9xn4WhGZCn5GXpvUdgCdFo9JM0UzgkNOXcBzOwvAly8WWqhup\nqn3/eTlkR8QQomnQsVyXcZHZFqXtSGyaDObzs9JlVFcVg+AKB+S5YZt4n0XAIHlk\nOi88fvunIl5bQDM+7ckouTcmcK1S+ZcrgfOjEKEriTSokXUMRVI/Ebcks4pApCp7\n9+pJORF+CtTsPtT9MRoSXmUsU8aLplP7cBs5QBqxZY4xNgwYwbqoLbCVDtZ0kngq\nJaER8DtoZB/Ohbzk9nI3VDNHLNbW/Jdjkwo1ewIDAQABAoIBAQDAhzIu3HTQe/zp\n1CfSPzT9PLixETM9Q+K7+Qgf4vTWEA7/biUpnzTH6sHG+S1fT0FXir/XqbBwRiM5\nGM2IqOhcKLRv2v4e7OmTtup35IhpJui2rE8fquD5gLNOJ2p8HmItjyhhp4UQhZ3r\nNOFKsyDtVacypK2MF9EwwFgCykeeCbVwBj8PxmTunDuVlD/hSTUa9JK/lh3ssGfN\novPPiuL2TTPiZUy415y4Kd4dVvuQW0QIlm4yV+qMynSNgU7p8f6jyKY5R7FAG3MF\nBEO07LZSYzTOhfhd2zz6lQpi5t2srScIUs+6hQSyZlFqILwkKeIgKTYbvb1q+o+s\nk5XijuwBAoGBAP2JjcENSCUj5c1bRiOvczS/XkDrE5YIOwJMFLWjuwuMNZN++PV3\ngZ58cM4leJ2TmpVidfu9Rkcp6q95M5fZa7Ms7qtS8sFu2JZhWBpydToRrsFAApsG\nBU09cfoq7DhlqJyE+p/X0lmmKyTw2HNeWUGU5AjSPeVa/4da9CSjDn97AoGBAO8E\nHevtVODXCJC+VJXOjlvCff4zGQi3pFo5s9hCXUDk8KNmfieJwcLXlqz3MG/kbmHi\n0ywToeBK5wA/npImAFKXGR5U6ivOeTBML8dM0wG+2gVCwazxm3qO1d2jyD6d/o2H\n/pH3rmuj+Hi9QSZPVVp+nrJsrG2HtMkwGFWP+kIBAoGANtZsqafUxeu4xa0LQ6as\nNWl62nG9/8Jx+PI5vHvYdgvyfp+E+5rIl131DDGAoByP3+W2/ScYL0Y6s490gFCP\ngeajDL1ZMktmX0hYxQeioVe3w6azqZIozWcP4vsrspsSWCBPEQmePrO5Ozk4p+Nt\nTMkGdX3700LWaBFdIxt9hEcCgYAf7CvW49bPRMkHE/SWIYVP6hULy2VPjb9ssYI8\novhzf2BIYpr8yuBPFp4wMb+NYjP/7NyJaYHYRAjANr8GA/9NCJM5QtwXx7bV5YcI\nFlGkTQovY7AcWhSK9OLJfGN1QYLLAlvUwQDRrY+1CInYBQaAVKL7b5pD8rkJmdvW\nKamiAQKBgQCq3cjGNqvMBVWeq0qJRsouSXRaHFRTmS9hgrRD/GZ+GQJNuhsDwfV9\nhb74cz7EqaFhQVhLHZ2QP/AauqtZ+p9wIM5X9xJDCaEKGVR+95vvxilBIOqua5IK\nphkbWGfS/Wp0fWGDaBoc05xuMm3Zv7GdiNeKr6soT1dOOCA88L3eTQ==\n-----END RSA PRIVATE KEY-----\n");
        Network network = ServiceLoader.load(Network.class).getDefault();
        network.refresh(route);
        network.weave(route);
    }

}
```

原生接口（HTTP或GRPC协议）调用节点组网申请

```shell
curl --location 'http://10.99.27.33:7304/v1/interconn/net/refresh' --header 'Content-Type: application/json' --data '{
    "node_id": "LX0000010000280",
    "inst_id": "JG0100002800000000",
    "address": "192.168.100.63:27304",
    "guest_root": "",
    "guest_crt": ""
}'
```

### <a id="2.3">2.3 节点组网授权</a>

节点组网授权目前未做强制要求。

SDK调用节点组网授权接口，如Java语言

```java
import io.be.mesh.client.mpc.ServiceLoader;
import io.be.mesh.client.prisim.Network;
import io.be.mesh.client.struct.Route;

public class Example {

    public static void main(String[] args) {
        System.setProperty("MESH_ADDRESS", "127.0.0.1:7304");
        Route route = new Route();
        route.setNodeId("LX0000010000280");
        route.setInstId("JG0100002800000000");
        Network network = ServiceLoader.load(Network.class).getDefault();
        network.ack(route);
    }

}
```

原生接口（HTTP或GRPC协议）调用节点组网授权

```shell
curl --location 'http://10.99.27.33:7304/v1/interconn/net/ack' --header 'Content-Type: application/json' --data '{
    "node_id": "LX0000010000280",
    "inst_id": "JG0100002800000000",
}'
```

### <a id="3">3 管理面通信</a>

Mesh除提供数据面通信外，还额外提供管理面通信能力。管理面通信以两种方式提供

1. 动态服务路由注册，该模式仅提供管理面通信路由能力
2. 管理面服务化托管，该模式提供多语言接口全面服务化托管和治理

目前互联互通管理面只用到第一种模式。

### <a id="3.1">3.1 动态路由注册</a>

动态路由注册提供服务的路由规则动态注册能力，为管理面提供灵活、收敛、统一的通信路由能力。可以理解成提供一个加强版本Nginx的能力。
管理面通过该能力，实现所有隐私计算通信协议一个端口。 目前算子服务也是通过该能力实现分布式服务化。

SDK调用动态路由注册，如Python语言

```python
import mesh
from mesh import Mesh, log, Registry, MethodProxy, ServiceLoader, Metadata, Routable, Transport
from mesh.kinds import Registration, Service, Resource


def main():
    mesh.start()

    http1 = Service()
    http1.urn = '/a/b/c'
    http1.kind = 'Restful'
    http1.namespace = ''
    http1.name = ''
    http1.version = '1.0.0'
    http1.proto = 'http'
    http1.codec = 'json'
    http1.asyncable = False
    http1.timeout = 10000
    http1.retries = 0
    http1.address = "IP:PORT"
    http1.lang = 'Python3'
    http1.attrs = {}

    grpc = Service()
    grpc.urn = '/a/b/c'
    grpc.kind = 'Restful'
    grpc.namespace = ''
    grpc.name = ''
    grpc.version = '1.0.0'
    grpc.proto = 'grpc'
    grpc.codec = 'json'
    grpc.asyncable = False
    grpc.timeout = 10000
    grpc.retries = 0
    grpc.address = "IP:PORT"
    grpc.lang = 'Python3'
    grpc.attrs = {}

    metadata = Resource()
    metadata.services = [http1, grpc]

    registration = Registration()
    registration.name = 'tensorserving'
    registration.kind = "complex"
    registration.address = "IP:PORT"
    registration.instance_id = "IP:PORT"
    registration.timestamp = 5 * 60 * 1000
    registration.attachments = {}
    registration.content = metadata

    registry = ServiceLoader.load(MethodProxy).get_default().proxy(Registry)
    Routable.of(registry).with_address("10.99.27.33:570").local().register(registration)
```

原生接口（HTTP或GRPC协议）调用动态路由注册

```shell
curl --location 'http://10.99.27.33:7304/v1/interconn/registry/register' --header 'Content-Type: application/json' --data '{
    "name": "tensorserving",
    "kind": "complex",
    "address": "IP:PORT",
    "instance_id": "IP:PORT",
    "timestamp": 300000,
    "attachments": {},
    "content": {
    "services": [
            {
                "urn": "/a/b/c",
                "kind": "Restful",
                "namespace": "",
                "name": "",
                "version": "1.0.0",
                "proto": "http",
                "codec": "json",
                "asyncable": false,
                "timeout": 10000,
                "retries": 0,
                "address": "IP:PORT",
                "lang": "Python3",
                "attrs": {}
            },
            {
                "urn": "/a/b/c",
                "kind": "Restful",
                "namespace": "",
                "name": "",
                "version": "1.0.0",
                "proto": "grpc",
                "codec": "json",
                "asyncable": false,
                "timeout": 10000,
                "retries": 0,
                "address": "IP:PORT",
                "lang": "Python3",
                "attrs": {}
            }
        ]
    }
}'
```

### <a id="5">5 附录</a>

更多高级能力详见： [Mesh高级能力](https://github.com/be-io/mesh/tree/master/client/golang/prsim)

更多使用案例详见：[Mesh使用案例](https://github.com/be-io/mesh/tree/master/examples)