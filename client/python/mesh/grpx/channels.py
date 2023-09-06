#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

import random
from typing import Any, List, Dict, Callable

import grpc

import mesh.tool as tool
from mesh.grpx.interceptor import GrpcInterceptor
from mesh.grpx.marshaller import GrpcMarshaller


class GrpcChannels:

    def __init__(self):
        self.channels: Dict[str, GrpcMultiplexChannel] = {}
        self.marshaller = GrpcMarshaller()
        self.serializer = self.marshaller.serialize
        self.deserializer = self.marshaller.deserialize
        self.interceptor = GrpcInterceptor()
        self.path = "/mesh-rpc/v1"
        self.streams: Dict[str, GrpcMultiplexStream] = dict()

    def __exit__(self, exc_type, exc_val, exc_tb):
        for _, channel in self.channels.items():
            channel.close()

    def create_if_absent(self, address: str) -> "GrpcMultiplexChannel":
        channel = self.make_channel(address)
        if channel.counter > 100:
            self.channels.__delitem__(address)
            channel.close()
            channel = self.make_channel(address)
        return channel

    def make_channel(self, address: str) -> "GrpcMultiplexChannel":
        channel = self.channels.get(address, None)
        if not channel:
            self.channels[address] = channel = GrpcMultiplexChannel(address, self.interceptor)
        return channel

    def create_stream(self, address: str) -> grpc.StreamStreamMultiCallable:
        channel = self.create_if_absent(address)
        return channel.stream_stream(self.path, self.serializer, self.deserializer)

    def unary(self, address: str, inbound: bytes, timeout: Any, metadata: Any) -> bytes:
        channel = self.create_if_absent(address)
        conn: grpc.UnaryUnaryMultiCallable = channel.unary_unary(self.path, self.serializer, self.deserializer)
        return conn(inbound, timeout, metadata, None, True, False)

    def stream(self, address: str, inbound: bytes, timeout: Any, metadata: Any) -> bytes:
        stream = self.streams.get(address, None)
        if stream is None:
            stream = self.streams[address] = GrpcMultiplexStream(address, lambda addr: self.create_stream(addr))
        return stream.next(inbound, timeout, metadata)


class GrpcMultiplexStream:

    def __init__(self, address: str, streamer: Callable[[str], grpc.StreamStreamMultiCallable]):
        self.address = address
        self.stream = streamer(address)
        self.streamer = streamer

    def next(self, inbound: bytes, timeout: Any, metadata: Any) -> bytes:
        return self.stream(inbound, timeout, metadata, None, True, False)


class GrpcMultiplexChannel(grpc.Channel):
    """
    https://github.com/grpc/grpc/blob/598e171746576c5398388a4c170ddf3c8d72b80a/include/grpc/impl/codegen/grpc_types.h#L170
    """

    def close(self):
        for _, channel in self.channels:
            channel.close()

    def subscribe(self, callback, try_to_connect=False):
        return self.select().subscribe(callback, try_to_connect)

    def unsubscribe(self, callback):
        return self.select().unsubscribe(callback)

    def unary_unary(self, method, request_serializer=None, response_deserializer=None):
        return self.select().unary_unary(method, request_serializer, response_deserializer)

    def unary_stream(self, method, request_serializer=None, response_deserializer=None):
        return self.select().unary_stream(method, request_serializer, response_deserializer)

    def stream_unary(self, method, request_serializer=None, response_deserializer=None):
        return self.select().stream_unary(method, request_serializer, response_deserializer)

    def stream_stream(self, method, request_serializer=None, response_deserializer=None):
        return self.select().stream_stream(method, request_serializer, response_deserializer)

    def __enter__(self):
        for _, channel in self.channels:
            channel.__enter__()

    def __exit__(self, exc_type, exc_val, exc_tb):
        for _, channel in self.channels:
            channel.__exit__(exc_type, exc_val, exc_tb)

    def __init__(self, address: str, *interceptors):
        self.counter = 0
        self.address = address
        self.interceptors = interceptors
        self.channels: List[grpc.Channel] = []
        self.compression_options = {
            "none": grpc.Compression.NoCompression,
            "gzip": grpc.Compression.Gzip,
        }

    def default_compression(self) -> Any:
        return self.compression_options["gzip"]

    def reset(self, channel: grpc.Channel):
        self.channels.append(channel)

    def select(self) -> grpc.Channel:
        """
        grpc.ssl_target_name_override
        grpc.default_authority
        https://github.com/grpc/grpc/blob/master/include/grpc/impl/codegen/compression_types.h
        https://github.com/grpc/grpc/blob/master/include/grpc/impl/codegen/grpc_types.h
        """
        if self.channels.__len__() > 0:
            return self.channels[random.randint(0, self.channels.__len__() - 1)]

        for _ in range(tool.get_max_channels()):
            options = self.default_options()
            channel = grpc.insecure_channel(self.address, options, compression=self.default_compression())
            self.channels.append(grpc.intercept_channel(channel, *self.interceptors))

        return self.channels[random.randint(0, self.channels.__len__() - 1)]

    def default_options(self):
        return [
            ("grpc.enable_retries", True),
            ("grpc.default_compression_algorithm", self.default_compression()),
            ("grpc.client_idle_timeout_ms", 1000 * 32),
            ("grpc.max_send_message_length", 1 << 30),
            ("grpc.max_receive_message_length", 1 << 30),
            ("grpc.max_connection_age_ms", 1000 * 60),
            ("grpc.max_connection_age_grace_ms", 1000 * 12),
            ("grpc.per_message_compression", True),
            ("grpc.per_message_decompression", True),
            ("grpc.enable_deadline_checking", True),
            ("grpc.keepalive_time_ms", 1000 * 12),
            ("grpc.keepalive_timeout_ms", 1000 * 12),
        ]

    def incr(self):
        self.counter += 1
