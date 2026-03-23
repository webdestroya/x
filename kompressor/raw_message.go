package kompressor

import "encoding/json"

type RawMessage []byte

// below this size, dont bother compressing
const compressionSizeThreshold = 3000

func (m RawMessage) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}

	if defaultCompressionType == CompressionNone || len(m) < compressionSizeThreshold {
		return m, nil
	}

	compress, _, _ := GetCompressionFuncs(defaultCompressionType)

	compData, err := compress(m)
	if err != nil {
		// compression failed, pass through raw
		return m, nil
	}

	if b64enc.EncodedLen(len(compData))+4 > len(m) {
		// compressed version is bigger, not worth it
		return m, nil
	}

	encoded := `"` + string(defaultCompressionType) + `:` + b64enc.EncodeToString(compData) + `"`
	return []byte(encoded), nil
}

func (m *RawMessage) UnmarshalJSON(data []byte) error {
	if m == nil {
		return &json.InvalidUnmarshalError{}
	}

	// null
	if len(data) == 4 && string(data) == "null" {
		return nil
	}

	// not a string, just store the raw bytes (objects, arrays, booleans, numbers)
	if len(data) == 0 || data[0] != '"' {
		*m = append((*m)[0:0], data...)
		return nil
	}

	// it's a JSON string - unquote it
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		*m = append((*m)[0:0], data...)
		return nil
	}

	sBytes := []byte(s)

	if len(sBytes) > 3 && sBytes[1] == ':' {
		// check if it's a valid compression method
		if _, uncompress, err := GetCompressionFuncs(CompressionType(sBytes[0])); err == nil {
			rawData := make([]byte, b64enc.DecodedLen(len(sBytes)-2))

			if n, err := b64enc.Decode(rawData, sBytes[2:]); err == nil {
				if realData, err := uncompress(rawData[:n]); err == nil {
					*m = realData
					return nil
				}
			}
		}
	}

	// not compressed, return as the original JSON string token
	*m = append((*m)[0:0], data...)
	return nil
}
