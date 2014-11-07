package u2f

import (
	"bytes"
	"crypto/sha256"
	"crypto/x509"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
)

// SignResponse dictionary from the fido u2f javascript api.
// Serves as input to SignFin()
type SignResponse struct {
	KeyHandle     string `json:"keyHandle"`
	ClientData    string `json:"clientData"`
	SignatureData string `json:"signatureData"`
}

// SignFinResult is the result of a successful SignFin operation.
type SignFinResult struct {
	Touch byte `json:"touch"`
	// Counter current counter value
	Counter uint32 `json:"counter"`
}

// SignFin Finalize a Sign/Login operation. If this succeeds everything is
// good and the usb token has been validated.
// r should contain an SignResponse JSON Object.
func (s U2FServer) SignFin(u User, r io.Reader) (sf SignFinResult, err error) {
	if !u.Enrolled {
		return sf, fmt.Errorf("User '%s' not enrolled", u.User)
	}

	j := json.NewDecoder(r)
	b := SignResponse{}
	err = j.Decode(&b)
	if err != nil {
		return sf, err
	}

	d, err := s.validateClientData("navigator.id.getAssertion", b.ClientData, u.Devices)
	if err != nil {
		return sf, err
	}

	t, c, err := s.validateSignResponse(b, d)
	if err != nil {
		return sf, err
	}

	err = s.Users.PutUser(u)
	if err != nil {
		return sf, fmt.Errorf("failed to put user")
	}

	sf.Touch = t
	sf.Counter = c

	return sf, nil
}

func (s U2FServer) validateSignResponse(b SignResponse, d *Device) (up byte, counter uint32, err error) {
	data, err := unb64u(b.SignatureData)
	if err != nil {
		return up, counter, err
	}

	// userPresence / touch
	up = data[0]

	err = binary.Read(bytes.NewReader(data[1:6]), binary.BigEndian, &counter)
	if err != nil {
		return up, counter, err
	}

	sig := data[5:len(data)]

	// xxx we already have done this up above
	cd, err := unb64u(b.ClientData)
	if err != nil {
		return up, counter, err
	}
	cdHash := sha256.Sum256([]byte(cd))
	appHash := sha256.Sum256([]byte(s.AppID))

	var verify []byte
	verify = append(verify, appHash[:]...)
	verify = append(verify, data[0:5]...)
	verify = append(verify, cdHash[:]...)

	cert, err := pubKeyCert(d.PubKey)
	if err != nil {
		return up, counter, err
	}

	err = cert.CheckSignature(x509.ECDSAWithSHA256, verify, sig)
	if err != nil {
		return up, counter, err
	}

	if d.Counter >= counter {
		return up, counter, fmt.Errorf("Counter mismatch")
	}
	d.Counter = counter

	return up, counter, nil
}

func pubKeyCert(pub string) (*x509.Certificate, error) {
	data, err := unb64u(pub)
	if err != nil {
		return nil, err
	}

	key, err := x509.ParsePKIXPublicKey(data)
	if err != nil {
		return nil, err
	}

	return &x509.Certificate{PublicKey: key}, nil
}
