package jwt

import (
	"bytes"
	"compress/flate"
	"github.com/realglobe-Inc/go-lib-rg/erro"
	"io/ioutil"
)

func defCompress(data []byte) ([]byte, error) {
	var buff bytes.Buffer
	if compressor, err := flate.NewWriter(&buff, flate.DefaultCompression); err != nil {
		return nil, erro.Wrap(err)
	} else if _, err := compressor.Write(data); err != nil {
		return nil, erro.Wrap(err)
	} else if compressor.Close(); err != nil {
		return nil, erro.Wrap(err)
	}
	return buff.Bytes(), nil
}

func defDecompress(data []byte) ([]byte, error) {
	decompressor := flate.NewReader(bytes.NewReader(data))
	defer decompressor.Close()
	buff, err := ioutil.ReadAll(decompressor)
	if err != nil {
		return nil, erro.Wrap(err)
	}
	return buff, nil
}
