package u2f

import (
	"bytes"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
)

// input
type verifyJSON struct {
	ClientData    string `json:"clientData"`
	SignatureData string `json:"signatureData"`
}

type VerifyJSON struct {
	Touch   byte   `json:"touch"`
	Counter uint32 `json:"counter"`
}

func (u2f U2F) Verify(u User, r io.Reader) (vj VerifyJSON, err error) {
	if !u.Enrolled {
		return vj, fmt.Errorf("User '%s' not enrolled", u.User)
	}

	buf := make([]byte, len("data="))
	n, err := r.Read(buf)
	if err != nil {
		return vj, err
	}
	if n != cap(buf) {
		return vj, fmt.Errorf("failed to read all of data")
	}

	j := json.NewDecoder(r)
	b := verifyJSON{}
	err = j.Decode(&b)
	if err != nil {
		return vj, err
	}

	d, err := u2f.validateClientData("navigator.id.getAssertion", b.ClientData, u.Devices)
	if err != nil {
		return vj, err
	}

	t, c, err := u2f.validateSignatureData(b, d)
	if err != nil {
		return vj, err
	}

	err = u2f.Users.PutUser(u)
	if err != nil {
		return vj, fmt.Errorf("failed to put user")
	}

	vj.Touch = t
	vj.Counter = c

	return vj, nil
}

func (u2f U2F) validateSignatureData(b verifyJSON, d *Device) (up byte, counter uint32, err error) {
	data, err := base64.URLEncoding.DecodeString(b.SignatureData)
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
	cd, err := base64.URLEncoding.DecodeString(b.ClientData)
	if err != nil {
		return up, counter, err
	}
	cdHash := sha256.Sum256([]byte(cd))
	appHash := sha256.Sum256([]byte(u2f.AppID))

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
	data, err := base64.URLEncoding.DecodeString(pub)
	if err != nil {
		return nil, err
	}

	key, err := x509.ParsePKIXPublicKey(data)
	if err != nil {
		return nil, err
	}

	return &x509.Certificate{PublicKey: key}, nil
}
