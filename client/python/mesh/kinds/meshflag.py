#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

from enum import Enum


class MeshFlag(Enum):
    HTTP = "00", "http", 0
    GRPC = "01", "grpc", 0
    MQTT = "02", "mqtt", 0
    TCP = "03", "tcp", 0

    JSON = "00", "json", 1
    PROTOBUF = "01", "protobuf", 1
    XML = "02", "xml", 1
    THRIFT = "03", "thrift", 1
    YAML = "04", "yaml", 1

    def get_code(self):
        return self.value[0]

    def get_name(self):
        return self.value[1]

    @staticmethod
    def of_proto(code: str) -> "MeshFlag":
        for member in MeshFlag:
            if member.value and member.value[2] == 0 and member.value[0] == code:
                return member
        return MeshFlag.HTTP

    @staticmethod
    def of_code(code: str) -> "MeshFlag":
        for member in MeshFlag:
            if member.value and member.value[2] == 1 and member.value[0] == code:
                return member
        return MeshFlag.JSON

    @staticmethod
    def of_name(name: str) -> "MeshFlag":
        for member in MeshFlag:
            if member.value and member.value[1] == name:
                return member
        return MeshFlag.HTTP
