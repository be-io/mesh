/*
 * Copyright (c) 2000, 2023, trustbe and/or its affiliates. All rights reserved.
 * TRUSTBE PROPRIETARY/CONFIDENTIAL. Use is subject to license terms.
 *
 *
 */

package prsim

import "context"

const (
	PublicKey  CryptorFeat = "public_key"
	PrivateKey CryptorFeat = "private_key"
	Signature  CryptorFeat = "signature"
)

var ICryptor = (*Cryptor)(nil)

type CryptorFeat string

// Cryptor spi
// @MPI(macro.MeshMPI)
// @SPI(macro.MeshSPI)
type Cryptor interface {

	// Encrypt binary to encrypted binary.
	// @MPI("mesh.crypt.encrypt")
	Encrypt(ctx context.Context, buff []byte, features map[string][]byte) ([]byte, error)

	// Decrypt binary to decrypted binary.
	// @MPI("mesh.crypt.decrypt")
	Decrypt(ctx context.Context, buff []byte, features map[string][]byte) ([]byte, error)

	// Hash compute the hash value.
	// @MPI("mesh.crypt.hash")
	Hash(ctx context.Context, buff []byte, features map[string][]byte) ([]byte, error)

	// Sign compute the signature value.
	// @MPI("mesh.crypt.sign")
	Sign(ctx context.Context, buff []byte, features map[string][]byte) ([]byte, error)

	// Verify the signature value.
	// @MPI("mesh.crypt.verify")
	Verify(ctx context.Context, buff []byte, features map[string][]byte) (bool, error)
}
