/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package crypt

type Mapper interface {

	// Map compute as map.
	Map() map[string]string
}

type Format interface {

	// FormKey format the input binary.
	FormKey(bytes []byte, headers ...string) (string, error)

	// InformKey informat the input binary
	InformKey(data string) ([]byte, error)

	// FormSIG format the signature.
	FormSIG(bytes []byte, headers ...string) (string, error)

	// InformSIG informat the signature.
	InformSIG(data string) ([]byte, error)

	// FormCipher format the cipher.
	FormCipher(bytes []byte, headers ...string) (string, error)

	// InformCipher informat the cipher.
	InformCipher(data string) ([]byte, error)

	// FormPlain format the plain.
	FormPlain(bytes []byte, headers ...string) (string, error)

	// InformPlain informat the plain.
	InformPlain(data string) ([]byte, error)
}
