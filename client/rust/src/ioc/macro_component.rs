/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

//! @author coyzeng@gmail.com
//!
//! Implementation of the `#[derive(Component)]` procedural macro

use proc_macro::{quote, TokenStream};
use quote::format_ident;
use syn::{DeriveInput, Ident, Visibility};
use crate::ioc::derive::get_debug_level;
use crate::ioc::macro_common_output::create_dependency;
use crate::ioc::service::{Property, PropertyDefault, ServiceData};

pub fn expand_derive_component(input: &DeriveInput) -> syn::Result<TokenStream> {
    let service = ServiceData::from_derive_input(input)?;

    let debug_level = get_debug_level();
    if debug_level > 1 {
        println!("Service data parsed from Component input: {:#?}", service);
    }

    let resolve_properties: Vec<TokenStream> = service
        .properties
        .iter()
        .map(create_resolve_property)
        .collect();

    let dependencies: Vec<TokenStream> = service
        .properties
        .iter()
        .filter_map(create_dependency)
        .collect();

    let visibility = &service.metadata.visibility;
    let parameters_properties: Vec<TokenStream> = service
        .properties
        .iter()
        .filter_map(|property| create_parameters_property(property, visibility))
        .collect();

    let parameters_defaults: Vec<TokenStream> = service
        .properties
        .iter()
        .filter_map(|property| create_parameters_default(property, &service.metadata.identifier))
        .collect();

    // Component implementation
    let component_name = service.metadata.identifier;
    let parameters_name = format_ident!("{}Parameters", component_name);
    let parameters_doc = format!(" Parameters for {}", component_name);
    let interface = service.metadata.interface;
    let (generic_impls, generic_tys, generic_where) = service.metadata.generics.split_for_impl();
    let generic_impls_no_parens = &service.metadata.generics.params;
    let output = quote! {
        impl<
            M: ::shaku::Module #(+ #dependencies)*,
            #generic_impls_no_parens
        > ::shaku::Component<M> for #component_name #generic_tys #generic_where {
            type Interface = dyn #interface;
            type Parameters = #parameters_name #generic_tys;

            fn build(context: &mut ::shaku::ModuleBuildContext<M>, params: Self::Parameters) -> Box<Self::Interface> {
                Box::new(Self {
                    #(#resolve_properties),*
                })
            }
        }

        #[doc = #parameters_doc]
        #visibility struct #parameters_name #generic_impls #generic_where {
            #(#parameters_properties),*
        }

        impl #generic_impls ::std::default::Default for #parameters_name #generic_tys #generic_where {
            #[allow(unreachable_code)]
            fn default() -> Self {
                Self {
                    #(#parameters_defaults),*
                }
            }
        }
    };

    if debug_level > 0 {
        println!("{}", output);
    }

    Ok(output)
}

fn create_resolve_property(property: &Property) -> TokenStream {
    let property_name = &property.property_name;

    if property.is_service() {
        quote! {
            #property_name: M::build_component(context)
        }
    } else {
        quote! {
            #property_name: params.#property_name
        }
    }
}

fn create_parameters_property(property: &Property, vis: &Visibility) -> Option<TokenStream> {
    if property.is_service() {
        return None;
    }

    let property_name = &property.property_name;
    let property_type = &property.ty;
    let doc_comment = &property.doc_comment;

    Some(quote! {
        #(#doc_comment)*
        #vis #property_name: #property_type
    })
}

fn create_parameters_default(property: &Property, component_ident: &Ident) -> Option<TokenStream> {
    if property.is_service() {
        return None;
    }

    let property_name = &property.property_name;

    match &property.default {
        PropertyDefault::Provided(default_expr) => Some(quote! {
            #property_name: #default_expr
        }),
        PropertyDefault::NotProvided => Some(quote! {
            #property_name: Default::default()
        }),
        PropertyDefault::NoDefault => {
            let unreachable_msg = format!(
                "There is no default value for `{}::{}`",
                component_ident, property_name
            );

            Some(quote! {
                #property_name: unreachable!(#unreachable_msg)
            })
        }
    }
}
