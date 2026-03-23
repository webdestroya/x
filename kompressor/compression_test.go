package kompressor_test

import (
	"crypto/rand"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/webdestroya/x/kompressor"
)

func TestCompressionFuncs(t *testing.T) {
	tables := []struct {
		label string
		ctype kompressor.CompressionType
	}{
		{"none", kompressor.CompressionNone},
		{"zlib", kompressor.CompressionZLIB},
	}

	for _, table := range tables {
		t.Run(table.label, func(t *testing.T) {

			data := make([]byte, 32*1024)

			_, err := rand.Read(data)
			require.NoError(t, err)

			compFunc, uncompFunc, err := kompressor.GetCompressionFuncs(table.ctype)
			require.NoError(t, err)

			compData, err := compFunc(data)
			require.NoError(t, err)
			require.NotNil(t, compData)

			if len(compData) > len(data) {
				t.Logf("compression is actually bigger: compressed=%d uncompressed=%d", len(compData), len(data))
			}

			uncompData, err := uncompFunc(compData)
			require.NoError(t, err)
			require.NotNil(t, uncompData)

			require.Exactly(t, data, uncompData)

		})
	}

	t.Run("bad type", func(t *testing.T) {
		compFunc, uncompFunc, err := kompressor.GetCompressionFuncs(kompressor.CompressionType(50))
		require.ErrorIs(t, err, kompressor.ErrUnsupportedCompressionError)
		require.Nil(t, compFunc)
		require.Nil(t, uncompFunc)
	})
}

func BenchmarkCompressions(b *testing.B) {
	compTypesToTest := []kompressor.CompressionType{
		// kompressor.CompressionNone,
		kompressor.CompressionZLIB,
		kompressor.CompressionFLATE,
		kompressor.CompressionGZIP,
		kompressor.CompressionNone,
	}

	// HAPPY CASE:
	dummyBytesMSG1 := mustMinifyJSON(b, filepath.Join("testdata", "data.json"))
	bigboy := append(dummyBytesMSG1, dummyBytesMSG1...) //nolint:gocritic // yes i know its not the same slice
	bigboy = append(bigboy, bigboy...)
	bigboy = append(bigboy, bigboy...)

	for _, compType := range compTypesToTest {
		compress, uncompress, _ := kompressor.GetCompressionFuncs(compType)

		b.Run("comp_"+string(compType), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				compData, err := compress(bigboy)
				if err != nil {
					b.Fatalf("got err: %v", err.Error())
				}
				realData, err := uncompress(compData)
				if err != nil {
					b.Fatalf("got err: %v", err.Error())
				}
				if len(realData) != len(bigboy) {
					b.Fatalf("lengths do not match")
				}
			}
		})

	}
}
