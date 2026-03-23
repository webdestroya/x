package kompressor

/*
import (
	"bytes"
	"compress/lzw"
	"io"
)

func uncompressLzw(buf []byte) ([]byte, error) {
	reader := lzw.NewReader(bytes.NewReader(buf), lzw.LSB, 8)

	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func compressLzw(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := lzw.NewWriter(&buf, lzw.LSB, 8)

	_, err := writer.Write(data)
	if err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
*/
