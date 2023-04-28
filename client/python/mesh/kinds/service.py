#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from typing import Dict

from mesh.macro import index, serializable


@serializable
class Service:
    """
    Service definition.
            service = Service()
            service.urn = ""
            service.namespace = ""
            service.name = ""
            service.version = version
            service.proto = proto
            service.codec = codec
            service.asyncable = asyncable
            service.timeout = 3000
            service.retries = 3
            service.node = ""
            service.inst = ""
            service.zone = ""
            service.cluster = ""
            service.cell = ""
            service.group = ""
            service.address = ""
            service.kind = ""
            service.lang = ""
            service.attrs = {}
    """

    @index(0)
    def urn(self) -> str:
        """
        Uniform domain name.
        """
        pass

    @index(5)
    def namespace(self) -> str:
        """
        Service namespace.
        """
        pass

    @index(10)
    def name(self) -> str:
        """
        Service name.
        """
        pass

    @index(15)
    def version(self) -> str:
        """
        Service version.
        """
        pass

    @index(20)
    def proto(self) -> str:
        """
        Service net protocol.
        """
        pass

    @index(25)
    def codec(self) -> str:
        """
        Service payload codec protocol.
        """
        pass

    @index(30)
    def flags(self) -> int:
        """
        Service serve asyncable.
        """
        pass

    @index(35)
    def timeout(self) -> int:
        """
        Service serve timeout in mill seconds.
        """
        pass

    @index(40)
    def retries(self) -> int:
        """
        Service serve retry times.
        """
        pass

    @index(45)
    def node(self) -> str:
        """
        Service customized zone.
        """
        pass

    @index(50)
    def inst(self) -> str:
        """
        Service customized institution.
        """
        pass

    @index(55)
    def zone(self) -> str:
        """
        Service customized zone.
        """
        pass

    @index(60)
    def cluster(self) -> str:
        """
        Service customized cluster.
        """
        pass

    @index(65)
    def cell(self) -> str:
        """
        Service customized cell.
        """
        pass

    @index(70)
    def group(self) -> str:
        """
        Service customized group.
        """
        pass

    @index(75)
    def group(self) -> str:
        """
        Service customized sets.
        """
        pass

    @index(80)
    def address(self) -> str:
        """
        Service customized address.
        """
        pass

    @index(85)
    def kind(self) -> str:
        """
        Service kind.
        """
        pass

    @index(90)
    def lang(self) -> str:
        """
        Service lang.
        """
        pass

    @index(95)
    def attrs(self) -> Dict[str, str]:
        """
        Service customized attrs.
        """
        pass
