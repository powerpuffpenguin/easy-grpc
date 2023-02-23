package core

import "unsafe"

// 將 []byte 轉爲 字符串
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// 將字符串轉爲 只讀的 []byte
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}
