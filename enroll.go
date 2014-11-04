package u2f

import (
	"fmt"
)

type EnrollJSON struct {
	Version   string `json:"version"`
	AppID     string `json:"appId"`
	Challenge string `json:"challenge"`
}

func (u2f U2F) Enroll(u User) (r EnrollJSON, err error) {
	if u.Enrolled {
		return r, fmt.Errorf("User '%s' already enrolled", u.User)
	}

	c, err := u2f.Challenge()
	if err != nil {
		return r, err
	}

	u.Devices = append(u.Devices, Device{
		Challenge: c,
	})
	err = u2f.Users.PutUser(u)
	if err != nil {
		return r, err
	}

	r = EnrollJSON{
		Challenge: c,
		AppID:     u2f.AppID,
		Version:   u2f.Version,
	}
	return r, nil
}
