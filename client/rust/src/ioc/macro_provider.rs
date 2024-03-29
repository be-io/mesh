/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

//! @author coyzeng@gmail.com
//!
//! Implementation of the `#[derive(Provider)]` procedural macro

use proc_macro::{quote, TokenStream};

use syn::{DeriveInput, Error};

use crate::ioc::derive::get_debug_level;
use crate::ioc::macro_common_output::create_dependency;
use crate::ioc::service::{Property, PropertyType, ServiceData};

pub fn expand_derive_provider(input: &DeriveInput) -> syn::Result<TokenStream> {
    let service = ServiceData::from_derive_input(input)?;

    let debug_level = get_debug_level();
    if debug_level > 1 {
        println!("Service data parsed from Provider input: {:#?}", service);
    }

    let resolve_properties: Vec<TokenStream> = service
        .properties
        .iter()
        .map(create_property_assignment)
        .collect::<Result<_, _>>()?;

    let dependencies: Vec<TokenStream> = service
        .properties
        .iter()
        .filter_map(create_dependency)
        .collect();

    // Provider implementation
    let provider_name = service.metadata.identifier;
    let interface = service.metadata.interface;
    let (_, generic_tys, generic_where) = service.metadata.generics.split_for_impl();
    let generic_impls_no_parens = &service.metadata.generics.params;
    let output = quote! {
        impl<
            M: ::shaku::Module #(+ #dependencies)*,
            #generic_impls_no_parens
        > ::shaku::Provider<M> for #provider_name #generic_tys #generic_where {
            type Interface = dyn #interface;

            fn provide(module: &M) -> ::std::result::Result<
                Box<Self::Interface>,
                Box<dyn ::std::error::Error>
            > {
                Ok(Box::new(Self {
                    #(#resolve_properties),*
                }))
            }
        }
    };

    if debug_level > 0 {
        println!("{}", output);
    }

    Ok(output)
}

fn create_property_assignment(property: &Property) -> syn::Result<TokenStream> {
    let property_name = &property.property_name;

    match property.property_type {
        PropertyType::Component => Ok(quote! {
            #property_name: module.resolve()
        }),
        PropertyType::Provided => Ok(quote! {
            #property_name: module.provide()?
        }),
        PropertyType::Parameter => Err(Error::new(
            property.property_name.span(),
            "Parameters are not allowed in Providers",
        )),
    }
}
