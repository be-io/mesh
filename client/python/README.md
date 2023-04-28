# Mesh Python Client

[![Build Status](https://travis-ci.org/ducesoft/babel.svg?branch=master)](https://travis-ci.org/ducesoft/babel)
[![Financial Contributors on Open Collective](https://opencollective.com/babel/all/badge.svg?label=financial+contributors)](https://opencollective.com/babel) [![codecov](https://codecov.io/gh/babel/babel/branch/master/graph/badge.svg)](https://codecov.io/gh/babel/babel)
![license](https://img.shields.io/github/license/ducesoft/babel.svg)

中文版 [README](README_CN.md)

## Introduction

Mesh is a standard implementation for [Private Transmission Protocol](Specifications.md) specification.

Mesh Python develop kits base on Python3.6. Recommend use [poetry](https://github.com/python-poetry/poetry) to manage
dependencies.

## Features

As an open source Internet of Data infrastructure develop kits, Mesh has the following core functions:

* Minimal kernel with SPI plugin architecture, everything is replacement.
* Support full stack of service mesh architecture.
* Support full stack of service oriented architecture.
* Support transport with TCP, HTTP, or other RPC protocols.
* Support rich routing features.
* Support reliable upstream management and load balancing capabilities.
* Support network and protocol layer observability.
* Support mTLS and protocols on TLS.
* Support rich extension mechanism to provide highly customizable expansion capabilities.
* Support process smooth upgrade.

## Get Started

```bash
poetry add imesh
```

or

```bash
pip install imesh
```

### RPC

Declared rpc interface Facade.

```python

from abc import ABC, abstractmethod

from mesh import spi, mpi


@spi("mesh")
class Tokenizer(ABC):

    @abstractmethod
    @mpi("mesh.trust.apply")
    def apply(self, kind: str, duration: int) -> str:
        """
        Apply a node token.
        :param kind:
        :param duration:
        :return:
        """
        pass

    @abstractmethod
    @mpi("mesh.trust.verify")
    def verify(self, token: str) -> bool:
        """
        Verify some token verifiable.
        :param token:
        :return:
        """
        pass
```

Declared rpc service Implement.

```python

from mesh import mps, Tokenizer


@mps
class MeshTokenizer(Tokenizer):

    def apply(self, kind: str, duration: int) -> str:
        return "foo"

    def verify(self, token: str) -> bool:
        return True
```

Remote reference procedure call.

```python

from mesh import mpi, Tokenizer


class Component:

    @mpi
    def tokenizer(self) -> Tokenizer:
        pass

    def invoke(self) -> bool:
        token = self.tokenizer().apply('PERMIT', 1000 * 60 * 5)
        return self.tokenizer().verify(token)


```

### Transport

Transport is a full duplex communication stream implement.

```python
import mesh
from mesh import Mesh, log, ServiceLoader, Transport, Routable
from mesh.prsim import Metadata


def main():
    mesh.start()

    transport = Routable.of(ServiceLoader.load(Transport).get("mesh"))
    session = transport.with_address("10.99.1.33:570").local().open('session_id_008', {
        Metadata.MESH_VERSION.key(): '',
        Metadata.MESH_TECH_PROVIDER_CODE.key(): 'LX',
        Metadata.MESH_TRACE_ID.key(): Mesh.context().get_trace_id(),
        Metadata.MESH_TOKEN.key(): 'x',
        Metadata.MESH_SESSION_ID.key(): 'session_id_008',
        Metadata.MESH_TARGET_INST_ID.key(): 'JG0100000100000000',
    })
    for index in range(100):
        inbound = f"节点4发送给节点5报文{index}"
        log.info(f"节点4发送:{inbound}")
        session.push(inbound.encode('utf-8'), {}, "topic")
        outbound = session.pop(10000, "topic")
        if outbound:
            log.info(f"节点4接收:{outbound.decode('utf-8')}")

```
