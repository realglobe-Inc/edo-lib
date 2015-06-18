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
	"crypto/ecdsa"
	"github.com/realglobe-Inc/edo-lib/jwk"
)

// args は順に kid, kty, use, op, alg として使う。
// 以下に該当する最初の鍵を返す。
// kid を指定している場合、ID が一致する。
// 指定していなければ、ID は任意。
// use または op が指定されている場合、用途がいずれかに合致する。
// 指定していなければ、用途は任意。
// alg を指定している場合、使用アルゴリズムが合致する。
// 指定していなければ、使用アルゴリズムは任意。
func findKey(keys []jwk.Key, args ...string) jwk.Key {
	var kid, kty, use, op, alg string
	if len(args) > 0 {
		kid = args[0]
	}
	if len(args) > 1 {
		kty = args[1]
	}
	if len(args) > 2 {
		use = args[2]
	}
	if len(args) > 3 {
		op = args[3]
	}
	if len(args) > 4 {
		alg = args[4]
	}

	for _, key := range keys {
		if (kid == "" || key.Id() == kid) && match(key, kty, use, op, alg) {
			return key
		}
	}
	return nil
}

func match(key jwk.Key, kty, use, op, alg string) bool {
	return matchTypeAndType(key.Type(), kty) &&
		matchTypeAndUse(key.Type(), use) &&
		matchTypeAndOperation(key.Type(), op) &&
		matchTypeAndAlgorithm(key.Type(), alg) &&
		matchUseAndType(key.Use(), kty) &&
		matchUseAndUse(key.Use(), use) &&
		matchUseAndOperation(key.Use(), op) &&
		matchUseAndAlgorithm(key.Use(), alg) &&
		func() bool {
			if len(key.Operations()) == 0 {
				return true
			}
			for keyOp := range key.Operations() {
				if matchOperationAndType(keyOp, kty) &&
					matchOperationAndUse(keyOp, use) &&
					matchOperationAndOperation(keyOp, op) &&
					matchOperationAndAlgorithm(keyOp, alg) {
					return true
				}
			}
			return false
		}() &&
		matchAlgorithmAndType(key.Algorithm(), kty) &&
		matchAlgorithmAndUse(key.Algorithm(), use) &&
		matchAlgorithmAndOperation(key.Algorithm(), op) &&
		matchAlgorithmAndAlgorithm(key.Algorithm(), alg) &&
		matchDetail(key, kty, use, op, alg)
}

func matchTypeAndType(kty, kty2 string) bool {
	return kty == "" || kty2 == "" || kty == kty2
}

func matchTypeAndUse(kty, use string) bool {
	return true
}

func matchTypeAndOperation(kty, op string) bool {
	return true
}

func matchTypeAndAlgorithm(kty, alg string) bool {
	if kty == "" || alg == "" {
		return true
	}

	switch kty {
	case tagEc:
		switch alg {
		case tagEs256,
			tagEs384,
			tagEs512,
			tagEcdh_es,
			tagEcdh_es_a128Kw,
			tagEcdh_es_a192Kw,
			tagEcdh_es_a256Kw:
			return true
		}
	case tagRsa:
		switch alg {
		case tagRs256,
			tagRs384,
			tagRs512,
			tagPs256,
			tagPs384,
			tagPs512,
			tagRsa1_5,
			tagRsa_oaep,
			tagRsa_oaep_256:
			return true
		}
	case tagOct:
		switch alg {
		case tagHs256,
			tagHs384,
			tagHs512,
			tagA128Kw,
			tagA192Kw,
			tagA256Kw,
			tagDir,
			tagA128Gcmkw,
			tagA192Gcmkw,
			tagA256Gcmkw:
			return true
		}
	}

	return false
}

func matchUseAndType(use, kty string) bool {
	return matchTypeAndUse(kty, use)
}

func matchUseAndUse(use, use2 string) bool {
	return use == "" || use2 == "" || use == use2
}

func matchUseAndOperation(use, op string) bool {
	if use == "" || op == "" {
		return true
	}

	switch use {
	case tagSig:
		switch op {
		case tagSign,
			tagVerify:
			return true
		}
	case tagEnc:
		switch op {
		case tagEncrypt,
			tagDecrypt,
			tagWrapKey,
			tagUnwrapKey,
			tagDeriveKey,
			tagDeriveBits:
			return true
		}
	}

	return false
}

