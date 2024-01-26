#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

import argparse
import os
import random
import socket
from contextlib import closing
from enum import Enum
from typing import List

from mesh.environ.cache import Cache

DEFAULT_MESH_PORT = 570


class StatefulServer:
    def __init__(self, available: bool, address: str):
        self.available = available
        self.address = address


class Addrs:

    def __init__(self, servers: str):
        self.servers: List[StatefulServer] = []
        self.available_addrs = []
        if not servers:
            return
        for server in servers.split(","):
            if server and server.__len__() > 0 and not self.available_addrs.__contains__(server):
                available_addr = server if server.__contains__(":") else f"{server}:{DEFAULT_MESH_PORT}"
                self.available_addrs.append(available_addr)
                self.servers.append(StatefulServer(True, available_addr))

    def any(self) -> str:
        addrs = self.available_addrs
        if not addrs or addrs.__len__() < 1:
            return ''
        return addrs[random.randint(0, addrs.__len__() - 1)]

    def many(self) -> List[str]:
        addrs = self.available_addrs
        if not addrs or addrs.__len__() < 1:
            return []
        return addrs.copy()

    def all_servers(self) -> List[StatefulServer]:
        return self.servers

    def available(self, addr: str, available: bool):
        for server in self.servers:
            if addr == server.address and server.available is not available:
                server.available = available
                available_addrs = []
                for serv in self.servers:
                    if serv.available:
                        available_addrs.append(serv.address)
                if available_addrs.__len__() > 0:
                    self.available_addrs = available_addrs
                else:
                    available_addrs.append(self.servers[random.randint(0, self.servers.__len__() - 1)].address)
                    self.available_addrs = available_addrs
                return


class URI:
    def __init__(self, hostname: str, port: int):
        self.hostname = hostname
        self.port = port

    def __str__(self) -> str:
        return f"{self.hostname}:{self.port}"

    @staticmethod
    def parse(addr: str) -> "URI":
        if not addr.__contains__(":"):
            return URI(addr, 0)
        pair = addr.split(":")
        return URI(pair[0], int(pair[1]) if pair[1].isdigit() else 0)


class Mode(Enum):
    Disable = 1
    Failfast = 2
    Nolog = 4
    JsonLogFormat = 8
    RCache = 16
    PHeader = 32
    Metrics = 64
    RLog = 128
    MGrpc = 256
    PermitCirculate = 512
    NoStdColor = 1024

    def __init__(self, v: int) -> None:
        self.v = v

    def enable(self) -> bool:
        return self.match(System.mode())

    def match(self, mode: int) -> bool:
        return (self.v & mode) == self.v


class Environ:
    __mesh_address_keys = ['mesh.address', 'mesh_address', 'mesh-address', 'MESH_ADDRESS']
    __mesh_runtime_keys = ['mesh.runtime', 'mesh_runtime', 'mesh-runtime', 'MESH_RUNTIME']
    __mesh_mode_keys = ['mesh.mode', 'mesh_mode', 'mesh-mode', 'MESH_MODE']
    __mesh_name_keys = ['mesh.name', 'mesh_name', 'mesh-name', 'MESH_NAME']
    __mesh_direct_keys = ['mesh.direct', 'mesh_direct', 'mesh-direct', 'MESH_DIRECT']
    __mesh_subset_keys = ['mesh.subset', 'mesh_subset', 'mesh-subset', 'MESH_SUBSET']
    __mesh_min_channels = ['mesh.grpc.channel.min', 'mesh_grpc_channel_min', 'MESH_GRPC_CHANNEL_MIN']
    __mesh_max_channels = ['mesh.grpc.channel.max', 'mesh_grpc_channel_max', 'MESH_GRPC_CHANNEL_MAX']
    __mesh_packet_size = ['mesh.packet.size', 'mesh_packet_size', 'MESH_PACKET_SIZE']

    @staticmethod
    def get_default_mesh_port() -> int:
        return DEFAULT_MESH_PORT

    def get_property(self, dft: str, keys: List[str]) -> str:
        if not keys:
            return dft
        for key in keys:
            if '' != self.get_arguments.get(key, ''):
                return self.get_arguments[key]
            if os.getenv(key):
                return os.getenv(key)
        return dft

    @Cache
    def get_arguments(self):
        parser = argparse.ArgumentParser()
        parser.add_argument("--mesh.address", type=str, default='', required=False)
        parser.add_argument("--mesh.runtime", type=str, default='', required=False)
        parser.add_argument("--mesh.enable", type=str, default='', required=False)
        parser.add_argument("--mesh.name", type=str, default='', required=False)
        parser.add_argument("--mesh.direct", type=str, default='', required=False)
        parser.add_argument("--mesh.grpc.channel.min", type=int, default=20, required=False)
        parser.add_argument("--mesh.grpc.channel.max", type=int, default=200, required=False)
        parser.add_argument("--mesh.packet.size", type=int, default=1 << 26, required=False)
        args, _ = parser.parse_known_args()
        return vars(args)

    @Cache
    def get_mesh_address(self) -> Addrs:
        return Addrs(self.get_property('127.0.0.1', self.__mesh_address_keys))

    @Cache
    def get_mesh_runtime(self) -> URI:
        mr = self.get_property('', self.__mesh_runtime_keys)
        if mr is not None and '' != mr:
            return URI.parse(mr)
        return URI.parse(f'{self.get_ip()}:{self.get_available_port()}')

    @Cache
    def get_mesh_name(self) -> str:
        return self.get_property('unknown', self.__mesh_name_keys)

    @Cache
    def get_mesh_mode(self) -> str:
        return self.get_property('0', self.__mesh_mode_keys)

    @Cache
    def get_mesh_direct(self) -> str:
        return self.get_property('', self.__mesh_direct_keys)

    @Cache
    def get_mesh_subset(self) -> str:
        return self.get_property('', self.__mesh_subset_keys)

    @Cache
    def get_max_channels(self) -> int:
        return int(self.get_property('3', self.__mesh_max_channels))

    @Cache
    def get_min_channels(self) -> int:
        return int(self.get_property(os.cpu_count() * 2, self.__mesh_min_channels))

    @Cache
    def get_packet_size(self) -> int:
        return int(self.get_property(f"{1 << 26}", self.__mesh_packet_size))

    @Cache
    def get_hostname(self) -> str:
        return socket.gethostname()

    @Cache
    def get_ip(self) -> str:
        s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        s.settimeout(0)
        try:
            # doesn't even have to be reachable
            s.connect(('10.255.255.255', 1))
            ip = s.getsockname()[0]
        except Exception as e:
            print(e)
            ip = '127.0.0.1'
        finally:
            s.close()
        return ip

    @Cache
    def get_ip_hex(self) -> str:
        return '{:02X}{:02X}{:02X}{:02X}'.format(*map(int, self.get_ip().split('.')))

    @Cache
    def get_available_port(self) -> str:
        with closing(socket.socket()) as s:
            s.bind(('', 0))
            return s.getsockname()[1]


__environ__ = Environ()
__mode__ = int(__environ__.get_mesh_mode())


class System:

    @staticmethod
    def mode() -> int:
        return __mode__

    @staticmethod
    def environ() -> Environ:
        return __environ__
