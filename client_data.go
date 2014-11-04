package u2f

import (
	"encoding/json"
	"fmt"
)

type clientData struct {
	Typ       string
	Origin    string
	Challenge string
}

func (f U2F) validateClientData(typ, cd string, devs []Device) (dev *Device, err error) {
	if cd == "" {
		return dev, fmt.Errorf("Missing ClientData")
	}

	data, err := unb64u(cd)
	if err != nil {
		return dev, err
	}

	c := clientData{}
	err = json.Unmarshal(data, &c)
	if err != nil {
		return dev, err
	}

	if c.Typ != typ {
		return dev, fmt.Errorf("Typ should be %s", typ)
	}
	if c.Origin != f.AppID {
		return dev, fmt.Errorf("Origin does not match appID")
	}

	for idx := range devs {
		if c.Challenge == devs[idx].Challenge {
			return &devs[idx], nil
		}
	}

	return dev, fmt.Errorf("no matching challenge found")
}
