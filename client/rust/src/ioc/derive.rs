/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

//! @author coyzeng@gmail.com
//!
extern crate proc_macro;
#[macro_use]
extern crate quote;

use proc_macro::TokenStream;
use std::env;

use syn::Attribute;

use crate::structures::module::ModuleData;

pub(self) use self::key_value::KeyValue;

pub const ATTR_NAME: &str = "shaku";
pub const INTERFACE_ATTR_NAME: &str = "interface";
pub const INJECT_ATTR_NAME: &str = "inject";
pub const PROVIDE_ATTR_NAME: &str = "provide";
pub const DEFAULT_ATTR_NAME: &str = "default";
pub const DEBUG_ENV_VAR: &str = "SHAKU_CODEGEN_DEBUG";


pub fn get_debug_level() -> usize {
    env::var(consts::DEBUG_ENV_VAR)
        .ok()
        .and_then(|value| value.parse().ok())
        .unwrap_or(0)
}


#[proc_macro_derive(Component, attributes(shaku))]
pub fn component(input: TokenStream) -> TokenStream {
    let input = syn::parse_macro_input!(input as syn::DeriveInput);

    macros::component::expand_derive_component(&input)
        .unwrap_or_else(|e| e.to_compile_error())
        .into()
}

#[proc_macro_derive(Provider, attributes(shaku))]
pub fn provider(input: TokenStream) -> TokenStream {
    let input = syn::parse_macro_input!(input as syn::DeriveInput);

    macros::provider::expand_derive_provider(&input)
        .unwrap_or_else(|e| e.to_compile_error())
        .into()
}

/// Create a [`Module`] which is associated with some components and providers.
///
/// ## Builder
/// A `fn builder(submodules...) -> ModuleBuilder<Self>` associated function will be created to make
/// instantiating the module convenient. The arguments are the submodules the module uses.
///
/// ## Module interfaces
/// After the module name, you can add `: MyModuleInterface` where `MyModuleInterface` is the trait
/// that you want this module to implement (ex. `trait MyModuleInterface: HasComponent<MyComponent> {}`).
/// The macro will implement this trait for the module automatically. That is, it is the same as
/// manually adding the line: `impl MyModuleInterface for MyModule {}`. See `MyModuleImpl` in the
/// example below. See also [`ModuleInterface`].
///
/// ## Submodules
/// A module can use components/providers from other modules by explicitly listing the interfaces
/// from each submodule they want to use. Submodules can be abstracted by depending on traits
/// instead of implementations. See `MySecondModule` in the example below.
///
/// See also the [submodules getting started guide].
///
/// ## Generics
/// This macro supports generics at the module level:
/// ```rust
/// use shaku::{module, Component, Interface, HasComponent};
///
/// trait MyComponent<T: Interface>: Interface {}
///
/// #[derive(Component)]
/// #[shaku(interface = MyComponent<T>)]
/// struct MyComponentImpl<T: Interface + Default> {
///     value: T
/// }
/// impl<T: Interface + Default> MyComponent<T> for MyComponentImpl<T> {}
///
/// // MyModuleImpl implements Module and HasComponent<dyn MyComponent<T>>
/// module! {
///     MyModule<T: Interface> where T: Default {
///         components = [MyComponentImpl<T>],
///         providers = []
///     }
/// }
/// # fn main() {}
/// ```
///
/// ## Circular dependencies
/// This macro will detect circular dependencies at compile time. The error that is thrown will be
/// something like
/// "overflow evaluating the requirement `TestModule: HasComponent<(dyn Component1Trait + 'static)>`".
///
/// It is still possible to compile with a circular dependency if the module is manually implemented
/// in a certain way. In that case, there will be a panic during module creation with more details.
///
/// ## Lazy Components
/// Components can be lazily created by annotating them with `#[lazy]` in the module declaration.
/// The component will not be built until it is required, such as when `resolve_ref` is called for
/// the first time.
///
/// ```rust
/// use shaku::{module, Component, Interface};
///
/// trait Service: Interface {}
///
/// #[derive(Component)]
/// #[shaku(interface = Service)]
/// struct ServiceImpl;
/// impl Service for ServiceImpl {}
///
/// module! {
///     MyModule {
///         components = [#[lazy] ServiceImpl],
///         providers = []
///     }
/// }
/// # fn main() {}
/// ```
///
/// # Examples
/// ```
/// use shaku::{module, Component, Interface, HasComponent};
///
/// trait MyComponent: Interface {}
/// trait MyModule: HasComponent<dyn MyComponent> {}
///
/// #[derive(Component)]
/// #[shaku(interface = MyComponent)]
/// struct MyComponentImpl;
/// impl MyComponent for MyComponentImpl {}
///
/// // MyModuleImpl implements Module, MyModule, and HasComponent<dyn MyComponent>
/// module! {
///     MyModuleImpl: MyModule {
///         components = [MyComponentImpl],
///         providers = []
///     }
/// }
///
/// // MySecondModule implements HasComponent<dyn MyComponent> by using
/// // MyModule's implementation.
/// module! {
///     MySecondModule {
///         components = [],
///         providers = [],
///
///         use MyModule {
///             components = [MyComponent],
///             providers = []
///         }
///     }
/// }
/// # fn main() {}
/// ```
///
/// [`Module`]: trait.Module.html
/// [`ModuleInterface`]: trait.ModuleInterface.html
/// [submodules getting started guide]: guide/submodules/index.html
#[proc_macro]
pub fn module(input: TokenStream) -> TokenStream {
    let module = syn::parse_macro_input!(input as ModuleData);

    macros::module::expand_module_macro(module)
        .unwrap_or_else(|e| e.to_compile_error())
        .into()
}

/// Generic parser for syn structures
// Note: Can't use `std::convert::From` here because we don't want to consume `T`
pub trait Parser<T: Sized> {
    fn parse_as(&self) -> syn::Result<T>;
}

/// Find the #[shaku(...)] attribute
fn get_shaku_attribute(attrs: &[Attribute]) -> Option<&Attribute> {
    attrs.iter().find(|a| a.path.is_ident(consts::ATTR_NAME))
}
