package u2f

import (
	"encoding/base64"
	"strings"
)

func unb64u(s string) ([]byte, error) {
	// fix padding
	if l := len(s) % 4; l != 0 {
		s += strings.Repeat("=", 4-l)
	}
	return base64.URLEncoding.DecodeString(s)
}

func b64u(s []byte) string {
	b := base64.URLEncoding.EncodeToString(s)
	return strings.Trim(b, "=")
}
