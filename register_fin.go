package u2f

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"encoding/json"
	"fmt"
	"io"
)

// RegisterResponse dictionary from the fido u2f javascript api.
// Serves as input to RegisterFin()
type RegisterResponse struct {
	ClientData       string `json:"clientData"`
	RegistrationData string `json:"registrationData"`
}

// RegisterFin finish up, validate and store a Registration.
// r should contain an RegisterResponse JSON Object.
func (s Server) RegisterFin(u User, r io.Reader) error {
	j := json.NewDecoder(r)
	b := RegisterResponse{}
	err := j.Decode(&b)
	if err != nil {
		return err
	}

	if b.RegistrationData == "" {
		return fmt.Errorf("malformed JSON, missing registrationData")
	}

	d, err := s.validateClientData("navigator.id.finishEnrollment", b.ClientData, u.Devices)
	if err != nil {
		return err
	}

	err = s.validateRegisterResponse(b, d)
	if err != nil {
		return err
	}

	u.Enrolled = true
	err = s.Users.PutUser(u)
	if err != nil {
		return err
	}

	return nil
}

var pubKeyLen = 65

type cert struct {
	Raw asn1.RawContent
}

func (s Server) validateRegisterResponse(b RegisterResponse, d *Device) error {
	data, err := unb64u(b.RegistrationData)
	if err != nil {
		return err
	}

	// format is 0x05, pubKey, len, keyHandle, cert, sig
	if data[0] != 0x05 {
		return fmt.Errorf("invalid format, expected 0x05: %x", data[0])
	}

	pubKey := data[1 : pubKeyLen+1]
	pubDer, err := pubKeyDer(pubKey)
	if err != nil {
		return err
	}

	data = data[1+pubKeyLen : len(data)]

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

	app := sha256.Sum256([]byte(s.AppID))

	// xxx we already have done this up above
	cd, err := unb64u(b.ClientData)
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

	d.KeyHandle = b64u(keyHandle)
	d.Cert = b64u(cert.Raw)
	d.PubKey = b64u(pubDer)
	d.Challenge = ""

	return nil
}

var derPrefix = []byte("\x30\x59\x30\x13\x06\x07\x2a\x86\x48\xce\x3d\x02\x01\x06\x08\x2a\x86\x48\xce\x3d\x03\x01\x07\x03\x42\x00")

func pubKeyDer(cert []byte) ([]byte, error) {
	c := append(derPrefix, cert...)
	// validate it
	_, err := x509.ParsePKIXPublicKey(c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
