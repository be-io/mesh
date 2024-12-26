## 自签证书安装说明

### 双向认证证书签发说明

* 每个节点需要有一对根证书，该根证书在节点初始化时写入到节点内
* 每对节点间通信需一对通信证书（可复用一对证书或各自独立使用一对），该证书由本节点根证书签发
* 每个节点把自身的通信证书和对方的根证书写入组网信息内

所以两种模式：

* 每个节点至少2对证书，1对根证书，1对通信证书，通信证书由根证书派生。
* 每个节点至少1+N（N为组网节点数）对证书，1对根证书，N对通信证书，通信证书由根证书派生。

### 参考流程

#### 1 初始化节点写入根证书

Mesh默认会内置节点根证书，如下接口可进行初始化重置。

```shell
curl --location 'https://<address>/v1/interconn/node/refresh' --header 'Content-Type: application/json' --data '{
    "version": "1.0.0",
    "node_id": "<节点编号>",
    "inst_id": "<机构编号>",
    "inst_name": "<节点名称>",
    "root_crt": "<节点根证书>", 
    "root_key": "<节点根证书私钥>"
}'
```

#### 2 覆盖默认节点间通信证书

Mesh默认会内置通信证书，如下接口可进行初始化重置。

```shell
curl --location 'https://<address>/v1/interconn/net/refresh' --header 'Content-Type: application/json' --data '{
    "node_id": "<对方节点编号>",
    "inst_id": "<对方机构编号>",
    "address": "<对方节点地址>",
    "host_crt": "<本节点服务端通信证书，只有本节点持有>",
    "host_key": "<本节点服务端通信证书私钥，只有本节点持有>",
    "host_root": "<本节点服务端根证书，节点初始化默认加载节点根证书，此参数可覆盖默认行为>",
    "guest_crt": "<本节点客户端通信证书>",
    "guest_key": "<本节点客户端通信证书>",
    "guest_root": "<对方节点根证书，用于双向认证验证>",
}'
```

### 其他注意事项

* 签发证书的时候需要注意域名和IP SANs，需要和实际的保持一致。