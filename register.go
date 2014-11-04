package u2f

// RegisterRequest dictionary from the fido u2f javascript api spec
type RegisterRequest struct {
	Version   string `json:"version"`
	Challenge string `json:"challenge"`
	AppID     string `json:"appId"`
}

// Register a user to a device. Returns a RegisterRequest Object for the device
// to sign. The result of which is passed to RegisterFin().
func (f U2F) Register(u User) (r RegisterRequest, err error) {
	c, err := challenge()
	if err != nil {
		return r, err
	}

	u.Devices = append(u.Devices, Device{
		Challenge: c,
	})
	err = f.Users.PutUser(u)
	if err != nil {
		return r, err
	}

	r = RegisterRequest{
		Version:   f.Version,
		Challenge: c,
		AppID:     f.AppID,
	}
	return r, nil
}
