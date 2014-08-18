package util

import (
	"bufio"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"io"
)

func Read(scanner *bufio.Scanner) (token string, err error) {
	if ok := scanner.Scan(); !ok {
		if err := scanner.Err(); err != nil {
			return "", erro.Wrap(err)
		} else {
			return "", erro.Wrap(io.EOF)
		}
	}

	return scanner.Text(), nil
}
