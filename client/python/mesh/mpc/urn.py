#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

from itertools import repeat
from typing import List

import mesh.log as log

MESH_DOMAIN = "trustbe.cn"
LOCAL_NODE_ID = "LX0000000000000"
LOCAL_INST_ID = "JG0000000000000000"


class URN:
    """
    Like:

    create.tenant.omega . 0001 . json . http2 . lx000001 . trustbe.net
    -------------------   ----   ----   -----   --------   -----------
    name            flag   codec  proto     node       domain

    """

    @staticmethod
    def parse(urn: str) -> "URN":
        instance = URN()
        if urn is None or urn.isspace():
            log.warn(f'Unresolved urn {urn}')
            return instance
        names = URN.as_array(urn)
        if len(names) < 5:
            log.warn(f'Unresolved urn {urn}')
            return instance
        names.reverse()
        instance.name = "{}.{}".format(names[1], names[0])
        instance.node_id = names[2]
        instance.flag = URNFlag.parse(names[3])
        instance.name = ".".join(names[4:])
        return instance

    def __init__(self):
        self.domain = MESH_DOMAIN
        self.node_id = ""
        self.flag = URNFlag()
        self.name = ""

    def __str__(self):
        return self.string()

    def match_node(self, node_id: str) -> bool:
        return self.node_id.lower() == node_id.lower()

    def string(self):
        urn = self.as_array(self.name)
        urn.reverse()
        urn.append(self.flag.string())
        urn.append(self.node_id.lower())
        if self.domain.isspace():
            urn.append(MESH_DOMAIN)
        else:
            urn.append(self.domain.lower())
        return ".".join(urn)

    @staticmethod
    def as_array(urn: str) -> List[str]:
        pairs = urn.split(".")
        names = []
        buff = []
        for pair in pairs:
            if pair and pair.startswith("${"):
                buff.append(pair)
                buff.append('.')
                continue
            if pair and pair.endswith("}"):
                buff.append(pair)
                names.append(''.join(buff))
                buff.clear()
                continue
            names.append(pair)
        return names

    @staticmethod
    def loc_urn(name: str) -> str:
        return URN.any_urn(name, LOCAL_NODE_ID)

    @staticmethod
    def any_urn(name: str, node_id: str) -> str:
        urn = URN()
        urn.domain = MESH_DOMAIN
        urn.node_id = node_id
        urn.flag = URNFlag.parse("0001000000000000000000000000000000000")
        urn.name = name
        return urn.string()


class URNFlag:

    @staticmethod
    def parse(value: str):
        flags = URNFlag()
        if value is None or value.isspace():
            log.error(f'Unresolved Flag {value}')
            return flags
        flags.v = value[0:2]
        flags.proto = value[2:4]
        flags.codec = value[4:6]
        flags.version = URNFlag.reduce(value[6:12], 2)
        flags.zone = value[12:14]
        flags.cluster = value[14:16]
        flags.cell = value[16:18]
        flags.group = value[18:20]
        flags.address = URNFlag.reduce(value[20:32], 3)
        flags.port = URNFlag.reduce(value[32:37], 5)
        return flags

    @staticmethod
    def reduce(value: str, length: int):
        buff = []
        for index, code in enumerate(value):
            if (index + length) % length != 0:
                continue
            has_none_zero = True
            for offset in range(index, index + length):
                if value[offset] != '0':
                    has_none_zero = False
                if not has_none_zero:
                    buff.append(value[offset])
            if has_none_zero:
                buff.append('0')
            buff.append('.')

        return ''.join(buff[0:len(buff) - 1])

    def __init__(self):
        self.v = '00'  # #                 2      00
        self.proto = '00'  # #             2      00
        self.codec = '00'  # #             2      00
        self.version = '000000'  # #       6      000000
        self.zone = '00'  # #              2      00
        self.cluster = '00'  # #           2      00
        self.cell = '00'  # #              2      00
        self.group = '00'  # #             2      00
        self.address = '000000000000'  # # 12     000000000000
        self.port = '00080'  # #           5      00080   Max port 65535

    @staticmethod
    def padding(v: str, length: int) -> str:
        value = v.replace(".", "")
        if len(value) == length:
            return value
        if len(value) < length:
            return f'{"".join(repeat("0", length - len(v)))}{value}'
        return v[0:length]

    def padding_chain(self, v: str, length: int, size: int) -> str:
        def chain(values: [str]) -> str:
            chars = []
            for frag in values:
                chars.append(self.padding(frag, length))
            return "".join(chars)

        frags = v.split(".")
        minimal = len(frags)
        if len(frags) == size:
            return chain(frags)
        if len(frags) < size:
            for _ in range(0, size - minimal):
                frags.append("")
            return chain(frags)
        return chain(frags[0:size])

    def string(self):
        return "{}{}{}{}{}{}{}{}{}{}".format(
            self.padding(self.v, 2),
            self.padding(self.proto, 2),
            self.padding(self.codec, 2),
            self.padding_chain(self.version, 2, 3),
            self.padding(self.zone, 2),
            self.padding(self.cluster, 2),
            self.padding(self.cell, 2),
            self.padding(self.group, 2),
            self.padding_chain(self.address, 3, 4),
            self.padding(self.port, 5)
        )

    def reset(self, proto: str, codec: str, version: str, zone: str, cluster: str, cell: str, group: str, address: str,
              port: str) -> "URNFlag":
        self.proto = proto
        self.codec = codec
        self.version = version
        self.zone = zone
        self.cluster = cluster
        self.cell = cell
        self.group = group
        self.address = address
        self.port = port
        return self
