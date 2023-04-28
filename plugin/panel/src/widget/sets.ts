/*
 * Copyright (c) 2000, 2099, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */


import {
    AssetIcon,
    CubeIcon,
    DashboardIcon,
    DataxIcon,
    GrafanaIcon,
    JanusIcon,
    MeshIcon,
    OmegaIcon,
    PandoraIcon,
    RedisIcon,
    ServingIcon,
    ThetaIcon,
    TiKVIcon
} from "@/widget/icons";
import {SetOption} from "@/widget/block";


const Sets: SetOption[] = [{
    name: 'Mesh',
    icon: MeshIcon,
    color: 'success',
    describe: 'A lightweight, distributed, relational network architecture for MPC',
    replicas: 1,
    status: '1'
}, {
    name: 'Janus',
    icon: JanusIcon,
    color: 'success',
    describe: 'Janus is an Integration/Frontend/Openapi gateway.',
    replicas: 1,
    status: '1'
}, {
    name: 'Pandora',
    icon: PandoraIcon,
    color: 'success',
    describe: 'The Control Center of Private Compute.',
    replicas: 1,
    status: '1'
}, {
    name: 'Omega',
    icon: OmegaIcon,
    color: 'success',
    describe: 'The master data system.',
    replicas: 1,
    status: '1'
}, {
    name: 'Asset',
    icon: AssetIcon,
    color: 'success',
    describe: 'Metadata and entity management system.',
    replicas: 1,
    status: '1'
}, {
    name: 'Serving',
    icon: ServingIcon,
    color: 'success',
    describe: 'Realtime private compute and data scene develop platform.',
    replicas: 1,
    status: '1'
}, {
    name: 'Datax',
    icon: DataxIcon,
    color: 'success',
    describe: 'DataX connect MySQL、Oracle、OceanBase、SqlServer、Postgre、HDFS、Hive、ADS、HBase、TableStore(OTS)、MaxCompute(ODPS)、Hologres、DRDS etc.',
    replicas: 1,
    status: '1'
}, {
    name: 'Theta',
    icon: ThetaIcon,
    color: 'success',
    describe: 'Federal Learning product frontend.',
    replicas: 1,
    status: '1'
}, {
    name: 'Cube',
    icon: CubeIcon,
    color: 'success',
    describe: 'Federal Learning backend engine.',
    replicas: 1,
    status: '1'
}, {
    name: 'Loki',
    icon: GrafanaIcon,
    color: 'success',
    describe: 'Like Prometheus, but for logs.',
    replicas: 1,
    status: '1'
}, {
    name: 'IData',
    icon: TiKVIcon,
    color: 'success',
    describe: 'Distributed transactional key-value database.',
    replicas: 1,
    status: '1'
}, {
    name: 'Loghub',
    icon: DashboardIcon,
    color: 'success',
    describe: 'Log center.',
    replicas: 1,
    status: '1'
}, {
    name: 'Mysql',
    icon: MeshIcon,
    color: 'success',
    describe: 'One MySQL Database service for OLTP, OLAP, and ML.',
    replicas: 1,
    status: '1'
}, {
    name: 'Redis',
    icon: RedisIcon,
    color: 'success',
    describe: 'Redis is an in-memory database that persists on disk. The data model is key-value, but many different kind of values are supported: Strings, Lists, Sets, Sorted Sets, Hashes, Streams, HyperLogLogs, Bitmaps.',
    replicas: 1,
    status: '1'
}];

export default Sets;