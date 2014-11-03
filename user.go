package u2f

import (
	"crypto/x509"
)

type User struct {
	User      string
	Enrolled  bool
	Enrolling EnrollJSON
	KeyHandle string
	//PubKey        string
	SignChallenge string
	Cert          *x509.Certificate
	PubKey        *x509.Certificate
	Counter       uint32
}

type Users interface {
	GetUser(user string) (User, error)
	PutUser(u User)
}
