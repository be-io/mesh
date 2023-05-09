/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */


import icons from "@/widget/icons";
import {SetOption} from "@/widget/block";


const Sets: SetOption[] = [{
    name: 'Mesh',
    icon: icons.MeshIcon,
    color: 'success',
    describe: 'A lightweight, distributed, relational network architecture for MPC',
    replicas: 1,
    status: '1'
}, {
    name: 'Mysql',
    icon: icons.MeshIcon,
    color: 'success',
    describe: 'One MySQL Database service for OLTP, OLAP, and ML.',
    replicas: 1,
    status: '1'
}, {
    name: 'Redis',
    icon: icons.RedisIcon,
    color: 'success',
    describe: 'Redis is an in-memory database that persists on disk. The data model is key-value, but many different kind of values are supported: Strings, Lists, Sets, Sorted Sets, Hashes, Streams, HyperLogLogs, Bitmaps.',
    replicas: 1,
    status: '1'
}];

export default Sets;