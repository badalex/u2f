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

func (u2f U2F) validateClientData(typ, clientData, challenge string) error {
	if clientData == "" {
		return fmt.Errorf("Missing ClientData")
	}

	data, err := base64.URLEncoding.DecodeString(clientData)
	if err != nil {
		return err
	}

	cd := clientDataJSON{}
	err = json.Unmarshal(data, &cd)
	if err != nil {
		return err
	}

	if cd.Typ != typ {
		return fmt.Errorf("Typ should be %s", typ)
	}
	if cd.Challenge != challenge {
		return fmt.Errorf("challenges dont match")
	}
	if cd.Origin != u2f.AppID {
		return fmt.Errorf("Origin does not match appID")
	}

	return nil
}
