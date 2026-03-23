package kompressor_test

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	mathrand "math/rand"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/webdestroya/x/kompressor"
)

func TestKompressor(t *testing.T) {

	totalChunks := 500
	maxSizeBytes := 256 * 1024 * 1000

	chunks := make([][]byte, totalChunks)

	for i := 0; i < totalChunks; i++ {
		chunkSize := 1000 + mathrand.Intn(31000)

		data := make([]byte, chunkSize)

		_, err := rand.Read(data)
		if err != nil {
			t.Fatalf("could not read random: %s", err.Error())
		}

		chunks[i] = data
	}

	chunks[0] = bytes.Repeat([]byte(`x`), maxSizeBytes-5)
	chunks[5] = bytes.Repeat([]byte(`x`), maxSizeBytes-4)

	// chunks[0] = make([]byte, maxSizeBytes-4)
	// _, err := rand.Read(chunks[totalChunks-1])
	// require.NoError(t, err)

	in := make(chan []byte, 100)

	checkerIndex := 0

	onEachBatchFunc := func(data []byte) error {

		require.Exactly(t, chunks[checkerIndex], data)
		checkerIndex++

		return nil
	}

	opts := kompressor.Options{
		CompressionType: kompressor.CompressionNone,
		MaxSizeBytes:    maxSizeBytes,
		OnBatchReadyFunc: func(ctx context.Context, b []byte) error {

			err := kompressor.Unbatcher(b, onEachBatchFunc)
			require.NoError(t, err, "unbatcher err")

			return err
		},
	}

	ctx := context.TODO()

	done := make(chan struct{})

	go func() {
		defer close(done)
		kompresserr := kompressor.Batcher(opts, ctx, in)
		require.NoError(t, kompresserr)
	}()

	go func() {
		defer close(in)

		for i := range chunks {
			in <- chunks[i]
		}
	}()

	<-done

}

func TestAssumptions(t *testing.T) {

	// t.Logf("256KiB base64 = %d", base64.RawStdEncoding.DecodedLen(256*1024*1000))

	t.Run("b64 padding", func(t *testing.T) {
		t.Parallel()
		input := `sadasdsdss`
		// outputPadded := `c2FkYXNkc2Rzcw==`
		outputNoPadded := `c2FkYXNkc2Rzcw`

		b64enc := base64.RawStdEncoding

		dst := make([]byte, b64enc.EncodedLen(len(input)))

		base64.RawStdEncoding.Encode(dst, []byte(input))

		require.Equal(t, outputNoPadded, string(dst))
		require.Equal(t, len(input), b64enc.DecodedLen(len(dst)))

	})

	t.Run("lengthcoders", func(t *testing.T) {
		t.Parallel()
		// data := make([]byte, 8)

	})

}
