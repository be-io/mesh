# 节点1设置节点信息
curl --location 'http://127.0.0.1:7304/v1/interconn/node/refresh' --header 'Content-Type: application/json' --data '{
    "version": "1.0.0",
    "node_id": "321",
    "inst_id": "321",
    "inst_name": "321"
}'

# 节点2设置节点信息
curl --location 'http://127.0.0.1:17304/v1/interconn/node/refresh' --header 'Content-Type: application/json' --data '{
    "version": "1.0.0",
    "node_id": "123",
    "inst_id": "123",
    "inst_name": "123"
}'

# 节点1组网
curl --location 'http://127.0.0.1:7304/v1/interconn/net/refresh' --header 'Content-Type: application/json' --data '{
    "node_id": "123",
    "inst_id": "123",
    "address": "192.168.100.63:17304"
}'

# 节点2组网
curl --location 'http://127.0.0.1:17304/v1/interconn/net/refresh' --header 'Content-Type: application/json' --data '{
    "node_id": "321",
    "inst_id": "321",
    "address": "192.168.100.63:7304"
}'

# 节点1发送消息
curl --location 'http://127.0.0.1:7304/v1/interconn/chan/push' --header 'Content-Type: application/json' --header 'x-ptp-target-node-id: 123' --header 'x-ptp-session-id: 1' --data '{
    "topic": "0",
    "metadata": {
        "a": "b"
    },
    "payload": "MTkyLjE2OC4xMDAuNjM6MTczMDQ="
}'

# 节点2接收消息
curl --location 'http://127.0.0.1:17304/v1/interconn/chan/peek' --header 'Content-Type: application/json' --header 'x-ptp-target-node-id: 123' --header 'x-ptp-session-id: 1' --data '{
    "topic": "0",
    "timeout": 10000
}'