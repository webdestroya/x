package kompressor

import (
	"bytes"
	"compress/flate"
	"io"
)

func uncompressFlate(buf []byte) ([]byte, error) {
	reader := flate.NewReader(bytes.NewReader(buf))

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func compressFlate(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer, err := flate.NewWriter(&buf, flate.DefaultCompression)
	if err != nil {
		return nil, err
	}
	_, err = writer.Write(data)
	if err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
