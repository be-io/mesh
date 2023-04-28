#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

import inspect
from abc import ABCMeta, abstractmethod


class Interface(metaclass=ABCMeta):
    pass


def interface(cls):
    attrs = {n: abstractmethod(f)
             for n, f in inspect.getmembers(cls, predicate=inspect.isfunction)}

    return type(cls.__name__, (Interface, cls), attrs)
