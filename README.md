![logo](./LOGO.svg)

# Mesh: A lightweight, distributed, relational network architecture for MPC

[![Financial Contributors on Open Collective](https://opencollective.com/babel/all/badge.svg?label=financial+contributors)](https://opencollective.com/babel) [![codecov](https://codecov.io/gh/babel/babel/branch/master/graph/badge.svg)](https://codecov.io/gh/babel/babel)
![license](https://img.shields.io/github/license/ducesoft/babel.svg)

Private Transmission Protocol -- PTP/1.0 [README](specifications/api/Specifications_CN.md)

中文版 [README](README_CN.md)

高可用 [HA](HA.md)

## Introduction

Mesh is an MPC network solution to serve resource as service.

## Features

As an open source network proxy, Mesh has the following core functions:

* Support full dynamic resource configuration through xDS API integrated with Service Mesh.
* Support proxy with TCP, HTTP, and RPC protocols.
* Support rich routing features.
* Support reliable upstream management and load balancing capabilities.
* Support network and protocol layer observability.
* Support mTLS and protocols on TLS.
* Support rich extension mechanism to provide highly customizable expansion capabilities.
* Support process smooth upgrade.

## Get Started

Install go compile time.

```bash
brew install go
```

Set the go dependency proxy host.

```bash
go env -w GOPROXY="https://mirrors.aliyun.com/goproxy/,https://goproxy.cn,https://goproxy.io,https://proxy.golang.org,direct"
```

Compile mesh project.

```bash
go mod tidy
```

Test project.

```bash
make test
```

Build distribute package for multiplatform.

```bash
make build
```

Generate software develop kits for C++, Python, Java, Golang language.

```bash
make babel
```

Release packages.

```bash
make publish
```

## Documentation

Mesh default listened port list. Mesh server boot with data port 7700 and control port 7305. These port communication
from node to node.

* 570 Mesh port in node
* 7304 Mesh port out node
* 80/443 Mesh port in panel

### Cli usage

```bash
# Container
docker run -d -v 7304:7304 bfia/mesh:latest

# MacOS
brew install mesh

# CentOS
yum install mesh -y

# Alpine
apk add mesh
```

Mesh man usage

```bash
ducer@coyzeng mesh % mesh --help

Usage:
  Mesh [flags]
  Mesh [command]

Examples:
mesh COMMAND

Available Commands:
  access      Mesh data access commands.
  babel       Mesh babel.
  completion  Generate the autocompletion script for the specified shell
  domain      Remove mesh net domains.
  dump        Dump the mesh cluster metadata tables.
  exec        Execute the mesh actor, Such mesh exec actor -i '{}'.
  help        Help about any command
  http        Mesh curl.
  inspect     Return low-level information on mesh plugin objects.
  issue       Mesh issue.
  kv          Mesh kv store.
  net         Display a live stream of mesh cluster resource usage statistics.
  os          Mesh OS 🚀🚀🚀.
  start       Start mesh with the run mode(server/operator/node/panel).
  status      Display a live stream of mesh cluster resource usage statistics.
  stop        Stop the mesh process if available.
  token       Mesh token center.

Flags:
  -c, --config string   Mesh configuration input, it can be url/path/body, default is {}. (default "{}")
  -d, --debug           Set the run mode to debug with more information.
  -f, --format string   Mesh configuration format, it can be json/yaml, default is yaml. (default "yaml")
  -h, --help            help for Mesh
  -v, --version         version for Mesh

Use "Mesh [command] --help" for more information about a command.
```

### Client usage

Golang

```bash
go get github.com/be-io/mesh/client/golang
```

Java

```xml

<dependency>
    <groupId>io.be.mesh</groupId>
    <artfactId>mesh</artfactId>
    <version>0.0.17</version>
</dependency>
```

Python

```bash
poetry add imesh==0.0.17

or

pip install imesh==0.0.17
```

Typescript

```bash
yarn add imesh
```

Rust

```bash
cargo add imesh
```

### Setting node information

```bash
curl --location 'http://10.99.28.33:7304/v1/interconn/node/refresh' --header 'Content-Type: application/json' --data '{
    "version": "1.0.0",
    "node_id": "LX0000010000280",
    "inst_id": "JG0100002800000000",
    "inst_name": "互联互通节点X",
    "root_crt": "-----BEGIN CERTIFICATE-----\nMIIEZTCCA02gAwIBAgIRAKkfFud6oGcbEmGOZ2i917wwDQYJKoZIhvcNAQELBQAw\ngZIxCzAJBgNVBAYTAkNOMQswCQYDVQQIEwJaSjELMAkGA1UEBxMCSFoxJDAiBgNV\nBAoTG0xYMDAwMDAxMDAwMDI3MC50cnVzdGJlLm5ldDEkMCIGA1UECxMbTFgwMDAw\nMDEwMDAwMjcwLnRydXN0YmUubmV0MR0wGwYDVQQDDBTok53osaHpm4blm6IyN+WF\nrOWPuDAgFw0yMjA5MjEwNjU1MTRaGA8yMTIyMDkyMTA2NTUxNFowgZIxCzAJBgNV\nBAYTAkNOMQswCQYDVQQIEwJaSjELMAkGA1UEBxMCSFoxJDAiBgNVBAoTG0xYMDAw\nMDAxMDAwMDI3MC50cnVzdGJlLm5ldDEkMCIGA1UECxMbTFgwMDAwMDEwMDAwMjcw\nLnRydXN0YmUubmV0MR0wGwYDVQQDDBTok53osaHpm4blm6IyN+WFrOWPuDCCASIw\nDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAOy3fx+vyvaRBTg9hXSXW/XEz5A9\nazavO4KNMPB/Rp6UtQ9+ilireBQXqCM1Zzo5/cZ+FoRmQp+Rl6b1HYAnRaPSTNFM\n4JDTl3AczsLwJcvFlqobqap9/3k5ZEfEEKJp0LFcl3GR2Ral7Uhsmgzm87PSZVRX\nFYPgCgfkuWGbeJ9FwCB5ZDovPH77pyJeW0AzPu3JKLk3JnCtUvmXK4HzoxChK4k0\nqJF1DEVSPxG3JLOKQKQqe/fqSTkRfgrU7D7U/TEaEl5lLFPGi6ZT+3AbOUAasWWO\nMTYMGMG6qC2wlQ7WdJJ4KiWhEfA7aGQfzoW85PZyN1QzRyzW1vyXY5MKNXsCAwEA\nAaOBsTCBrjAOBgNVHQ8BAf8EBAMCArwwJwYDVR0lBCAwHgYIKwYBBQUHAwEGCCsG\nAQUFBwMCBggrBgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBSrXjQL\nGNwhXBQxqx5Sh+LQwXK3mzBDBgNVHREEPDA6ghtMWDAwMDAwMTAwMDAyNzAudHJ1\nc3RiZS5uZXSBG0xYMDAwMDAxMDAwMDI3MEB0cnVzdGJlLm5ldDANBgkqhkiG9w0B\nAQsFAAOCAQEAcRhabDb2O55lQfeMkHozqFz3V4fQwQKQzYW6gwf7FEUxPMaxmFrk\nPKK2rE5Li49mUxk/norXMVmEpBQdYgsu2PnY5J39RivFku0nObzIMFkyDPer05tp\nx7sKXS6sIy9gfaN6uZkzHZH+0MD2NlOj19SI5BNJxiwHzB9BZwGwbKtq3DVZzWRq\nVyZnAcXVTq1Pk9gypatvBgD7r6edXSBXEz2d8HMaidJfChweZPVp98uY8s9EzHnJ\nowsia5sPaVNqVFE72IgQ5TUrRQsQIBGOF81i37Bpsc9g/pf2s6eNGW4CV9a8yfN1\n+BLRh/2eCYXUNQOm5IovVXqBJ4MlhchAOQ==\n-----END CERTIFICATE-----\n",
    "root_key": "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA7Ld/H6/K9pEFOD2FdJdb9cTPkD1rNq87go0w8H9GnpS1D36K\nWKt4FBeoIzVnOjn9xn4WhGZCn5GXpvUdgCdFo9JM0UzgkNOXcBzOwvAly8WWqhup\nqn3/eTlkR8QQomnQsVyXcZHZFqXtSGyaDObzs9JlVFcVg+AKB+S5YZt4n0XAIHlk\nOi88fvunIl5bQDM+7ckouTcmcK1S+ZcrgfOjEKEriTSokXUMRVI/Ebcks4pApCp7\n9+pJORF+CtTsPtT9MRoSXmUsU8aLplP7cBs5QBqxZY4xNgwYwbqoLbCVDtZ0kngq\nJaER8DtoZB/Ohbzk9nI3VDNHLNbW/Jdjkwo1ewIDAQABAoIBAQDAhzIu3HTQe/zp\n1CfSPzT9PLixETM9Q+K7+Qgf4vTWEA7/biUpnzTH6sHG+S1fT0FXir/XqbBwRiM5\nGM2IqOhcKLRv2v4e7OmTtup35IhpJui2rE8fquD5gLNOJ2p8HmItjyhhp4UQhZ3r\nNOFKsyDtVacypK2MF9EwwFgCykeeCbVwBj8PxmTunDuVlD/hSTUa9JK/lh3ssGfN\novPPiuL2TTPiZUy415y4Kd4dVvuQW0QIlm4yV+qMynSNgU7p8f6jyKY5R7FAG3MF\nBEO07LZSYzTOhfhd2zz6lQpi5t2srScIUs+6hQSyZlFqILwkKeIgKTYbvb1q+o+s\nk5XijuwBAoGBAP2JjcENSCUj5c1bRiOvczS/XkDrE5YIOwJMFLWjuwuMNZN++PV3\ngZ58cM4leJ2TmpVidfu9Rkcp6q95M5fZa7Ms7qtS8sFu2JZhWBpydToRrsFAApsG\nBU09cfoq7DhlqJyE+p/X0lmmKyTw2HNeWUGU5AjSPeVa/4da9CSjDn97AoGBAO8E\nHevtVODXCJC+VJXOjlvCff4zGQi3pFo5s9hCXUDk8KNmfieJwcLXlqz3MG/kbmHi\n0ywToeBK5wA/npImAFKXGR5U6ivOeTBML8dM0wG+2gVCwazxm3qO1d2jyD6d/o2H\n/pH3rmuj+Hi9QSZPVVp+nrJsrG2HtMkwGFWP+kIBAoGANtZsqafUxeu4xa0LQ6as\nNWl62nG9/8Jx+PI5vHvYdgvyfp+E+5rIl131DDGAoByP3+W2/ScYL0Y6s490gFCP\ngeajDL1ZMktmX0hYxQeioVe3w6azqZIozWcP4vsrspsSWCBPEQmePrO5Ozk4p+Nt\nTMkGdX3700LWaBFdIxt9hEcCgYAf7CvW49bPRMkHE/SWIYVP6hULy2VPjb9ssYI8\novhzf2BIYpr8yuBPFp4wMb+NYjP/7NyJaYHYRAjANr8GA/9NCJM5QtwXx7bV5YcI\nFlGkTQovY7AcWhSK9OLJfGN1QYLLAlvUwQDRrY+1CInYBQaAVKL7b5pD8rkJmdvW\nKamiAQKBgQCq3cjGNqvMBVWeq0qJRsouSXRaHFRTmS9hgrRD/GZ+GQJNuhsDwfV9\nhb74cz7EqaFhQVhLHZ2QP/AauqtZ+p9wIM5X9xJDCaEKGVR+95vvxilBIOqua5IK\nphkbWGfS/Wp0fWGDaBoc05xuMm3Zv7GdiNeKr6soT1dOOCA88L3eTQ==\n-----END RSA PRIVATE KEY-----\n"
}'
```

### Weave partner network

```bash
curl --location 'http://10.99.27.33:7304/v1/interconn/net/refresh' --header 'Content-Type: application/json' --data '{
    "node_id": "LX0000010000280",
    "inst_id": "JG0100002800000000",
    "address": "192.168.100.63:27304",
    "guest_root": "",
    "guest_crt": ""
}'
```