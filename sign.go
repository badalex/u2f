package u2f

import (
	"fmt"
)

type SignJSON struct {
	KeyHandle string `json:"keyHandle"`
	Challenge string `json:"challenge"`
	AppID     string `json:"appId"`
	Version   string `json:"version"`
}

func (u2f *U2F) Sign(u User) (SignJSON, error) {
	if !u.Enrolled {
		return SignJSON{}, fmt.Errorf("User '%s' is not enrolled", u.User)
	}

	if u.KeyHandle == "" {
		return SignJSON{}, fmt.Errorf("User '%s' has no keyhandle", u.User)
	}

	c, err := u2f.NewChallenge()
	if err != nil {
		return SignJSON{}, err
	}

	e := SignJSON{
		KeyHandle: u.KeyHandle,
		Challenge: c,
		AppID:     u2f.AppID,
		Version:   u2f.Version,
	}
	u.SignChallenge = c

	u2f.UserList.PutUser(u)

	return e, nil

}
