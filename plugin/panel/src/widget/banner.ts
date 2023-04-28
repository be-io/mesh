/*
 * Copyright (c) 2000, 2099, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

const banner = (version: string, commitId: string) => {
    return `
  __  __ _____ ____  _   _ 
 |  \\/  | ____/ ___|| | | |
 | |\\/| |  _| \\___ \\| |_| |
 | |  | | |___ ___) |  _  |
 |_|  |_|_____|____/|_| |_|

A lightweight, distributed, relational network architecture for MPC (v${version}, build ${commitId})

`
}

export default banner