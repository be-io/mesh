#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from asyncio import Future
from typing import Any

import mesh.log as log
from mesh.context import Mesh
from mesh.macro import T, ServiceLoader
from mesh.mpc.eden import Eden
from mesh.mpc.filter import Filter, PROVIDER
from mesh.mpc.invoker import Invoker, Invocation


class ServiceInvoker(Invoker[T]):

    def __init__(self, service: Any):
        self.service = service

    def run(self, invocation: Invocation) -> Any:
        result = invocation.get_inspector().invoke(self.service, invocation.get_arguments())
        if not isinstance(result, Future):
            return result
        try:
            eden = ServiceLoader.load(Eden).get_default()
            execution = eden.infer(Mesh.context().get_urn())
            timeout = execution.schema().timeout
            return result.result(timeout=max(timeout, 3000))
        except BaseException as e:
            log.error(f'Invoke {Mesh.context().get_urn()} fault because of {e}')
            raise e


class ServiceInvokeHandler(Invoker):
    """

    """

    def __init__(self, service: Any):
        self.invoker = Filter.composite(ServiceInvoker(service), PROVIDER)

    def run(self, invocation: Invocation) -> Any:
        return self.invoker.run(invocation)
