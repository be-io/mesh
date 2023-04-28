#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

from mesh.log.nop import NopLogger
from mesh.log.types import Level, Logger

__all__ = ("Level", "Logger")

mog = NopLogger()


def init():
    """ init function """
    pass


def info(fmt: str, *args: object):
    mog.info(fmt, args)


def warn(fmt: str, *args: object):
    mog.warn(fmt, args)


def error(fmt: str, *args: object):
    mog.error(fmt, args)


def debug(fmt: str, *args: object):
    mog.debug(fmt, args)


def fatal(fmt: str, *args: object):
    mog.fatal(fmt, args)


def stack(fmt: str, *args: object):
    mog.debug(fmt, args)


def level(lev: Level):
    mog.level(lev)
