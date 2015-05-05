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

import ()

const (
	tagEnc = "enc"
	tagAlg = "alg"
	tagKid = "kid"
	tagZip = "zip"
	tagIv  = "iv"
	tagTag = "tag"
)

const (
	ktyEc  = "EC"
	ktyRsa = "RSA"
	ktyOct = "oct"
)

const (
	useSig = "sig"
	useEnc = "enc"
)

const (
	opSign       = "sign"
	opVerify     = "verify"
	opEncrypt    = "encrypt"
	opDecrypt    = "decrypt"
	opWrapKey    = "wrapKey"
	opUnwrapKey  = "unwrapKey"
	opDeriveKey  = "deriveKey"
	opDeriveBits = "deriveBits"
)

const (
	algHs256              = "HS256"
	algHs384              = "HS384"
	algHs512              = "HS512"
	algRs256              = "RS256"
	algRs384              = "RS384"
	algRs512              = "RS512"
	algEs256              = "ES256"
	algEs384              = "ES384"
	algEs512              = "ES512"
	algPs256              = "PS256"
	algPs384              = "PS384"
	algPs512              = "PS512"
	algNone               = "none"
	algRsa1_5             = "RSA1_5"
	algRsa_oaep           = "RSA-OAEP"
	algRsa_oaep_256       = "RSA-OAEP-256"
	algA128Kw             = "A128KW"
	algA192Kw             = "A192KW"
	algA256Kw             = "A256KW"
	algDir                = "dir"
	algEcdh_es            = "ECDH-ES"
	algEcdh_es_a128Kw     = "ECDH-ES+A128KW"
	algEcdh_es_a192Kw     = "ECDH-ES+A192KW"
	algEcdh_es_a256Kw     = "ECDH-ES+A256KW"
	algA128Gcmkw          = "A128GCMKW"
	algA192Gcmkw          = "A192GCMKW"
	algA256Gcmkw          = "A256GCMKW"
	algPbes2_hs256_a128Kw = "PBES2-HS256+A128KW"
	algPbes2_hs384_a192Kw = "PBES2-HS384+A192KW"
	algPbes2_hs512_a256Kw = "PBES2-HS512+A256KW"
)

const (
	encA128Cbc_Hs256 = "A128CBC-HS256"
	encA192Cbc_Hs384 = "A192CBC-HS384"
	encA256Cbc_Hs512 = "A256CBC-HS512"
	encA128Gcm       = "A128GCM"
	encA192Gcm       = "A192GCM"
	encA256Gcm       = "A256GCM"
)

const (
	zipDef = "DEF"
)

var keySizes = map[string]int{
	algEs256:         32,
	algEs384:         48,
	algEs512:         66, // (521 + 7) / 8
	algA128Kw:        16,
	algA192Kw:        24,
	algA256Kw:        32,
	algA128Gcmkw:     16,
	algA192Gcmkw:     24,
	algA256Gcmkw:     32,
	encA128Cbc_Hs256: 32,
	encA192Cbc_Hs384: 48,
	encA256Cbc_Hs512: 64,
	encA128Gcm:       16,
	encA192Gcm:       24,
	encA256Gcm:       32,
}
