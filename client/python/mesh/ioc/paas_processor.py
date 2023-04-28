#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from abc import abstractmethod
from typing import Any


class PaaSProcessor:

    @abstractmethod
    def before_initialization(self, cls: type, name: str) -> Any:
        pass

    @abstractmethod
    def after_initialization(self, instance, name: str) -> Any:
        pass

    @abstractmethod
    def before_instantiation(self, instance, name: str) -> Any:
        pass

    @abstractmethod
    def after_instantiation(self, instance, name: str) -> Any:
        pass
