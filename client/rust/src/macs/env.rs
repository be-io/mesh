/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

//! @author coyzeng@gmail.com

macro_rules! varargs {
    ($ui:expr, $f:path, ($($p:expr),*), $content:expr) => {
        $f($ui, $($p,)* $content);
    };
}

pub struct Environ {
    pub x(&self)->&str{

    }
}

#[varargs]
fn env(backoff: &str, keys: &str) -> &str {

}

