#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from typing import Any


class GrpcMarshaller:

    @staticmethod
    def deserialize(buffer: Any):
        return buffer

    @staticmethod
    def serialize(buffer: Any):
        return buffer
