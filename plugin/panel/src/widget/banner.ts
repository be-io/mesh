/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

const banner = (version: string, commitId: string) => {
    return `
 ██████   ██████ ██████████  █████████  █████   █████
▒▒██████ ██████ ▒▒███▒▒▒▒▒█ ███▒▒▒▒▒███▒▒███   ▒▒███ 
 ▒███▒█████▒███  ▒███  █ ▒ ▒███    ▒▒▒  ▒███    ▒███ 
 ▒███▒▒███ ▒███  ▒██████   ▒▒█████████  ▒███████████ 
 ▒███ ▒▒▒  ▒███  ▒███▒▒█    ▒▒▒▒▒▒▒▒███ ▒███▒▒▒▒▒███ 
 ▒███      ▒███  ▒███ ▒   █ ███    ▒███ ▒███    ▒███ 
 █████     █████ ██████████▒▒█████████  █████   █████
▒▒▒▒▒     ▒▒▒▒▒ ▒▒▒▒▒▒▒▒▒▒  ▒▒▒▒▒▒▒▒▒  ▒▒▒▒▒   ▒▒▒▒▒ 

A lightweight, distributed, relational network architecture for MPC (v${version}, build ${commitId})
`
}

export default banner