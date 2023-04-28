/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

//! @author coyzeng@gmail.com
//!

use syn::{DeriveInput, Error, Type};
use syn::spanned::Spanned;

use crate::ioc::derive::{ATTR_NAME, INTERFACE_ATTR_NAME, Parser};
use crate::ioc::parser_key_value::KeyValue;
use crate::ioc::service::MetaData;

impl Parser<MetaData> for DeriveInput {
    fn parse_as(&self) -> syn::Result<MetaData> {
        // Find the shaku(interface = ?) attribute
        let shaku_attribute = get_shaku_attribute(&self.attrs).ok_or_else(|| {
            Error::new(
                self.ident.span(),
                format!(
                    "Unable to find interface. Please add a '#[{}({} = <your trait>)]'",
                    ATTR_NAME,
                    INTERFACE_ATTR_NAME
                ),
            )
        })?;

        // Get the interface key/value
        let interface_kv: KeyValue<Type> = shaku_attribute.parse_args().map_err(|_| {
            Error::new(
                shaku_attribute.span(),
                format!(
                    "Invalid attribute format. The attribute must be in name-value form. \
                     Example: #[{}({} = <your trait>)]",
                    ATTR_NAME,
                    INTERFACE_ATTR_NAME
                ),
            )
        })?;

        if interface_kv.key != INTERFACE_ATTR_NAME {
            return Err(Error::new(
                self.ident.span(),
                format!(
                    "Unable to find interface. Please add a '#[{}({} = <your trait>)]'",
                    ATTR_NAME,
                    INTERFACE_ATTR_NAME
                ),
            ));
        }

        Ok(MetaData {
            identifier: self.ident.clone(),
            generics: self.generics.clone(),
            interface: interface_kv.value,
            visibility: self.vis.clone(),
        })
    }
}
