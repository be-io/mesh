/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

//! @author coyzeng@gmail.com
//!
//! This module contains trait definitions for provided services and interfaces

use std::error::Error;

use crate::ioc::module::{Module, ModuleInterface};

/// Like [`Component`]s, providers provide a service by implementing an interface.
///
/// Unlike [`Component`], `Provider` represents a temporary service. Examples include a connection
/// to a remote service or pooled database connection. Because only providers can have other
/// providers as dependencies, services which use these provided services must also be `Provider`s
/// (ex. DB repository, service using a DB repository, etc).
///
/// See also the [provider getting started guide].
///
/// [`Component`]: trait.Component.html
/// [provider getting started guide]: guide/provider/index.html
pub trait Provider<M: Module>: 'static {
    /// The trait/interface which this provider implements
    type Interface: ?Sized;

    /// Provides the service, possibly resolving other components/providers
    /// to do so.
    fn provide(module: &M) -> Result<Box<Self::Interface>, Box<dyn Error>>;
}

/// The type signature of [`Provider::provide`]. This is used when overriding a
/// provider via [`ModuleBuilder::with_provider_override`]
///
/// [`Provider::provide`]: trait.Provider.html#tymethod.provide
/// [`ModuleBuilder::with_provider_override`]: struct.ModuleBuilder.html#method.with_provider_override
#[cfg(not(feature = "thread_safe"))]
pub type ProviderFn<M, I> = Box<dyn (Fn(&M) -> Result<Box<I>, Box<dyn Error>>)>;
/// The type signature of [`Provider::provide`]. This is used when overriding a
/// provider via [`ModuleBuilder::with_provider_override`]
///
/// [`Provider::provide`]: trait.Provider.html#tymethod.provide
/// [`ModuleBuilder::with_provider_override`]: struct.ModuleBuilder.html#method.with_provider_override
#[cfg(feature = "thread_safe")]
pub type ProviderFn<M, I> = Box<dyn (Fn(&M) -> Result<Box<I>, Box<dyn Error>>) + Send + Sync>;

/// Indicates that a module contains a provider which implements the interface.
pub trait HasProvider<I: ?Sized>: ModuleInterface {
    /// Create a service using the provider registered with the interface `I`.
    /// Each call will create a new instance of the service.
    ///
    /// # Examples
    /// ```
    /// # use shaku::{module, HasProvider, Provider};
    /// # use std::sync::Arc;
    /// #
    /// # trait Foo {}
    /// #
    /// # #[derive(Provider)]
    /// # #[shaku(interface = Foo)]
    /// # struct FooImpl;
    /// # impl Foo for FooImpl {}
    /// #
    /// # module! {
    /// #     TestModule {
    /// #         components = [],
    /// #         providers = [FooImpl]
    /// #     }
    /// # }
    /// #
    /// # fn main() {
    /// # let module = TestModule::builder().build();
    /// #
    /// let foo: Box<dyn Foo> = module.provide().unwrap();
    /// # }
    /// ```
    fn provide(&self) -> Result<Box<I>, Box<dyn Error>>;
}
