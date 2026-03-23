// This will let you send stuff to be compressed up to a target size
// once the target is reached, the channel is called with the payload
package kompressor

import (
	"bytes"
	"context"
)

const preambleLength = 1 // <compressionType>

type Options struct {
	// CompressionType
	CompressionType CompressionType // unused

	// max payload limit in bytes
	MaxSizeBytes int

	// called for each batch
	OnBatchReadyFunc func(context.Context, []byte) error

	OnEntryTooBigFunc func(context.Context, []byte) error
}

// Deprecated: Probably should not use this. It doesn't save that much
func Base64Batcher(opts Options, ctx context.Context, in <-chan []byte) error {

	newOpts := Options{
		MaxSizeBytes: b64enc.DecodedLen(opts.MaxSizeBytes),
		OnBatchReadyFunc: func(c context.Context, b []byte) error {
			dst := make([]byte, b64enc.EncodedLen(len(b)))
			b64enc.Encode(dst, b)
			return opts.OnBatchReadyFunc(c, dst)
		},
	}

	return Batcher(newOpts, ctx, in)
}

// Batch up requests and compress them
//
// Deprecated: Probably should not use this. It doesn't save that much
func Batcher(opts Options, ctx context.Context, in <-chan []byte) error {

	if opts.CompressionType == compressionDefault {
		opts.CompressionType = defaultCompressionType
	}

	if opts.OnEntryTooBigFunc == nil {
		opts.OnEntryTooBigFunc = func(_ context.Context, _ []byte) error {
			return ErrPayloadTooBigError
		}
	}

	maxSize := opts.MaxSizeBytes
	if maxSize == 0 {
		return ErrMaxBytesNotSetError
	}

	compressFunc, _, err := GetCompressionFuncs(opts.CompressionType)
	if err != nil {
		return err
	}

	compressData := func(data []byte) ([]byte, error) {
		compData, err := compressFunc(data)
		if err != nil {
			return nil, err
		}

		sizeBuf := make([]byte, lengthValueByteSize)
		endian.PutUint32(sizeBuf, uint32(len(compData))) //nolint:gosec // this will not overflow

		//nolint:makezero // im relying on it resizing the underlying slice
		return append(sizeBuf, compData...), nil
	}

	buf := bytes.NewBuffer(make([]byte, 0, maxSize))
	if err := buf.WriteByte(byte(opts.CompressionType)); err != nil {
		return err
	}

	batchSend := func() error {
		if buf.Len() <= preambleLength {
			// it only has the compression type, abandon
			return nil
		}
		newBuf := make([]byte, buf.Len())
		copy(newBuf, buf.Bytes())
		if err := opts.OnBatchReadyFunc(ctx, newBuf); err != nil {
			return err
		}
		buf.Reset()
		buf.WriteByte(byte(opts.CompressionType))

		return nil
	}

	for {
		select {
		case <-ctx.Done():
			return batchSend()

		case v, ok := <-in:
			if !ok {
				return batchSend()
			}

			msgBytes, err := compressData(v)
			if err != nil {
				return err
			}

			if len(msgBytes) > maxSize {
				if err := opts.OnEntryTooBigFunc(ctx, v); err != nil {
					return err
				}
				continue
			}

			if len(msgBytes) > buf.Available() {
				if err := batchSend(); err != nil {
					return err
				}
			}

			if _, err := buf.Write(msgBytes); err != nil {
				return err
			}

			// if not enough to write another entry
			if (lengthValueByteSize + 1) >= buf.Available() {
				if err := batchSend(); err != nil {
					return err
				}
			}
		}
	}
}
