#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from mesh.macro.ark import ark, T, A, B, C
from mesh.macro.binding import binding, Binding
from mesh.macro.cause import Cause
from mesh.macro.codec import serializable
from mesh.macro.compatible import Compatible
from mesh.macro.index import index, Index
from mesh.macro.inspect import Returns, Parameters, Inspector
from mesh.macro.loader import ServiceLoader
from mesh.macro.mpi import mpi, MPI, MethodProxy
from mesh.macro.mps import mps, MPS
from mesh.macro.proxy import Proxy, InvocationHandler
from mesh.macro.spi import spi, SPI
from mesh.macro.types import Types

__all__ = ("mpi",
           "mps",
           "binding",
           "index",
           "spi",
           "serializable",
           "ark",
           "T",
           "A",
           "B",
           "C",
           "Index",
           "SPI",
           "Proxy",
           "InvocationHandler",
           "MPI",
           "MPS",
           "Binding",
           "ServiceLoader",
           "Types",
           "Cause",
           "Returns",
           "Parameters",
           "Inspector",
           "MethodProxy",
           "Compatible",
           )


def init():
    """ init function """
    pass
