package util

import (
	"encoding/base64"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"unicode"
)

// 引用符を考慮した引数分割。
// " か ' で括られていたら中身だけにする。ただし、一番外側だけ。
// " ' \ と半角スペースは \ でエスケープ可能。
func Fields(s string) []string {
	fields := []string{}

	field := ""     // 現在のフィールドに加えられている文字。
	quo := uint8(0) // 現在の引用符。
	quoted := false // 引用されたかどうか。空文字列 "" を認識するため。
	esc := false    // 前の文字が '\' だったかどうか。
	for i := 0; i < len(s); i++ {
		if esc {
			esc = false
			switch s[i] {
			case '"', '\'', '\\', ' ':
				field += s[i : i+1]
				continue
			}
			field += "\\"
		}

		// esc == false

		if unicode.IsSpace(rune(s[i])) {
			if quo == 0 {
				if field != "" || quoted {
					fields = append(fields, field)
				}
				field = ""
				quo = 0
				quoted = false
				continue
			}
		}

		switch s[i] {
		case '"', '\'':
			if s[i] == quo {
				// 引用オワタ。
				quo = 0
				continue
			} else if quo != 0 {
				// 引用中。
				field += s[i : i+1]
				continue
			} else {
				// 引用ハジマタ。
				quo = s[i]
				quoted = true
				continue
			}
		case '\\':
			esc = true
			continue
		default:
			field += s[i : i+1]
		}
	}

	if field != "" || quoted {
		fields = append(fields, field)
	}

	return fields
}

func SecureRandomString(length int) (string, error) {
	buff, err := SecureRandomBytes((length*6 + 7) / 8)
	if err != nil {
		return "", erro.Wrap(err)
	}
	return base64.URLEncoding.EncodeToString(buff)[:length], nil
}
