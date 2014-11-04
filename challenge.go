package u2f

import (
	"crypto/rand"
	"encoding/base64"
)

func (u2f U2F) Challenge() (string, error) {
	c := make([]byte, 32)
	_, err := rand.Read(c)
	if err != nil {
		return "", err
	}
	cs := base64.URLEncoding.EncodeToString(c)
	return cs, nil
}
