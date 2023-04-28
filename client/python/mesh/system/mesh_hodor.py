#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

from typing import Dict, List

from mesh.macro import spi
from mesh.prsim import Hodor


@spi("mesh")
class MeshHodor(Hodor):

    def stats(self, features: List[str]) -> Dict[str, str]:
        return {"status": "true"}

    def debug(self, features: Dict[str, str]):
        """ debug """
        pass
