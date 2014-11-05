package u2f

import "crypto/rand"

// MakeChallenge interface for creating a New challenge
// Generally the only reason to change this is for testing
type MakeChallenge interface {
	New() (string, error)
}

// RandChallenge default implementation of the MakeChallenge interfaces, uses
// golangs crypto/rand Read() interface
type RandChallenge struct {
}

func (_ RandChallenge) New() (string, error) {
	c := make([]byte, 32)
	_, err := rand.Read(c)
	if err != nil {
		return "", err
	}
	return b64u(c), nil
}
