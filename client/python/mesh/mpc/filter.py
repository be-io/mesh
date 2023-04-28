#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from abc import abstractmethod
from typing import Any

from mesh.macro import spi, T, ServiceLoader
from mesh.mpc.invoker import Invocation, Invoker

PROVIDER = "PROVIDER"
CONSUMER = "CONSUMER"


@spi(pattern="*")
class Filter:

    @abstractmethod
    def invoke(self, invoker: Invoker, invocation: Invocation) -> Any:
        pass

    @staticmethod
    def composite(invoker: Invoker[T], pattern: str) -> Invoker[T]:
        """
        Composite the filter spi providers as a invoker.
        :param invoker: Delegate invoker.
        :param pattern: Filter pattern.
        :return: Composite invoker.
        """
        filters = ServiceLoader.load(Filter).list(pattern)
        last: Invoker[T] = invoker
        for index, filter_ in enumerate(filters):
            next_: Invoker[T] = last
            last = FilteredInvoker(next_, filter_)
        return last


class FilteredInvoker(Invoker[T]):

    def __init__(self, next_: Invoker[T], filter_: Filter):
        self.next_ = next_
        self.filter_ = filter_

    def run(self, invocation: Invocation) -> Any:
        return self.filter_.invoke(self.next_, invocation)


def invoke(invoker: Invoker[T], filtering: Filter, invocation: Invocation):
    return filtering.invoke(invoker, invocation)
