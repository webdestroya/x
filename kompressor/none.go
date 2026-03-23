package kompressor

func uncompressNone(buf []byte) ([]byte, error) {
	return buf, nil
}

func compressNone(data []byte) ([]byte, error) {
	return data, nil
}
