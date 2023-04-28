#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from typing import Any

import mesh.log as log
from mesh.macro import spi
from mesh.mpc import Provider


@spi("http")
class HTTPProvider(Provider):

    def start(self, address: str, tc: Any):
        log.info(f"Listening and serving HTTP 1.x on {address}")

    def close(self):
        log.info(f"Graceful stop HTTP 1.x serving")
