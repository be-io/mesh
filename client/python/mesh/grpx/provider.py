#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
import multiprocessing
from concurrent import futures
from typing import Any

import grpc

import mesh.log as log
import mesh.tool as tool
from mesh.cause import ValidationException
from mesh.grpx import GrpcBindableService, GrpcInterceptor
from mesh.macro import spi
from mesh.mpc import Provider


@spi("grpc")
class GrpcProvider(Provider):
    """
    Grpc provider wrap the grpc protocol implementation.
    """

    def __init__(self):
        self.interceptor = GrpcInterceptor()
        self.service = GrpcBindableService()
        self.server = grpc.server(
            futures.ThreadPoolExecutor(max_workers=multiprocessing.cpu_count() * 10),
            [self.service],
            [self.interceptor],
            [('grpc.max_message_length', 1 << 30),
             ('grpc.max_send_message_length', 1 << 30),
             ('grpc.max_receive_message_length', 1 << 30)
             ],
            None,
            grpc.Compression.Gzip,
            False
        )

    def start(self, address: str, tc: Any):
        if tool.optional(address):
            raise ValidationException("GRPC address cant be empty.")
        self.server.add_insecure_port(address)
        self.server.add_generic_rpc_handlers([self.service])
        self.server.start()
        log.info(f"Listening and serving HTTP 2.0 on {address}")

    def close(self):
        self.server.stop(grace=12)
        log.info(f"Graceful stop HTTP 2.0 serving")
