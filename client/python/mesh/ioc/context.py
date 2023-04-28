#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from abc import abstractmethod
from typing import Any


class Context:

    @abstractmethod
    def inject_properties(self, properties: bytes):
        pass

    @abstractmethod
    def register_processor(self, processor: Any):
        pass
