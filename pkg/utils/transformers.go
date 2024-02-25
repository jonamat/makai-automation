package utils

func BoolToBytes(b bool) []byte {
	if b {
		return []byte("1")
	}
	return []byte("0")
}

func BytesToBool(b []byte) bool {
	return string(b) == "1"
}
