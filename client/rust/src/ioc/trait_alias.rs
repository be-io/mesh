/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

//! @author coyzeng@gmail.com
//!
//! Adapted from https://stackoverflow.com/a/30293051/3267834

// FIXME: Use real trait aliases when they are stabilized:
//        https://github.com/rust-lang/rust/issues/41517
macro_rules! trait_alias {
    ($(#[$attributes:meta])* $visibility:vis $name:ident = $base1:ident $(+ $base2:ident)*) => {
        $(#[$attributes])*
        $visibility trait $name: $base1 $(+ $base2)* { }
        impl<T: $base1 $(+ $base2)*> $name for T { }
    };
}
