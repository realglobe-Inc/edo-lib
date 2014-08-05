package util

import ()

// ハッシュ値計算。
// Java と eclipse を参考にしている。

func DigestString(s string) int {
	prime := 31
	dig := 0
	for _, r := range s {
		dig = prime*dig + int(r)
	}
	return dig
}

func DigestBool(b bool) int {
	if b {
		return 1231
	} else {
		return 1237
	}
}
