#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#

from typing import Any, List

from mesh.kinds.entity import Entity
from mesh.kinds.principal import Principal
from mesh.macro import ServiceLoader
from mesh.macro import index, serializable, Binding


@serializable
class Topic:

    @index(0)
    def topic(self) -> str:
        return ''

    @index(5)
    def code(self) -> str:
        return ''

    @index(10)
    def group(self) -> str:
        return ''

    @index(15)
    def sets(self) -> str:
        return ''

    def matches(self, bindings: List[Binding]) -> bool:
        if not bindings:
            return False
        for b in bindings:
            if b.topic == self.topic and b.code == self.code:
                return True
        return False


class Event:
    """Any fixed information of principal."""

    @index(0)
    def version(self) -> str:
        return ""

    @index(5)
    def tid(self) -> str:
        return ""

    @index(10)
    def sid(self) -> str:
        return ""

    @index(15)
    def eid(self) -> str:
        return ""

    @index(20)
    def mid(self) -> str:
        return ""

    @index(25)
    def timestamp(self) -> str:
        return ""

    @index(30)
    def source(self) -> Principal:
        return Principal()

    @index(35)
    def target(self) -> Principal:
        return Principal()

    @index(40)
    def binding(self) -> Topic:
        return Topic()

    @index(45)
    def entity(self) -> Entity:
        return Entity()

    @staticmethod
    def new_instance(payload: Any, topic: Topic) -> "Event":
        """
        Create local event instance.
        :param payload:
        :param topic:
        :return:
        """
        from mesh.prsim.network import Network
        network = ServiceLoader.load(Network).get_default()
        target = Principal()
        target.node_id = network.get_environ().node_id
        target.inst_id = network.get_environ().inst_id
        return Event.new_instance_with_target(payload, topic, target)

    @staticmethod
    def new_instance_with_target(payload: Any, topic: Topic, target: Principal) -> "Event":
        """
        Create any node event instance.
        :param payload:
        :param topic:
        :param target:
        :return:
        """
        from mesh.prsim.network import Network
        network = ServiceLoader.load(Network).get_default()
        source = Principal()
        source.node_id = network.get_environ().node_id
        source.inst_id = network.get_environ().inst_id
        return Event.new_instance_with_target_source(payload, topic, target, source)

    @staticmethod
    def new_instance_with_target_source(payload: Any, topic: Topic, target: Principal, source: Principal) -> "Event":
        """
        Create an event with source node and target node.
        :param payload:
        :param topic:
        :param target:
        :param source:
        :return:
        """
        event = Event()
        event.version = "1.0.0"
        event.tid = ""
        event.sid = ""
        event.eid = ""
        event.mid = ""
        event.timestamp = ""
        event.source = source
        event.target = target
        event.binding = topic
        event.entity = Entity.wrap(payload)
        return event
