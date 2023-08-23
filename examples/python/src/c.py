#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

import requests

from mesh.ptp.x_pb2 import Inbound
from mesh.ptp.y_pb2 import PushInbound, PopInbound, TransportOutbound


def main():
    push = PushInbound()
    push.topic = "321"
    push.payload = b'1234567'
    inbound = Inbound()
    inbound.payload = push.SerializeToString()
    response = requests.post(url='http://127.0.0.1:17304/v1/interconn/chan/invoke', headers={
        'Content-Type': 'application/x-protobuf',
        'x-ptp-tech-provider-code': 'LX',
        'x-ptp-trace-id': '1',
        'x-ptp-session-id': '1',
        'x-ptp-token': 'x',
        'x-ptp-source-node-id': '2',
        'x-ptp-source-inst-id': '2',
        'x-ptp-target-node-id': '1',
        'x-ptp-target-inst-id': '1',
        'x-ptp-uri': 'https://ptp.cn/v1/interconn/chan/push',
    }, data=inbound.SerializeToString())
    print(response)

    pop = PopInbound()
    pop.topic = "321"
    pop.timeout = 1000
    response = requests.post(url='http://127.0.0.1:17304/v1/interconn/chan/pop', headers={
        'Content-Type': 'application/x-protobuf',
        'x-ptp-tech-provider-code': 'LX',
        'x-ptp-trace-id': '1',
        'x-ptp-session-id': '1',
        'x-ptp-token': 'x',
        'x-ptp-source-node-id': '2',
        'x-ptp-source-inst-id': '2',
        'x-ptp-target-node-id': '1',
        'x-ptp-target-inst-id': '1',
    }, data=pop.SerializeToString())

    t = TransportOutbound()
    t.ParseFromString(response.content)
    print(t.code)
    print(t.message)
    print(str(t.payload))


if __name__ == "__main__":
    main()
