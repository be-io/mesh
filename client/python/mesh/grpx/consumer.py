#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

from concurrent.futures import Future
from typing import Iterator, Union, Any

import grpc

import mesh.tool as tool
from mesh.cause import TimeoutException, MeshException, MeshCode
from mesh.context import Mesh
from mesh.environ import System
from mesh.grpx.bfia import BfiaBinding, bfia_consume
from mesh.grpx.channels import GrpcChannels
from mesh.kinds import Reference
from mesh.macro import spi
from mesh.mpc.consumer import Consumer
from mesh.mpc.invoker import Execution


@spi("grpc")
class GrpcConsumer(Consumer):

    def __init__(self):
        self.channel = GrpcChannels()

    def __enter__(self):
        self.start()

    def __exit__(self, exc_type, exc_val, exc_tb):
        """ Close """
        self.channel.__exit__(exc_type, exc_val, exc_tb)

    def start(self):
        """ Do not need to implement """
        pass

    def consume(self, address: str, urn: str, execution: Execution[Reference], inbound: bytes, args: Any) -> Future:
        try:
            binding = BfiaBinding(urn, args)
            if "" != binding.path:
                return bfia_consume(address, binding, execution, inbound, self.channel)

            ct = Mesh.context().get_attribute(Mesh.TIMEOUT)
            timeout = (ct if ct else execution.schema().timeout) / 1000
            outbound = self.channel.unary(address, inbound, timeout, [])
            # GrpcFuture(future, self.reset_if_unconnected(channel))
            future = Future()
            future.set_result(outbound)
            return future
        except grpc.FutureTimeoutError:
            raise TimeoutException(f'Invoke {execution.schema().urn} timeout with {execution.schema().timeout}ms')
        except grpc.FutureCancelledError:
            raise MeshException(MeshCode.SYSTEM_ERROR.get_code(), f'Invoke {execution.schema().urn} canceled')
        except grpc.RpcError as e:
            if isinstance(e, grpc.Call):
                if e.code() == grpc.StatusCode.DEADLINE_EXCEEDED:
                    raise TimeoutException(
                        f'Invoke {execution.schema().urn} timeout with {execution.schema().timeout}ms')
                if e.code() == grpc.StatusCode.UNIMPLEMENTED:
                    raise MeshException(MeshCode.NOT_FOUND.get_code(), MeshCode.NOT_FOUND.get_message())
                raise MeshException(MeshCode.SYSTEM_ERROR.get_code(),
                                    f'Invoke {execution.schema().urn} {e.code()} {e.details()}')
            raise e

    @staticmethod
    def resolve(address: str):
        hosts = tool.split(address, ":")
        if hosts.__len__() < 2:
            return f'{hosts[0]}:{System.environ().get_default_mesh_port()}'
        return address

    def reset_if_unconnected(self, channel: GrpcChannels):
        def x(e: BaseException):
            pass

        return x


class AsyncInvokeObserver(Future, grpc.Call):

    def __init__(self, rendezvous: Union[grpc.Future, grpc.Call]):
        super().__init__()
        self.rendezvous = rendezvous

    def initial_metadata(self):
        self.rendezvous.initial_metadata()

    def trailing_metadata(self):
        return self.rendezvous.trailing_metadata()

    def code(self):
        return self.rendezvous.code()

    def details(self):
        return self.rendezvous.details()

    def is_active(self):
        return self.rendezvous.is_active()

    def time_remaining(self):
        return self.rendezvous.time_remaining()

    def add_callback(self, callback):
        return self.rendezvous.add_callback(callback)


class IterableBuffer(Iterator):

    def __init__(self, buffer: bytes):
        self.buffer = list()
        self.buffer.append(buffer)

    def __next__(self) -> bytes:
        if self.buffer:
            return self.buffer.pop()
        raise StopIteration()


class GrpcFuture(grpc.Future):

    def __init__(self, future: grpc.Future, hooks: Any) -> None:
        self.future = future
        self.hooks = hooks

    def cancel(self, msg: Any = ...) -> bool:
        return self.future.cancel()

    def cancelled(self):
        return self.future.cancelled()

    def running(self):
        return self.future.running()

    def done(self):
        return self.future.done()

    def result(self, timeout=None):
        return self.future.result(timeout)

    def exception(self, timeout=None):
        return self.future.exception(timeout)

    def traceback(self, timeout=None):
        return self.future.traceback(timeout)

    def add_done_callback(self, fn):
        return self.future.add_done_callback(fn)
