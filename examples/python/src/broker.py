#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
import os
import pickle
import unittest

import grpc

import mesh.log as log
from patch_pb2 import Packet, Metadata
from patch_pb2_grpc import DataTransferServiceStub


class TestBroker(unittest.TestCase):

    @staticmethod
    def generate_packets():
        for _ in range(0, 10):
            packet = Packet(header=Metadata(taskId="1", fromPartyId="1", toPartyId="2", serviceType=5),
                            body=os.urandom(1024 * 1024 * 200))
            yield packet

    def test_invoke(self):
        with grpc.insecure_channel("127.0.0.1:17304") as channel:
            stub = DataTransferServiceStub(channel)
            packet = Packet()
            packet.header.operator = "POST"
            packet.body = b'{"access_token":"JG0100011500000000"}'
            x = stub.unaryCall(packet, 5, [
                ("mesh-urn", "grpc.route.tensor.0001000000000000000000000000000000000.JG0100001100000000.ducesoft.com"),
                ("object-key", "/v1/party/2023030313103945464992/host/JG0100011600000000/update")],
                               None, True, False)
            log.info(f"{x}")

    def test_push(self):
        with grpc.insecure_channel("127.0.0.1:17304") as channel:
            stub = DataTransferServiceStub(channel)
            packet = Packet()
            packet.header.operator = "POST"
            packet.body = b'{"access_token":"JG0100011500000000"}'

            packet_iterator = self.generate_packets()

            y = stub.push(packet_iterator, 300, [
                ("mesh-urn", "grpc.route.tensor.0001000000000000000000000000000000000.jg0100001100000000.ducesoft.com"),
                ("object-key", "1transfer_value#2#3")],
                          None, True, False)
            log.info(f"{y}")

    def test_pickle(self):
        v = pickle.dumps({"key": "1", "split_count": 1})
        log.info(v)
        x = pickle.loads(v)
        log.info(x)


if __name__ == '__main__':
    unittest.main()
