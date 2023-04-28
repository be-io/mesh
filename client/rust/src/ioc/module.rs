/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

//! @author coyzeng@gmail.com
//!
//! Structures to hold useful module data

use std::collections::HashSet;
use std::hash::Hash;

use syn::{Attribute, Generics, Ident, token, Type, Visibility};
use syn::parse::Parse;
use syn::punctuated::Punctuated;

use crate::ioc::derive::Parser;

pub type ComponentItem = ModuleItem<ComponentAttribute>;

mod kw {
    syn::custom_keyword!(components);
    syn::custom_keyword!(providers);
}

/// The main module data structure, parsed from the macro input
#[derive(Debug)]
pub struct ModuleData {
    pub metadata: ModuleMetadata,
    pub services: ModuleServices,
    pub submodules: Punctuated<Submodule, syn::Token![,]>,
}

/// Metadata about the module
#[derive(Debug)]
pub struct ModuleMetadata {
    pub visibility: Visibility,
    pub identifier: Ident,
    pub generics: Generics,
    pub interface: Option<Type>,
}

/// A submodule dependency
#[derive(Debug)]
pub struct Submodule {
    pub ty: Type,
    pub services: ModuleServices,
}

/// Services associated with a module/submodule
#[derive(Debug)]
pub struct ModuleServices {
    pub components: ModuleItems<kw::components, ComponentAttribute>,
    pub comma_token: syn::Token![,],
    pub providers: ModuleItems<kw::providers, ProviderAttribute>,
    pub trailing_comma: Option<syn::Token![,]>,
}

/// A list of components/providers
#[derive(Debug)]
pub struct ModuleItems<T: Parse, A: Eq + Hash>
    where
        Attribute: Parser<A>,
{
    pub keyword_token: T,
    pub eq_token: token::Eq,
    pub bracket_token: token::Bracket,
    // Can't use syn::Token![,] here because of
    // https://github.com/rust-lang/rust/issues/50676
    pub items: Punctuated<ModuleItem<A>, token::Comma>,
}

/// An annotated component/provider type
#[derive(Debug)]
pub struct ModuleItem<A: Eq + Hash>
    where
        Attribute: Parser<A>,
{
    pub attributes: HashSet<A>,
    pub ty: Type,
}

impl ModuleItem<ComponentAttribute> {
    /// Check if a component is marked with `#[lazy]`
    pub fn is_lazy(&self) -> bool {
        self.attributes.contains(&ComponentAttribute::Lazy)
    }
}

/// Valid component attributes
#[derive(Debug, Eq, PartialEq, Hash)]
pub enum ComponentAttribute {
    Lazy,
}

/// Valid provider attributes
#[derive(Debug, Eq, PartialEq, Hash)]
pub enum ProviderAttribute {
    // None currently
}
