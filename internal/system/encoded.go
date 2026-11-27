package system

import (
	"encoding/base64"
	"unicode/utf16"
)

// utf16LE 将字符串编码为 UTF-16LE 字节序列（无 BOM）。
func utf16LE(s string) []byte {
	u16 := utf16.Encode([]rune(s))
	b := make([]byte, 0, len(u16)*2)
	for _, c := range u16 {
		b = append(b, byte(c), byte(c>>8))
	}
	return b
}

// encodedCommand 生成 powershell -EncodedCommand 参数：Base64(UTF-16LE 脚本)。
func encodedCommand(script string) string {
	return base64.StdEncoding.EncodeToString(utf16LE(script))
}
