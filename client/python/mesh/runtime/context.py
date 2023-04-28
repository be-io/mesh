#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
from typing import Any

from mesh.ioc import Context, EnvsProcessor, PaaSProcessor
from mesh.kinds.reference import Reference
from mesh.kinds.service import Service


class MeshContext(Context):

    def __init__(self):
        self.env_processors = [EnvsProcessor]
        self.paas_processors = [PaaSProcessor]
        self.properties = []
        self.services = [Service]
        self.references = [Reference]

    def inject_properties(self, properties: bytes):
        self.properties.append(bytes)

    def register_processor(self, processor: Any):
        if issubclass(processor.__class__, EnvsProcessor.__class__):
            self.env_processors.append(processor)
        if issubclass(processor.__class__, PaaSProcessor.__class__):
            self.paas_processors.append(processor)

    def refresh(self):
        services = []
        for processor in self.env_processors:
            processor.post_properties(self)
        for service in self.services:
            for processor in self.paas_processors:
                service = processor.before_initialization(service.kind, service.name)
                service = processor.after_initialization(service, service.name)
                service = processor.before_instantiation(service, service.name)
                service = processor.after_instantiation(service, service.name)
                services = services.append(service)
