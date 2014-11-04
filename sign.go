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

func (u2f U2F) Sign(u User) (r []SignJSON, err error) {
	if !u.Enrolled {
		return r, fmt.Errorf("User '%s' is not enrolled", u.User)
	}

	c, err := u2f.Challenge()
	if err != nil {
		return r, err
	}

	for i, d := range u.Devices {
		d.Challenge = c
		u.Devices[i] = d

		r = append(r, SignJSON{
			KeyHandle: d.KeyHandle,
			Challenge: c,
			AppID:     u2f.AppID,
			Version:   u2f.Version,
		})
	}

	err = u2f.Users.PutUser(u)
	if err != nil {
		return r, err
	}

	return r, nil

}
