package kompressor

import (
	"bytes"
	"errors"
	"io"
)

func Base64Unbatcher(bufData []byte, onEachFunc func([]byte) error) error {
	dst := make([]byte, b64enc.DecodedLen(len(bufData)))
	if _, err := b64enc.Decode(dst, bufData); err != nil {
		return err
	}
	return Unbatcher(dst, onEachFunc)
}

// given a batch written by Batcher, call the function for each entry
func Unbatcher(bufData []byte, onEachFunc func([]byte) error) error {

	buf := bytes.NewReader(bufData)

	compTypeRaw, err := buf.ReadByte()
	if err != nil {
		return err
	}

	compType := CompressionType(compTypeRaw)
	_, uncompressFunc, err := GetCompressionFuncs(compType)
	if err != nil {
		return err
	}

	sizeBuf := make([]byte, lengthValueByteSize)

	for {

		if buf.Len() < lengthValueByteSize { //nolint:staticcheck
			break
		}

		if _, err := buf.Read(sizeBuf); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}

		dataSize := int64(endian.Uint32(sizeBuf))

		// cant have zero or negative
		if dataSize <= 0 {
			break
		}

		dataBuf := make([]byte, dataSize)

		if _, err := buf.Read(dataBuf); err != nil {
			// uhh?
			return err
		}

		data, err := uncompressFunc(dataBuf)
		if err != nil {
			return err
		}

		if err := onEachFunc(data); err != nil {
			return err
		}
	}

	return nil
}
