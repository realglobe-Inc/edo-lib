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

const (
	tagAlg        = "alg"
	tagDecrypt    = "decrypt"
	tagDeriveBits = "deriveBits"
	tagDeriveKey  = "deriveKey"
	tagDir        = "dir"
	tagEnc        = "enc"
	tagEncrypt    = "encrypt"
	tagIv         = "iv"
	tagKid        = "kid"
	tagNone       = "none"
	tagOct        = "oct"
	tagSig        = "sig"
	tagSign       = "sign"
	tagTag        = "tag"
	tagUnwrapKey  = "unwrapKey"
	tagVerify     = "verify"
	tagWrapKey    = "wrapKey"
	tagZip        = "zip"

	// 大文字。
	tagA128Cbc_Hs256 = "A128CBC-HS256"
	tagA128Gcm       = "A128GCM"
	tagA128Gcmkw     = "A128GCMKW"
	tagA128Kw        = "A128KW"
	tagA192Cbc_Hs384 = "A192CBC-HS384"
	tagA192Gcm       = "A192GCM"
	tagA192Gcmkw     = "A192GCMKW"
	tagA192Kw        = "A192KW"
	tagA256Cbc_Hs512 = "A256CBC-HS512"
	tagA256Gcm       = "A256GCM"
	tagA256Gcmkw     = "A256GCMKW"
	tagA256Kw        = "A256KW"
	tagDef           = "DEF"
	tagEc            = "EC"
	tagEcdh_es       = "ECDH-ES"
	tagEs256         = "ES256"
	tagEs384         = "ES384"
	tagEs512         = "ES512"
	tagHs256         = "HS256"
	tagHs384         = "HS384"
	tagHs512         = "HS512"
	tagPs256         = "PS256"
	tagPs384         = "PS384"
	tagPs512         = "PS512"
	tagRs256         = "RS256"
	tagRs384         = "RS384"
	tagRs512         = "RS512"
	tagRsa           = "RSA"
	tagRsa1_5        = "RSA1_5"
	tagRsa_oaep      = "RSA-OAEP"
	tagRsa_oaep_256  = "RSA-OAEP-256"

	// 大文字、+ 交じり。
	tagEcdh_es_a128Kw     = "ECDH-ES+A128KW"
	tagEcdh_es_a192Kw     = "ECDH-ES+A192KW"
	tagEcdh_es_a256Kw     = "ECDH-ES+A256KW"
	tagPbes2_hs256_a128Kw = "PBES2-HS256+A128KW"
	tagPbes2_hs384_a192Kw = "PBES2-HS384+A192KW"
	tagPbes2_hs512_a256Kw = "PBES2-HS512+A256KW"
)

var keySizes = map[string]int{
	tagEs256:         32,
	tagEs384:         48,
	tagEs512:         66, // (521 + 7) / 8
	tagA128Kw:        16,
	tagA192Kw:        24,
	tagA256Kw:        32,
	tagA128Gcmkw:     16,
	tagA192Gcmkw:     24,
	tagA256Gcmkw:     32,
	tagA128Cbc_Hs256: 32,
	tagA192Cbc_Hs384: 48,
	tagA256Cbc_Hs512: 64,
	tagA128Gcm:       16,
	tagA192Gcm:       24,
	tagA256Gcm:       32,
}
