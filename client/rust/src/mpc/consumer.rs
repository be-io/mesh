/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

//! @author coyzeng@gmail.com
//!

use std::error::Error;
use std::io::Bytes;

use crate::mpc::invoker::Execution;

pub trait Consumer {
    //! Start the mesh broker.
    fn start() -> Result<T, dyn Error>;

    /**
    xxx
     */
    fn consume(urn: &str, execution: dyn Execution, inbound: Bytes<u8>) -> Bytes<u8>;
}