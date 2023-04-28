#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
import inspect
import io
import traceback
from datetime import datetime
from io import StringIO

from mesh.environ import Mode
from mesh.log.types import Logger, Level


def stack_tuple() -> (str, str):
    if Mode.Nolog:
        return "", ""
    frame = inspect.stack()[4]
    index = frame.filename.rindex("/")
    name = frame.filename[0 if index < 0 else index + 1:]
    return name, frame.lineno


class NopLogger(Logger):

    def __init__(self):
        self.__level__ = Level.INFO

    def name(self) -> str:
        return 'mesh'

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

    def write(self, level: Level, fmt: str, args):
        if not level.match(self.__level__):
            return

        # discard log
        if Mode.Nolog.enable():
            return

        if not args or len(args) < 1:
            self.do_write(level, fmt)
            return

        if len(args) == 1 and isinstance(args[0], BaseException):
            self.do_write(level, fmt, args[0])
            return

        frags = (fmt if fmt else "").split("{}")
        placeholders = []
        throwable = False
        for arg in args:
            if isinstance(arg, BaseException):
                throwable = True
                continue
            if isinstance(arg, str):
                placeholders.append(arg)
                continue
            placeholders.append(f"{arg}")

        msg = StringIO()
        for idx in range(len(frags)):
            msg.write(frags[idx])
            if len(placeholders) > idx:
                msg.write(placeholders[idx])
                continue
            if idx < len(frags) - 1:
                msg.write('{}')

        self.do_write(level, msg.getvalue(), throwable)

    def do_write(self, level: Level, msg: str, e: bool = False):
        print(f"{datetime.now().strftime('%Y-%m-%d %H:%M:%S')} {level.name} {msg}{self.stack_tuple()}")
        if e:
            print(traceback.format_exc())

    @staticmethod
    def stack_tuple() -> str:
        if not Mode.Metrics.enable():
            return ""
        frame = inspect.stack()[4]
        index = frame.filename.rindex("/")
        name = frame.filename[0 if index < 0 else index + 1:]
        return f" {name}:{frame.lineno}"