func matchUseAndAlgorithm(use, alg string) bool {
	if use == "" || alg == "" {
		return true
	}

	switch use {
	case tagSig:
		switch alg {
		case tagHs256,
			tagHs384,
			tagHs512,
			tagRs256,
			tagRs384,
			tagRs512,
			tagEs256,
			tagEs384,
			tagEs512,
			tagPs256,
			tagPs384,
			tagPs512,
			tagNone:
			return true
		}
	case tagEnc:
		switch alg {
		case tagRsa1_5,
			tagRsa_oaep,
			tagRsa_oaep_256,
			tagA128Kw,
			tagA192Kw,
			tagA256Kw,
			tagDir,
			tagEcdh_es,
			tagEcdh_es_a128Kw,
			tagEcdh_es_a192Kw,
			tagEcdh_es_a256Kw,
			tagA128Gcmkw,
			tagA192Gcmkw,
			tagA256Gcmkw,
			tagPbes2_hs256_a128Kw,
			tagPbes2_hs384_a192Kw,
			tagPbes2_hs512_a256Kw:
			return true
		}
	}

	return false
}

func matchOperationAndType(op, kty string) bool {
	return matchTypeAndOperation(kty, op)
}

func matchOperationAndUse(op, use string) bool {
	return matchUseAndOperation(use, op)
}

func matchOperationAndOperation(op, op2 string) bool {
	return op == "" || op2 == "" || op == op2
}

func matchOperationAndAlgorithm(op, alg string) bool {
	if op == "" || alg == "" {
		return true
	}

	switch op {
	case tagSign, tagVerify:
		switch alg {
		case tagHs256,
			tagHs384,
			tagHs512,
			tagRs256,
			tagRs384,
			tagRs512,
			tagEs256,
			tagEs384,
			tagEs512,
			tagPs256,
			tagPs384,
			tagPs512,
			tagNone:
			return true
		}
	case tagEncrypt, tagDecrypt:
		switch alg {
		case tagDir:
			return true
		}
	case tagWrapKey, tagUnwrapKey:
		switch alg {
		case tagRsa1_5,
			tagRsa_oaep,
			tagRsa_oaep_256,
			tagA128Kw,
			tagA192Kw,
			tagA256Kw,
			tagEcdh_es,
			tagEcdh_es_a128Kw,
			tagEcdh_es_a192Kw,
			tagEcdh_es_a256Kw,
			tagA128Gcmkw,
			tagA192Gcmkw,
			tagA256Gcmkw,
			tagPbes2_hs256_a128Kw,
			tagPbes2_hs384_a192Kw,
			tagPbes2_hs512_a256Kw:
			return true
		}
	}

	return false
}

func matchAlgorithmAndType(alg, kty string) bool {
	return matchTypeAndAlgorithm(kty, alg)
}

func matchAlgorithmAndUse(alg, use string) bool {
	return matchUseAndAlgorithm(use, alg)
}

func matchAlgorithmAndOperation(alg, op string) bool {
	return matchOperationAndAlgorithm(op, alg)
}

func matchAlgorithmAndAlgorithm(alg, alg2 string) bool {
	return alg == "" || alg2 == "" || alg == alg2
}

func matchDetail(key jwk.Key, kty, use, op, alg string) bool {
	switch key.Type() {
	case tagRsa, tagEc:
		switch op {
		case tagSign, tagUnwrapKey:
			if key.Private() == nil {
				// 署名・復号には秘密鍵が必要。
				return false
			}
		}
	}

	switch alg {
	case tagEs256, tagEs384, tagEs512:
		// JWA の仕様で ESxxx は鍵のサイズが決められている。
		if pub, ok := key.Public().(*ecdsa.PublicKey); !ok || (pub.Params().BitSize+7)/8 != keySizes[alg] {
			return false
		}
	case tagA128Kw,
		tagA192Kw,
		tagA256Kw,
		tagA128Gcmkw,
		tagA192Gcmkw,
		tagA256Gcmkw:
		if len(key.Common()) != keySizes[alg] {
			return false
		}
	}

	return true
}
