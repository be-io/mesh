#
# Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
# TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
#
#
import inspect
import os
from typing import Type, Any, Generic, Dict, List, Callable

import mesh.kinds as kinds
import mesh.log as log
import mesh.tool as tool
from mesh.kinds import Reference, Service, Environ, MeshFlag, Event, Registration, Resource, Topic, Versions
from mesh.macro import ark, T, A, mpi, MPI, spi, mps, MPS, binding, Binding, ServiceLoader, Returns, Parameters, \
    Inspector
from mesh.mpc.compiler import Compiler
from mesh.mpc.eden import Eden
from mesh.mpc.inspector import MethodInspector
from mesh.mpc.invoker import Invocation, Execution
from mesh.mpc.service import ServiceInvokeHandler
from mesh.mpc.service_proxy import ServiceProxy
from mesh.mpc.urn import URN, URNFlag, MESH_DOMAIN
from mesh.prsim import Network, RuntimeHook, Subscriber, Registry, Scheduler, Builtin


class Provider:

    def __init__(self, kind: Type[T], service: Any):
        self.making: bool = False
        self.kind: Type[T] = kind
        self.service: Any = service


class Consumer:
    def __init__(self, metadata: mpi, reference: Type[T]):
        self.metadata: mpi = metadata
        self.reference: Type[T] = reference
        self.proxy: Any = ServiceProxy.static_proxy(reference, metadata)


class Instance(Generic[T], Execution[T]):

    def __init__(self, urn: str, kind: Type[T], method: Any, target: Any, resource: T):
        declared_method = getattr(kind, method.__name__)
        self.urn = urn
        self.kind = kind
        self.method = method
        self.target = target
        self.resource = resource
        self.inspector = MethodInspector(MPI.get_mpi_if_present(method), kind, method)
        self.invoker = ServiceInvokeHandler(target)
        self.__intype = ServiceLoader.load(Compiler).get_default().intype(declared_method)
        self.__retype = ServiceLoader.load(Compiler).get_default().retype(declared_method)

    def schema(self) -> T:
        return self.resource

    def inspect(self) -> Inspector:
        return self.inspector

    def intype(self) -> Type[Parameters]:
        return self.__intype

    def retype(self) -> Type[Returns]:
        return self.__retype

    def inflect(self) -> Parameters:
        return self.__intype()

    def reflect(self) -> Returns:
        return self.__retype()

    def run(self, invocation: Invocation) -> Any:
        return self.invoker.run(invocation)


