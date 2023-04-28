/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

//! @author coyzeng@gmail.com
//!
//! Functions which create common tokenstream outputs

use proc_macro::{quote, TokenStream};
use crate::ioc::service::{Property, PropertyType};

pub fn create_dependency(property: &Property) -> Option<TokenStream> {
    let property_ty = &property.ty;

    match property.property_type {
        PropertyType::Parameter => None,
        PropertyType::Component => Some(quote! {
            ::shaku::HasComponent<#property_ty>
        }),
        PropertyType::Provided => Some(quote! {
            ::shaku::HasProvider<#property_ty>
        }),
    }
}
