package kompressor_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/webdestroya/x/kompressor"
)

func b64Compress(t *testing.T, ctype kompressor.CompressionType, data []byte) string {
	t.Helper()
	compress, _, _ := kompressor.GetCompressionFuncs(ctype)

	compData, err := compress(data)
	if err != nil {
		t.Fatal("compression error???")
		return ""
	}

	return `"` + string(ctype) + `:` + base64.RawStdEncoding.EncodeToString(compData) + `"`
}

func TestRawMessage(t *testing.T) {
	t.Parallel()
	type stuff struct {
		Name string                `json:"name"`
		Data kompressor.RawMessage `json:"data"`
	}

	jsonBytesMSG1 := mustMinifyJSON(t, filepath.Join("testdata", "data.json"))

	rawData := []byte(`{"thing": "yar", "stuff": "foobarbaz", "poop": "smears"}`)
	rawFullMessage := []byte(fmt.Sprintf(`{"name": "dummy", "data": %s}`, string(rawData)))

	// READING REGULAR STUFF
	var thing stuff
	err := json.Unmarshal(rawFullMessage, &thing)
	require.NoError(t, err)

	require.Equal(t, string(rawData), string(thing.Data))

	t.Run("unmarshal regulars", func(t *testing.T) {
		t.Parallel()
		sameAsInput := &struct{}{}

		tables := []struct {
			input       string
			expectation any
		}{
			{`{"thing": "yar", "stuff": "foobarbaz", "poop": "smears"}`, sameAsInput},
			{`{"thing": "yar", "stuff": {"foobarbaz": "inner", "thing": "yar"}, "poop": "smears"}`, sameAsInput},
			{`"Z:this is a goofy test"`, sameAsInput},
			{`"This is a test"`, sameAsInput},
			{`""`, sameAsInput},
			{`true`, sameAsInput},
			{`false`, sameAsInput},
			{`1234`, sameAsInput},
			{`{}`, sameAsInput},
			{`[]`, sameAsInput},
			{`["test"]`, sameAsInput},
			{`null`, nil},

			{
				input:       `"Z:eJyqVirJyMxLV7JSUKpMLFLSUVAqLilNSwPx0/LzkxKLkhKrQKIF+fkFIMHi3NTEomKlWkAAAAD///P9EcQ"`,
				expectation: `{"thing": "yar", "stuff": "foobarbaz", "poop": "smears"}`,
			},

			{
				input:       b64Compress(t, kompressor.CompressionZLIB, jsonBytesMSG1),
				expectation: jsonBytesMSG1,
			},
		}
		for testNum, table := range tables {
			t.Run(fmt.Sprintf("test_%02d", testNum), func(t *testing.T) {
				t.Parallel()
				msgStr := []byte(fmt.Sprintf(`{"name": "dummy", "data": %s}`, table.input))
				var thing stuff
				err := json.Unmarshal(msgStr, &thing)
				require.NoError(t, err)
				require.Equal(t, "dummy", thing.Name)

				if table.expectation == sameAsInput {
					require.Equal(t, table.input, string(thing.Data))
					return
				}

				switch v := table.expectation.(type) {
				case string:
					require.Equal(t, v, string(thing.Data))
				case []byte:
					require.Equal(t, string(v), string(thing.Data))
				case nil:
					require.Nil(t, thing.Data)
				default:
					require.Failf(t, "THIS IS NOT FINISHED", "for this type: %T", v)
				}

			})
		}
	})

	t.Run("marshal", func(t *testing.T) {
		t.Parallel()
		dummy := stuff{
			Name: "thinger",
			Data: rawData,
		}

		jsonBytes, err := json.Marshal(dummy)
		require.NoError(t, err)
		require.NotNil(t, jsonBytes)
		// t.Logf("JSON: %s", string(jsonBytes))

		var dummy2 stuff
		err = json.Unmarshal(jsonBytes, &dummy2)
		require.NoError(t, err)
		require.Equal(t, string(rawData), string(dummy.Data))

	})
}

func TestCompressionUsefulness(t *testing.T) {

	t.SkipNow()

	compTypesToTest := []kompressor.CompressionType{
		// kompressor.CompressionNone,
		kompressor.CompressionZLIB,
		kompressor.CompressionFLATE,
		kompressor.CompressionGZIP,
	}

	// HAPPY CASE:
	dummyBytesMSG1 := mustMinifyJSON(t, filepath.Join("testdata", "data.json"))
	bigboy := append(dummyBytesMSG1, dummyBytesMSG1...) //nolint:gocritic // yes i know its not the same slice
	bigboy = append(bigboy, bigboy...)
	bigboy = append(bigboy, bigboy...)

	sizelist := make([]int, 0, 50)
	sizelist = append(sizelist, 100, 200, 300, 500)
	for i := 0; i < 32; i++ {
		sizelist = append(sizelist, (i+1)*1000)
	}

	for _, compType := range compTypesToTest {
		compress, _, _ := kompressor.GetCompressionFuncs(compType)
		for _, endpos := range sizelist {

			compData, _ := compress(bigboy[:endpos])
			compSize := len(compData)
			b64Size := base64.RawStdEncoding.EncodedLen(compSize) + 2 + 2 // quote + compType + : <data> + quote

			b64Pct := (1.0 - (float64(b64Size) / float64(endpos))) * 100.0

			t.Logf("%s %5d => raw:%5d final: %5d (%.2f%% compression)", string(compType), endpos, compSize, b64Size, b64Pct)

		}
	}

}

type testHelper interface {
	Helper()
}

func mustMinifyJSON(t testHelper, filepath string) []byte {
	t.Helper()
	var tmp any

	data, err := os.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &tmp)
	if err != nil {
		panic(err)
	}

	out, err := json.Marshal(tmp)
	if err != nil {
		panic(err)
	}

	return out

}
