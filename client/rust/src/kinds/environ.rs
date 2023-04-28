/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

//! @author coyzeng@gmail.com
//!

pub struct Environ {
    pub version: str,
    pub node_id: str,
    pub inst_id: str,
    pub inst_name: str,
    pub root_crt: str,
    pub root_key: str,
    pub node_crt: str,
}

pub struct Lattice {
    pub zone: str,
    pub cluster: str,
    pub cell: str,
    pub group: str,
    pub address: str,
}
