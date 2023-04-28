#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from typing import Dict, Any, List

from mesh.kinds import Script, Paging, Page
from mesh.macro import spi
from mesh.mpc import ServiceProxy
from mesh.prsim import Evaluator


@spi("mesh")
class MeshEvaluator(Evaluator):

    def __init__(self):
        self.evaluator = ServiceProxy.default_proxy(Evaluator)

    def compile(self, script: Script) -> str:
        return self.evaluator.compile(script)

    def exec(self, code: str, args: Any, dft: str) -> str:
        return self.evaluator.exec(code, args, dft)

    def dump(self, feature: Dict[str, str]) -> List[Script]:
        return self.evaluator.dump(feature)

    def index(self, index: Paging) -> Page[Script]:
        return self.evaluator.index(index)
