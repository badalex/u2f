package u2f

import (
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

	c, err := u2f.Challenge()
	if err != nil {
		return EnrollJSON{}, err
	}

	u.Challenge = c
	u2f.UserList.PutUser(u)

	e := EnrollJSON{
		Challenge: c,
		AppID:     u2f.AppID,
		Version:   u2f.Version,
	}

	return e, nil
}
