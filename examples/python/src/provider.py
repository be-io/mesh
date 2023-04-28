#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

from consumer import Consumer
from mesh import mps


@mps
class Provider(Consumer):

    def who(self) -> str:
        return "I am terminator"
