package u2f

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

type EnrollJSON struct {
	Version   string `json:"version"`
	AppID     string `json:"appId"`
	Challenge string `json:"challenge"`
}

func (u2f *U2F) Enroll(u User) (EnrollJSON, error) {
	if u.Enrolled {
		return EnrollJSON{}, fmt.Errorf("User '%s' already enrolled", u.User)
	}

	c, err := u2f.NewChallenge()
	if err != nil {
		return EnrollJSON{}, err
	}

	e := EnrollJSON{
		Challenge: c,
		AppID:     u2f.AppID,
		Version:   u2f.Version,
	}
	u.Enrolling = e

	u2f.UserList.PutUser(u)

	return e, nil
}

func (u2f *U2F) NewChallenge() (string, error) {
	c := make([]byte, 32)
	_, err := rand.Read(c)
	if err != nil {
		return "", err
	}
	cs := base64.URLEncoding.EncodeToString(c)
	return cs, nil
}
