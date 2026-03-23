package kompressor

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
)

var (
	ErrUnsupportedCompressionError = errors.New("unsupported compression type")
	ErrMaxBytesNotSetError         = errors.New("max bytes was not set")
	ErrPayloadTooBigError          = errors.New("payload will never fit in buffer")
)

var endian = binary.BigEndian
var b64enc = base64.RawStdEncoding

const defaultCompressionType = CompressionFLATE

const lengthValueByteSize = 4 // uint32

type CompressionType byte

const (
	compressionDefault CompressionType = 0
	CompressionNone    CompressionType = 'N'
	CompressionFLATE   CompressionType = 'D' // the best
	CompressionZLIB    CompressionType = 'Z' // 2nd best
	CompressionGZIP    CompressionType = 'G' // has a huge header
	// CompressionLZW     CompressionType = 'L' // it's crap. dont use
)

func GetCompressionFuncs(ctype CompressionType) (func([]byte) ([]byte, error), func([]byte) ([]byte, error), error) {
	switch ctype {
	case CompressionZLIB:
		return compressZlib, uncompressZlib, nil

	case CompressionFLATE:
		return compressFlate, uncompressFlate, nil

	case CompressionGZIP:
		return compressGzip, uncompressGzip, nil

	// case CompressionLZW:
	// 	return compressLzw, uncompressLzw, nil

	case CompressionNone:
		return compressNone, uncompressNone, nil

	default:
		return nil, nil, ErrUnsupportedCompressionError
	}
}
