package u2f

import (
	"fmt"
)

// SignRequest dictionary from the fido u2f javascript api spec.
// Result of a valid Sign() operation
type SignRequest struct {
	Version   string `json:"version"`
	Challenge string `json:"challenge"`
	KeyHandle string `json:"keyHandle"`
	AppID     string `json:"appId"`
}

// Sign Returns SignRequests for the device to Sign. The result should then
// be passed to SignFin() for validation.
func (s U2FServer) Sign(u User) (r []SignRequest, err error) {
	if !u.Enrolled {
		return r, fmt.Errorf("User '%s' is not enrolled", u.User)
	}

	c, err := challenge()
	if err != nil {
		return r, err
	}

	for i, d := range u.Devices {
		d.Challenge = c
		u.Devices[i] = d

		r = append(r, SignRequest{
			Version:   s.Version,
			Challenge: c,
			KeyHandle: d.KeyHandle,
			AppID:     s.AppID,
		})
	}

	err = s.Users.PutUser(u)
	if err != nil {
		return r, err
	}

	return r, nil
}
