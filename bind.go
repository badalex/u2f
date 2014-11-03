package u2f

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
)

type bindJSON struct {
	ClientData       string `json:"clientData"`
	RegistrationData string `json:"registrationData"`
}

func (u2f *U2F) Bind(u User, r io.Reader) error {
	if u.Enrolled {
		return fmt.Errorf("User '%s' already enrolled", u.User)
	}

	if len(u.Enrolling.Challenge) == 0 {
		return fmt.Errorf("user has not started enroll")
	}

	buf := make([]byte, len("data="))
	n, err := r.Read(buf)
	if err != nil {
		return err
	}
	if n != cap(buf) {
		return fmt.Errorf("failed to read all of data")
	}

	j := json.NewDecoder(r)
	b := bindJSON{}
	err = j.Decode(&b)
	if err != nil {
		return err
	}

	if b.RegistrationData == "" {
		return fmt.Errorf("malformed JSON, missing registrationData")
	}

	err = u2f.validateClientData("navigator.id.finishEnrollment", b.ClientData, u.Enrolling.Challenge)
	if err != nil {
		return err
	}

	err = u2f.validateRegistrationData(b, u)
	if err != nil {
		return err
	}

	return nil
}

var PubKeyLen = 65

type cert struct {
	Raw asn1.RawContent
}

func (u2f *U2F) validateRegistrationData(b bindJSON, u User) error {
	data, err := base64.URLEncoding.DecodeString(b.RegistrationData)
	if err != nil {
		return err
	}

	// format is 0x05, pubKey, len, keyHandle, cert, sig
	if data[0] != 0x05 {
		return fmt.Errorf("invalid format, expected 0x05: %x", data[0])
	}

	pubKey := data[1 : PubKeyLen+1]
	pubKeyCert, err := pemToCert(pubKey)
	if err != nil {
		return fmt.Errorf("pubKeyBad: %s", err)
	}

	data = data[1+PubKeyLen : len(data)]

	khLen := data[0]
	keyHandle := data[1 : khLen+1]
	data = data[1+khLen : len(data)]

	c := cert{}
	rest, err := asn1.Unmarshal(data, &c)
	if err != nil {
		return err
	}
	der := data[0 : len(data)-len(rest)]
	sig := data[len(der):len(data)]

	cert, err := x509.ParseCertificate(der)
	if err != nil {
		return err
	}

	app := sha256.Sum256([]byte(u.Enrolling.AppID))

	// xxx we already have done this up above
	cd, err := base64.URLEncoding.DecodeString(b.ClientData)
	if err != nil {
		return err
	}
	chal := sha256.Sum256([]byte(cd))

	var verify = []byte{0}
	verify = append(verify, app[:]...)
	verify = append(verify, chal[:]...)
	verify = append(verify, keyHandle...)
	verify = append(verify, pubKey...)

	err = cert.CheckSignature(x509.ECDSAWithSHA256, verify, sig)
	if err != nil {
		return err
	}

	u.KeyHandle = base64.URLEncoding.EncodeToString(keyHandle)
	u.Cert = cert
	u.PubKey = pubKeyCert
	u.Enrolling = EnrollJSON{}
	u.Enrolled = true
	u2f.UserList.PutUser(u)

	return nil
}

var derPrefix = []byte("\x30\x59\x30\x13\x06\x07\x2a\x86\x48\xce\x3d\x02\x01\x06\x08\x2a\x86\x48\xce\x3d\x03\x01\x07\x03\x42\x00")

func pemToCert(cert []byte) (*x509.Certificate, error) {
	c := append(derPrefix, cert...)
	key, err := x509.ParsePKIXPublicKey(c)
	if err != nil {
		return nil, err
	}

	return &x509.Certificate{PublicKey: key}, err
}
