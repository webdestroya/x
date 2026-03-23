package kompressor

import (
	"bytes"
	"compress/zlib"
	"io"
)

func uncompressZlib(buf []byte) ([]byte, error) {
	reader, err := zlib.NewReader(bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func compressZlib(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := zlib.NewWriter(&buf)

	_, err := writer.Write(data)
	if err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
