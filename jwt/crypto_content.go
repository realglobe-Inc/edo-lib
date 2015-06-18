// Copyright 2015 realglobe, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package jwt

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	hash "github.com/realglobe-Inc/edo-lib/hash"
	"github.com/realglobe-Inc/go-lib/erro"
	"math/big"
)

// AES_CBC_HMAC_SHA2 暗号。
func encryptAesCbcHmacSha2(key []byte, hGen crypto.Hash, plain, authData, initVec []byte) (encrypted, authTag []byte, err error) {
	const authDataLenByteSize = 8

	keyLen := len(key) / 2
	macKey := key[:keyLen]
	cryKey := key[keyLen:]
	authTagLen := keyLen

	bl, err := aes.NewCipher(cryKey)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	}

	// 暗号化。
	encrypter := cipher.NewCBCEncrypter(bl, initVec)
	plain = paddingPkcs7(plain, encrypter.BlockSize())
	encrypted = make([]byte, len(plain))
	encrypter.CryptBlocks(encrypted, plain)

	// 署名？
	hVal := hash.Hashing(hmac.New(hGen.New, macKey), authData, initVec, encrypted, bigEndian(int64(8*len(authData)), authDataLenByteSize))

	return encrypted, hVal[:authTagLen], nil
}

// AES_CBC_HMAC_SHA2 復号。
func decryptAesCbcHmacSha2(key []byte, hGen crypto.Hash, authData, initVec, encrypted, authTag []byte) (plain []byte, err error) {
	if len(encrypted) == 0 {
		// PKCS#7 パディングするのでブロック長未満はおかしい。
		return nil, erro.New("no encrypted data")
	}

	const authDataLenByteSize = 8

	keyLen := len(key) / 2
	macKey := key[:keyLen]
	cryKey := key[keyLen:]
	authTagLen := keyLen

	bl, err := aes.NewCipher(cryKey)
	if err != nil {
		return nil, erro.Wrap(err)
	} else if len(encrypted)%bl.BlockSize() != 0 {
		return nil, erro.New("encrypted data size is not multiple of block size")
	}

	// 検証。
	hVal := hash.Hashing(hmac.New(hGen.New, macKey), authData, initVec, encrypted, bigEndian(int64(8*len(authData)), authDataLenByteSize))
	if !hmac.Equal(authTag, hVal[:authTagLen]) {
		return nil, erro.New("verification error")
	}

	// 復号。
	decrypter := cipher.NewCBCDecrypter(bl, initVec)
	plain = make([]byte, len(encrypted))
	decrypter.CryptBlocks(plain, encrypted)
	return unpaddingPkcs7(plain, decrypter.BlockSize()), nil
}

// パディングサイズからパディングするバイト列へのマップ。
var pkcs7Pads [][]byte

func init() {
	maxBlockSize := aes.BlockSize

	pkcs7Pads = make([][]byte, maxBlockSize+1)
	for padSize := 1; padSize < len(pkcs7Pads); padSize++ {
		pad := make([]byte, padSize)
		for j := 0; j < len(pad); j++ {
			pad[j] = byte(padSize)
		}
		pkcs7Pads[padSize] = pad
	}
}

func paddingPkcs7(data []byte, blSize int) []byte {
	if rem := len(data) % blSize; rem == 0 {
		return append(data, pkcs7Pads[blSize]...)
	} else {
		return append(data, pkcs7Pads[blSize-rem]...)
	}
}

func unpaddingPkcs7(data []byte, blSize int) []byte {
	return data[:len(data)-int(data[len(data)-1])]
}

func bigEndian(val int64, byteSize int) []byte {
	buff := big.NewInt(val).Bytes()
	if len(buff) < byteSize {
		b := make([]byte, byteSize)
		copy(b[byteSize-len(buff):], buff)
		buff = b
	} else if len(buff) > byteSize {
		buff = buff[len(buff)-byteSize:]
	}
	return buff
}

// AES_GCM Encryption.
func encryptAesGcm(key, plain, authData, initVec []byte) (encrypted, authTag []byte, err error) {
	bl, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	}
	encrypter, err := cipher.NewGCM(bl)
	if err != nil {
		return nil, nil, erro.Wrap(err)
	}

	// 暗号化と署名？
	buff := encrypter.Seal(nil, initVec, plain, authData)

	return buff[:len(buff)-encrypter.Overhead()], buff[len(buff)-encrypter.Overhead():], nil
}

// AES_GCM Decryption.
func decryptAesGcm(key, authData, initVec, encrypted, authTag []byte) (plain []byte, err error) {
	bl, err := aes.NewCipher(key)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	decrypter, err := cipher.NewGCM(bl)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	// 復号と検証。
	plain, err = decrypter.Open(nil, initVec, append(encrypted, authTag...), authData)
	if err != nil {
		return nil, erro.Wrap(err)
	}

	return plain, nil
}
