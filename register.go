package u2f

// RegisterRequest dictionary from the fido u2f javascript api spec
type RegisterRequest struct {
	Version   string `json:"version"`
	Challenge string `json:"challenge"`
	AppID     string `json:"appId"`
}

// Register a user to a device. Returns a RegisterRequest Object for the device
// to sign. The result of which is passed to RegisterFin().
func (s U2FServer) Register(u User) (r RegisterRequest, err error) {
	c, err := s.Challenge.New()
	if err != nil {
		return r, err
	}

	u.Devices = append(u.Devices, Device{
		Challenge: c,
	})
	err = s.Users.PutUser(u)
	if err != nil {
		return r, err
	}

	r = RegisterRequest{
		Version:   s.Version,
		Challenge: c,
		AppID:     s.AppID,
	}
	return r, nil
}