@spi("mesh")
@binding(topic="mesh.registry.event", code="refresh", meshable=False)
class MeshEden(Eden, RuntimeHook, Subscriber):
    """
    Like: create.tenant.omega.json.http2.lx000001.mpi.trustbe.net
    """

    def __init__(self):
        self.topic = Topic()
        self.topic.topic = "mesh.registry.event"
        self.topic.code = "refresh"
        self.task_id = ""
        self.providers: Dict[Type[T], Provider] = {}
        self.consumers: Dict[Type[T], Dict[mpi, Consumer]] = {}
        self.indies: Dict[Type[T], Dict[mpi, Dict[Any, Instance[Reference]]]] = {}
        self.services: Dict[str, Instance[Service]] = {}
        self.references: Dict[str, Instance[Reference]] = {}

    def define(self, metadata: mpi, reference: Type[T]) -> T:
        return self.make_consumer(metadata, reference)

    def refer(self, metadata: mpi, reference: Type[T], method: Inspector) -> Execution[Reference]:
        return self.make_reference_exec(metadata, reference, self.get_environ()).get(method.get_signature())

    def store(self, kind: Type[T], service: Any):
        if not self.providers.get(kind):
            self.providers[kind] = Provider(kind, service)
        return self.providers.get(kind)

    def infer(self, urn: str) -> Execution[Service]:
        return self.make_service_exec(self.get_environ()).get(URN.parse(urn).name)

    def refer_types(self) -> List[Type[T]]:
        return list(self.consumers.keys())

    def infer_types(self) -> List[Type[T]]:
        types: List[Type[T]] = []
        for provider in self.providers.values():
            types.append(provider.kind)
        return types

    def make_consumer(self, metadata: mpi, reference: Type[T]):
        if not self.consumers.get(reference):
            self.consumers[reference] = {}
        if not self.consumers.get(reference).get(metadata):
            self.consumers[reference][metadata] = Consumer(metadata, reference)
        return self.consumers.get(reference).get(metadata)

    def make_reference_exec(self, md: mpi, ref: Type[T], env: Environ) -> Dict[Any, Instance[Reference]]:
        consumer = self.make_consumer(md, ref)
        if not self.indies.get(ref):
            self.indies[ref] = {}
        if not self.indies.get(ref).get(md):
            executions: Dict[Any, Instance[Reference]] = {}
            for kind, methods in tool.get_declared_methods(ref).items():
                for method in methods:
                    refer = self.make_method_as_reference(env, md, kind, method)
                    instance = Instance(refer.urn, ref, method, consumer.proxy, refer)
                    self.references[refer.urn] = instance
                    executions[instance.inspect().get_signature()] = instance
            self.indies.get(ref)[md] = executions
        return self.indies[ref][md]

    def make_service_exec(self, environ: Environ) -> Dict[str, Instance[Service]]:
        for provider in self.providers.values():
            if provider.making:
                continue
            for kind, methods in tool.get_declared_methods(provider.kind).items():
                for method in methods:
                    if not hasattr(provider.service, method.__name__):
                        continue
                    mm = getattr(provider.service, method.__name__)
                    schemas: List[Service] = self.make_schema(environ, provider.kind, kind, method)
                    for schema in schemas:
                        urn = URN.parse(schema.urn)
                        self.services[urn.name] = Instance(schema.urn, provider.kind, mm, provider.service, schema)
            provider.making = True
        return self.services

    @staticmethod
    def get_environ() -> Environ:
        network = ServiceLoader.load(Network).get_default()
        return network.get_environ()

    def make_schema(self, environ: Environ, interface: T, kind: T, method: Any) -> List[Service]:
        schemas = []
        macro = MPS.get_mps_if_present(interface)
        if macro:
            schemas.append(self.make_mps_as_service(macro, kind, method, environ))
        bindings = Binding.get_binding_if_present(interface)
        for metadata in bindings:
            if metadata.meshable:
                schemas = schemas.append(self.make_binding_as_service(metadata, kind, method, environ))
        return schemas

    def make_method_as_reference(self, environ: Environ, metadata: MPI, interface: T, method: Any) -> Reference:
        om = MPI.get_mpi_if_present(method)
        ok = MPI.get_mpi_if_present(interface)
        name = self.choose_if_present(lambda x: x.name, "", om)
        alias = tool.anyone(name, str.lower(f'{interface.__name__}.{method.__name__}'))
        reference = Reference()
        reference.namespace = interface.__name__
        reference.name = method.__name__
        reference.version = self.choose_if_present(lambda x: x.version, '', metadata, om, ok)
        reference.proto = self.choose_if_present(lambda x: x.proto, MeshFlag.GRPC.name, metadata, om, ok)
        reference.codec = self.choose_if_present(lambda x: x.codec, MeshFlag.JSON.name, metadata, om, ok)
        reference.flags = self.choose_if_present(lambda x: x.flags, 0, om, ok) | metadata.flags
        reference.timeout = self.choose(lambda x: x.timeout, 10000, metadata, om, ok)
        reference.retries = self.choose(lambda x: x.retries, 3, metadata, om, ok)
        reference.node = self.choose_if_present(lambda x: x.node, "", metadata, om, ok)
        reference.inst = self.choose_if_present(lambda x: x.inst, "", metadata, om, ok)
        reference.zone = self.choose_if_present(lambda x: x.zone, "", metadata, om, ok)
        reference.cluster = self.choose_if_present(lambda x: x.cluster, "", metadata, om, ok)
        reference.cell = self.choose_if_present(lambda x: x.cell, "", metadata, om, ok)
        reference.group = self.choose_if_present(lambda x: x.group, "", metadata, om, ok)
        reference.address = self.choose_if_present(lambda x: x.address, "", metadata, om, ok)
        reference.urn = self.get_urn(alias, self.get_reference_urn_flag(reference), environ)
        return reference

    def make_mps_as_service(self, metadata: MPS, kind: T, method: Any, environ: Environ) -> Service:
        om = MPI.get_mpi_if_present(method)
        ok = MPI.get_mpi_if_present(kind)
        name = self.choose_if_present(lambda x: x.name, "", om).replace("${mesh.name}", tool.get_mesh_name()). \
            replace("${mesh.uname}", metadata.name)
        alias = tool.anyone(name, str.lower(f'{kind.__name__}.{method.__name__}'))
        service = self.init_service(environ)
        service.kind = kinds.MPS
        service.namespace = kind.__name__
        service.name = method.__name__
        service.proto = self.choose_if_present(lambda x: x.proto, MeshFlag.GRPC.name, metadata, om, ok)
        service.codec = self.choose_if_present(lambda x: x.codec, MeshFlag.JSON.name, metadata, om, ok)
        service.version = self.choose_if_present(lambda x: x.version, metadata.version, metadata, om, ok)
        service.flags = self.choose_if_present(lambda x: x.flags, 0, om, ok) | metadata.flags
        service.timeout = self.choose(lambda x: x.timeout, metadata.timeout, metadata, om, ok)
        service.urn = self.get_urn(alias, self.get_service_urn_flag(service), environ)
        return service

    def make_binding_as_service(self, metadata: Binding, kind: T, method: Any, environ: Environ) -> Service:
        service = self.init_service(environ)
        service.kind = kinds.Binding
        service.namespace = metadata.topic
        service.name = metadata.code
        service.proto = self.choose_if_present(lambda x: x.proto, MeshFlag.GRPC.name, metadata)
        service.codec = self.choose_if_present(lambda x: x.codec, MeshFlag.JSON.name, metadata)
        service.version = metadata.version
        service.flags = metadata.flags
        service.timeout = self.choose(lambda x: x.timeout, metadata.timeout, metadata)
        service.urn = self.get_urn(str.lower(f'{metadata.topic}.{metadata.code}'), self.get_service_urn_flag(service),
                                   environ)
        return service

    def init_service(self, environ: Environ) -> Service:
        service = Service()
        service.namespace = ''
        service.name = ''
        service.version = ''
        service.proto = ''
        service.codec = ''
        service.flags = 0
        service.timeout = 10000
        service.retries = 0
        service.node = environ.node_id
        service.inst = environ.inst_id
        service.zone = environ.lattice.zone
        service.cluster = environ.lattice.cluster
        service.cell = environ.lattice.cell
        service.group = environ.lattice.group
        service.address = self.get_listened_address()
        service.kind = kinds.MPS
        service.lang = 'Python3'
        service.attrs = {}
        service.urn = ''
        return service

    @staticmethod
    def get_urn(alias: str, flag: URNFlag, environ: Environ) -> str:
        urn = URN()
        urn.domain = MESH_DOMAIN
        urn.node_id = str.lower(environ.node_id)
        urn.name = alias
        urn.flag = flag
        return urn.string()

    @staticmethod
    def get_reference_urn_flag(reference: Reference) -> URNFlag:
        authority = str.split(f'{reference.address}:', ':')
        flag = URNFlag()
        return flag.reset(MeshFlag.of_name(reference.proto).get_code(),
                          MeshFlag.of_name(reference.codec).get_code(),
                          reference.version,
                          reference.zone,
                          reference.cluster,
                          reference.cell,
                          reference.group,
                          authority[0],
                          authority[1])

    @staticmethod
    def get_service_urn_flag(service: Service) -> URNFlag:
        flag = URNFlag()
        return flag.reset(MeshFlag.of_name(service.proto).get_code(),
                          MeshFlag.of_name(service.codec).get_code(),
                          service.version,
                          service.zone,
                          service.cluster,
                          service.cell,
                          service.group,
                          '',
                          '')

    @staticmethod
    def choose_if_present(fn: Callable[[A], T], defaults: T, *tags: A) -> T:
        for tag in tags:
            if tag and fn(tag):
                return fn(tag)

        return defaults

    @staticmethod
    def choose(fn: Callable[[A], T], defaults: T, *tags: A) -> T:
        for tag in tags:
            if tag and fn(tag) and fn(tag) != defaults:
                return fn(tag)

        return defaults

    @staticmethod
    def get_listened_address() -> str:
        return f'{tool.get_mesh_runtime().hostname}:{tool.get_mesh_runtime().port}'

    def new_registration(self, metadata: Resource, attachments: Dict[str, str]) -> Registration:
        registration = Registration()
        registration.instance_id = self.get_listened_address()
        registration.content = metadata
        registration.kind = Registration.METADATA
        registration.address = self.get_listened_address()
        registration.name = tool.get_mesh_name()
        registration.timestamp = 30 * 1000
        registration.attachments = attachments
        return registration

    @mpi
    def registry(self) -> Registry:
        """ """
        pass

    @spi("mesh")
    def scheduler(self) -> Scheduler:
        return ServiceLoader.load(Scheduler).get_default()

    @spi("mesh")
    def builtin(self) -> Builtin:
        return ServiceLoader.load(Builtin).get_default()

    def start(self):
        for kind, nss in ark.export(mps).items():
            for spec in ark.providers(mps, kind):
                if not isinstance(spec.metadata, MPS) or inspect.isabstract(spec.kind):
                    continue
                try:
                    self.store(kind, spec.kind())
                    break
                except BaseException as e:
                    log.warn(f"Mesh service {kind.__name__} with {e.__str__()}")

        if os.environ.get("MESH_REGISTER", "0") != "1":
            return

        self.task_id = self.scheduler().period(30 * 1000, self.topic)

    def stop(self):
        self.scheduler().cancel(self.task_id)
        self.scheduler().shutdown(3 * 1000)
        metadata = Resource()
        metadata.services = []
        metadata.references = []
        try:
            self.registry().unregister(self.new_registration(metadata, {}))
        except BaseException as e:
            log.warn("Graceful unregister metadata, {}", str(e))

    def refresh(self):
        self.start()

    def subscribe(self, event: Event):
        v = self.builtin().version()
        version = v if v else Versions()
        attachments = {}
        for k, v in (version.infos if version.infos else {}).items():
            attachments[k] = v
        environ = self.get_environ()
        self.make_service_exec(environ)
        for _, consumers in self.consumers.items():
            for _, consumer in consumers.items():
                self.make_reference_exec(consumer.metadata, consumer.reference, environ)
        metadata = Resource()
        metadata.services = []
        for _, reference in self.references.items():
            metadata.services.append(reference.resource)
        metadata.references = []
        for _, service in self.services.items():
            metadata.references.append(service.resource)
        self.registry().register(self.new_registration(metadata, attachments))
