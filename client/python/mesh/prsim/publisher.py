#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from abc import ABC, abstractmethod
from typing import List, Any

from mesh.kinds.event import Event
from mesh.kinds.event import Topic
from mesh.kinds.principal import Principal
from mesh.macro import mpi, spi


@spi("mesh")
class Publisher(ABC):

    @abstractmethod
    @mpi("mesh.queue.publish")
    def publish(self, events: List[Event]) -> List[str]:
        """
        Publish event to mesh.
        :param events:
        :return:
        """
        pass

    @abstractmethod
    @mpi("mesh.queue.broadcast")
    def broadcast(self, events: List[Event]) -> List[str]:
        """
        Synchronized broadcast the event to all subscriber. This maybe timeout with to many subscriber.
        :param events: Event payload
        :return: Synchronized subscriber return value
        """
        pass

    def publish_with_topic(self, binding: Topic, payload: Any) -> str:
        """
        Publish message to local node.
        :param binding:
        :param payload:
        :return:
        """
        event = Event.new_instance(payload, binding)
        return self.publish([event])[0]

    def unicast(self, binding: Topic, payload: Any, principal: Principal) -> str:
        """
        Unicast will publish to another node.
        :param binding:
        :param payload:
        :param principal:
        :return:
        """
        event = Event.new_instance_with_target(payload, binding, principal)
        return self.publish([event])[0]

    def multicast(self, binding: Topic, payload: Any, principals: List[Principal]) -> List[str]:
        """
        Multicast will publish event to principal groups.
        :param binding:
        :param payload:
        :param principals:
        :return:
        """
        events = []
        for principal in principals:
            events.append(Event.new_instance_with_target(payload, binding, principal))
        return self.publish(events)

    def broadcast_with_topic(self, binding: Topic, payload: Any) -> List[str]:
        """
        Synchronized broadcast the event to all subscriber. This maybe timeout with to many subscriber.
        :param binding:
        :param payload:
        :return:
        """
        return self.broadcast([Event.new_instance(payload, binding)])
