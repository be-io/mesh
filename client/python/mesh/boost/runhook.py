#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

import atexit
import signal
import sys
import traceback
from typing import Any, Type

import mesh.log as log


class Runtime(object):

    def __init__(self):
        self.sys_exist_hook = sys.exit
        self.sys_except_hook = sys.excepthook
        sys.exit = self.exit_hook
        sys.excepthook = self.except_hook

    def exit_hook(self, status=0, *args):
        if self.sys_exist_hook:
            self.sys_exist_hook(status)

    def except_hook(self, kind: Type[BaseException], e: BaseException, trace_type: Any, *args):
        if self.sys_except_hook:
            self.sys_except_hook(kind, e, trace_type)

    @staticmethod
    def safe_exec_hook(hook: Any):
        try:
            hook()
        except BaseException as e:
            log.error(f"{traceback.format_exception_only(type(e), e)}")

    def add_shutdown_hook(self, hook: Any):
        def safe_hook():
            return self.safe_exec_hook(hook)

        atexit.register(safe_hook)
        signal.signal(signal.SIGTERM, safe_hook)
        signal.signal(signal.SIGINT, safe_hook)

    @staticmethod
    def get_runtime() -> "Runtime":
        return runtime


runtime = Runtime()
