package u2f

import "crypto/rand"

func challenge() (string, error) {
	c := make([]byte, 32)
	_, err := rand.Read(c)
	if err != nil {
		return "", err
	}
	return b64u(c), nil
}
