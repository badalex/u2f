package u2f

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type clientDataJSON struct {
	Typ       string
	Origin    string
	Challenge string
}

func (u2f U2F) validateClientData(typ, clientData string, devs []Device) (dev *Device, err error) {
	if clientData == "" {
		return dev, fmt.Errorf("Missing ClientData")
	}

	data, err := base64.URLEncoding.DecodeString(clientData)
	if err != nil {
		return dev, err
	}

	cd := clientDataJSON{}
	err = json.Unmarshal(data, &cd)
	if err != nil {
		return dev, err
	}

	if cd.Typ != typ {
		return dev, fmt.Errorf("Typ should be %s", typ)
	}
	if cd.Origin != u2f.AppID {
		return dev, fmt.Errorf("Origin does not match appID")
	}

	for idx := range devs {
		if cd.Challenge == devs[idx].Challenge {
			return &devs[idx], nil
		}
	}

	return dev, fmt.Errorf("no matching challenge found")
}
