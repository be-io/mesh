#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

from abc import ABC, abstractmethod

from mesh.kinds import Event
from mesh.macro import mpi, spi


@spi("mesh")
class Subscriber(ABC):

    @abstractmethod
    @mpi("mesh.queue.subscribe")
    def subscribe(self, event: Event):
        """
        Subscribe the event with {@link com.be.mesh.client.annotate.Bindings} or {@link com.be.mesh.client.annotate.Binding}
        :param event:
        :return:
        """
        pass
