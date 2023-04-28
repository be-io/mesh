#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
#
# import glob
# from os.path import join, dirname, basename, isfileø
# modules = glob.glob(join(dirname(__file__), "*.py"))
# __all__ = [basename(f)[:-3] for f in modules if isfile(f) and not f.endswith('__init__.py')]
from mesh.kinds.commerce import CommerceLicense, CommerceEnviron
from mesh.kinds.document import Document
from mesh.kinds.entity import CacheEntity, Entity
from mesh.kinds.environ import Environ
from mesh.kinds.event import Event, Topic
from mesh.kinds.inbound import Inbound
from mesh.kinds.institution import Institution
from mesh.kinds.license import License
from mesh.kinds.location import Location
from mesh.kinds.meshflag import MeshFlag
from mesh.kinds.outbound import Outbound
from mesh.kinds.page import Page
from mesh.kinds.paging import Paging
from mesh.kinds.principal import Principal
from mesh.kinds.profile import Profile
from mesh.kinds.reference import Reference
from mesh.kinds.registration import Registration, Resource, Binding, MPS, Forward, Proxy
from mesh.kinds.route import Route
from mesh.kinds.script import Script
from mesh.kinds.service import Service
from mesh.kinds.timeout import Timeout
from mesh.kinds.versions import Versions

__all__ = (
    "Environ",
    "Route",
    "Location",
    "Principal",
    "Inbound",
    "Institution",
    "Outbound",
    "Reference",
    "Service",
    "Event",
    "Topic",
    "Entity",
    "MeshFlag",
    "CacheEntity",
    "Timeout",
    "Registration",
    "Versions",
    "Profile",
    "License",
    "CommerceLicense",
    "CommerceEnviron",
    "Document",
    "Paging",
    "Page",
    "Resource",
    "Binding",
    "MPS",
    "Forward",
    "Proxy",
    "Script",
)


def init():
    """ init function """
    pass
