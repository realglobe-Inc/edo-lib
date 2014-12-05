package util

import (
	"bufio"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"io"
)

func Read(scanner *bufio.Scanner) (token string, err error) {
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", erro.Wrap(err)
		} else {
			return "", io.EOF
		}
	}

	return scanner.Text(), nil
}
