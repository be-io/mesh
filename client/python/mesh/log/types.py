#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

import io
from abc import ABC, abstractmethod
from enum import Enum


def atoi(level: "Level") -> int:
    if isinstance(level, int):
        return level
    return int(level.code)


class Level(Enum):
    FATAL = 1
    ERROR = 2 | atoi(FATAL)
    WARN = 4 | atoi(ERROR)
    INFO = 8 | atoi(WARN)
    DEBUG = 16 | atoi(INFO)
    STACK = 32 | atoi(DEBUG)
    ALL = atoi(FATAL) | atoi(ERROR) | atoi(WARN) | atoi(INFO) | atoi(DEBUG) | atoi(STACK)

    def __init__(self, code: int) -> None:
        self.code = code

    def __str__(self) -> str:
        if Level.FATAL.match(self):
            return 'fatal'
        if Level.ERROR.match(self):
            return 'error'
        if Level.WARN.match(self):
            return 'warn'
        if Level.INFO.match(self):
            return 'info'
        if Level.DEBUG.match(self):
            return 'debug'
        if Level.STACK.match(self):
            return 'stack'
        return 'all'

    def match(self, level: "Level") -> bool:
        return self.code == (self.code & level.code)

    def match_name(self, level: str) -> bool:
        return self.name == level


class Logger(ABC):

    @abstractmethod
    def name(self) -> str:
        pass

    @abstractmethod
    def info(self, fmt: str, *args: object):
        pass

    @abstractmethod
    def warn(self, fmt: str, *args: object):
        pass

    @abstractmethod
    def error(self, fmt: str, *args: object):
        pass

    @abstractmethod
    def debug(self, fmt: str, *args: object):
        pass

    @abstractmethod
    def fatal(self, fmt: str, *args: object):
        pass

    @abstractmethod
    def stack(self, fmt: str, *args: object):
        pass

    @abstractmethod
    def writer(self) -> io.BytesIO:
        pass

    @abstractmethod
    def level(self, level: Level):
        pass
