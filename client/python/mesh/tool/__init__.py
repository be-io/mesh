#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

import inspect
import os
import random
import sys
from pathlib import Path
from typing import List, Any, Type, Dict

from mesh.environ import URI, Addrs, System
from mesh.macro import T, Compatible
from mesh.tool.snowflake import UUID


def init():
    """ init function """
    pass


__uuid = UUID(random.randint(0, 31), random.randint(0, 31))


def anyone(*args: str) -> str:
    if not args:
        return ""

    for arg in args:
        if arg and "" != arg:
            return arg

    return ""


def ternary(expr: bool, left: T, right: T) -> T:
    if expr:
        return left
    return right


def get_property(dft: str, name: List[str]) -> str:
    return System.environ().get_property(dft, name)


def get_ip() -> str:
    return System.environ().get_ip()


def get_hostname() -> str:
    return System.environ().get_hostname()


def get_mesh_address() -> Addrs:
    return System.environ().get_mesh_address()


def get_mesh_runtime() -> URI:
    return System.environ().get_mesh_runtime()


def get_mesh_name() -> str:
    return System.environ().get_mesh_name()


def get_mesh_mode() -> int:
    return System.environ().get_mesh_mode()


def get_mesh_direct() -> str:
    return System.environ().get_mesh_direct()


def get_max_channels() -> int:
    return System.environ().get_max_channels()


def get_min_channels() -> int:
    return System.environ().get_min_channels()


def get_packet_size() -> int:
    return System.environ().get_packet_size()


def required(*args: Any) -> bool:
    if not args:
        return False

    for arg in args:
        if arg is not None and "" != arg:
            continue
        return False

    return True


def optional(v: Any) -> bool:
    return not required(v)


def new_trace_id() -> str:
    return f"{System.environ().get_ip_hex()}{next_id()}"


def new_span_id(span_id: str, index: int) -> str:
    if optional(span_id):
        return "0"
    if span_id.__len__() > 255:
        return "0"
    return f"{span_id}.{index}"


def split(v: str, sep: str) -> [str]:
    if not v or "" == v:
        return []
    return v.split(sep)


def get_declared_methods(reference: Type[T]) -> Dict[Type[T], List[Any]]:
    kinds: List[Type[T]] = []
    for k, v in vars(reference).items():
        if inspect.isclass(v) and issubclass(reference, v):
            kinds.append(v)
    if inspect.isabstract(reference):
        kinds.append(reference)
    methods: Dict[Type[T], List[Any]] = {}
    for kind in kinds:
        methods[kind] = []
        for method in inspect.getmembers(kind, inspect.isfunction):
            methods[kind].append(method[1])
        for method in inspect.getmembers(kind, inspect.ismethod):
            methods[kind].append(method[1])
    return methods


def next_id() -> str:
    return __uuid.new_id()


def get_origin(tp):
    return Compatible.get_origin(tp)


def get_args(tp):
    return Compatible.get_args(tp)


def pwd() -> str:
    """
    Returns the name of the project root directory.
    :return: Project root directory name
    """

    # stack trace history related to the call of this function
    frame_stack: [inspect.FrameInfo] = inspect.stack()

    # get info about the module that has invoked this function
    # (index=0 is always this very module, index=1 is fine as long this function is not called by some other
    # function in this module)
    frame_info: inspect.FrameInfo = frame_stack[1]

    # if there are multiple calls in the stacktrace of this very module, we have to skip those and take the first
    # one which comes from another module
    if frame_info.filename == __file__:
        for frame in frame_stack:
            if frame.filename != __file__:
                frame_info = frame
                break

    # path of the module that has invoked this function
    caller_path: str = frame_info.filename

    # absolute path of the of the module that has invoked this function
    caller_absolute_path: str = os.path.abspath(caller_path)

    # get the top most directory path which contains the invoker module
    paths: [str] = [p for p in sys.path if p in caller_absolute_path]
    paths.sort(key=lambda p: len(p))
    caller_root_path: str = paths[0]

    if not os.path.isabs(caller_path):
        # file name of the invoker module (eg: "mymodule.py")
        caller_module_name: str = Path(caller_path).name

        # this piece represents a subpath in the project directory
        # (eg. if the root folder is "myproject" and this function has ben called from myproject/foo/bar/mymodule.py
        # this will be "foo/bar")
        project_related_folders: str = caller_path.replace(os.sep + caller_module_name, '')

        # fix root path by removing the undesired subpath
        caller_root_path = caller_root_path.replace(project_related_folders, '')

    dir_name: str = Path(caller_root_path).name

    return dir_name
