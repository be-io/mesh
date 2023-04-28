/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

//! @author coyzeng@gmail.com
//!

pub trait Network {
    // Get the meth network environment fixed information.
    #[mpi("mesh.net.environ")]
    fn get_environ(&self) -> Environ;

    // Check the mesh network is accessible.
    #[mpi("mesh.net.accessible")]
    fn accessible(&self, route: Route) -> bool;

    #[mpi("mesh.net.refresh")]
    fn refresh(&self, routes: Vec<Route>) -> void;
}