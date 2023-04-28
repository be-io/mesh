#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from typing import List

from mesh.kinds import Topic, Timeout
from mesh.macro import spi
from mesh.mpc import ServiceProxy
from mesh.prsim import Scheduler


@spi("mesh")
class MeshScheduler(Scheduler):

    def __init__(self):
        self.remote = ServiceProxy.default_proxy(Scheduler)

    def timeout(self, timeout: Timeout, duration: int) -> str:
        return self.remote.timeout(timeout, duration)

    def cron(self, cron: str, binding: Topic) -> str:
        return self.remote.cron(cron, binding)

    def period(self, duration: int, binding: Topic) -> str:
        return self.remote.period(duration, binding)

    def dump(self) -> List[str]:
        return self.remote.dump()

    def cancel(self, task_id: str) -> bool:
        return self.remote.cancel(task_id)

    def stop(self, task_id: str) -> bool:
        return self.remote.stop(task_id)

    def emit(self, topic: Topic) -> bool:
        return self.remote.emit(topic)

    def shutdown(self, duration: int) -> bool:
        return self.remote.shutdown(duration)
