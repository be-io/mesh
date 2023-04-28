#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
import io
import threading
import time
from typing import List

from mesh import tool, log
from mesh.boost.disruptor import Disruptor, Consumer
from mesh.kinds import Document
from mesh.log.types import Logger, Level
from mesh.macro import mpi, spi
from mesh.mpc import Mesh
from mesh.prsim import DataHouse, Key

DISCARD = Key("mesh.syslog.discard")
BUFFER_SIZE = 10
LIMITER = threading.Semaphore(3)


@spi("remote")
class RemoteLogger(Logger):

    def __init__(self, level: Level):
        self.__level__ = level
        self.disruptor = Disruptor(name='mesh', size=64)
        self.disruptor.register_consumer(DocumentConsumer())

    def name(self) -> str:
        return 'remote'

    def info(self, fmt: str, *args: object):
        self.write(Level.INFO, fmt, *args)

    def warn(self, fmt: str, *args: object):
        self.write(Level.WARN, fmt, *args)

    def error(self, fmt: str, *args: object):
        self.write(Level.ERROR, fmt, *args)

    def debug(self, fmt: str, *args: object):
        self.write(Level.DEBUG, fmt, *args)

    def fatal(self, fmt: str, *args: object):
        self.write(Level.FATAL, fmt, *args)

    def stack(self, fmt: str, *args: object):
        self.write(Level.STACK, fmt, *args)

    def writer(self) -> io.BytesIO:
        return io.BytesIO()

    def level(self, level: Level):
        self.__level__ = level

    def write(self, level: Level, fmt: str, *args):
        if not level.match(self.__level__):
            return
        is_discard = Mesh.context().get_attribute(DISCARD)
        if is_discard is not None and is_discard:
            return
        metadata = {
            "name": "",
            "level": level.name,
            "host": tool.get_hostname(),
            "ip": tool.get_ip(),
            "mesh_runtime": str(tool.get_mesh_runtime()),
            "mesh_trace_id": Mesh.context().get_trace_id(),
            "mesh_span_id": Mesh.context().get_span_id(),
            "mesh_name": tool.get_mesh_name(),
            "mesh_run_mode": str(Mesh.context().get_run_mode())
        }
        document = Document(metadata, fmt, time.time_ns())
        self.disruptor.produce([document])


class DocumentConsumer(Consumer):

    @mpi
    def data_house(self) -> DataHouse:
        """"""
        pass

    def consume(self, elements):
        return Mesh.context_safe(lambda: self.do_consume(elements))

    def do_consume(self, elements):
        Mesh.context().set_attribute(DISCARD, True)
        if len(elements) > BUFFER_SIZE - 1:
            try:
                self.data_house.writes(elements)
            except BaseException as e:
                self.exception(e, elements)
            finally:
                elements.clear()
            return
        if not LIMITER.acquire(True):
            return
        try:
            self.data_house.write(elements)
        except BaseException as e:
            self.exception(e, elements)
        finally:
            LIMITER.release(1)
            elements.clear()

    @staticmethod
    def exception(e: BaseException, event: List[Document]):
        try:
            for document in event:
                ln = document.metadata.get("level", '')
                if Level.ERROR.match_name(ln):
                    log.error(document.content)
                elif Level.WARN.match_name(ln):
                    log.warn(document.content)
                else:
                    log.info(document.content)
            log.error("", e)
        except BaseException as c:
            log.error("", c)
        finally:
            event.clear()
