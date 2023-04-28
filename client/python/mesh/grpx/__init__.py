#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

from mesh.grpx.bindservice import GrpcBindableService
from mesh.grpx.channels import GrpcChannels
from mesh.grpx.consumer import GrpcConsumer
from mesh.grpx.interceptor import GrpcInterceptor
from mesh.grpx.marshaller import GrpcMarshaller
from mesh.grpx.provider import GrpcProvider

MPC = "/mesh-rpc/v1"

__all__ = (
    "GrpcBindableService",
    "GrpcChannels",
    "GrpcConsumer",
    "GrpcInterceptor",
    "GrpcMarshaller",
    "GrpcProvider"
)


def init():
    """ init function """
    pass
