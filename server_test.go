package u2f

import (
	"fmt"
	"strings"
	"sync"
	"testing"
)

type userDB struct {
	Users map[string]User
	lock  sync.Mutex
}

func (ud *userDB) GetUser(user string) (User, error) {
	ud.lock.Lock()
	defer ud.lock.Unlock()

	if ud.Users == nil {
		ud.Users = make(map[string]User)
		return User{}, fmt.Errorf("no such user")
	}

	u, ok := ud.Users[user]
	if !ok {
		return u, fmt.Errorf("no such user")
	}
	return u, nil
}

func (ud *userDB) PutUser(u User) error {
	ud.lock.Lock()
	defer ud.lock.Unlock()

	if ud.Users == nil {
		ud.Users = make(map[string]User)
	}
	if u.User == "" {
		return fmt.Errorf("No username passed")
	}
	ud.Users[u.User] = u
	return nil
}

func TestAll(t *testing.T) {
	var udb = userDB{}
	var s = Server{
		Users:   &udb,
		AppID:   "http://demo.example.com",
		Version: "U2F_V2",
	}

	s.Users.PutUser(User{User: "test"})
	u := getUser(&s, t)

	_, err := s.Register(u)
	if err != nil {
		t.Fatal(err)
	}
	u = getUser(&s, t)

	_, err = s.Register(u)
	if err != nil {
		t.Fatal(err)
	}
	u = getUser(&s, t)

	u.Devices[0].Challenge = "yKA0x075tjJ-GE7fKTfnzTOSaNUOWQxRd9TWz5aFOg8"
	s.Users.PutUser(u)

	r := strings.NewReader(` {"registrationData": "BQQtEmhWVgvbh-8GpjsHbj_d5FB9iNoRL8mNEq34-ANufKWUpVdIj6BSB_m3eMoZ3GqnaDy3RA5eWP8mhTkT1Ht3QAk1GsmaPIQgXgvrBkCQoQtMFvmwYPfW5jpRgoMPFxquHS7MTt8lofZkWAK2caHD-YQQdaRBgd22yWIjPuWnHOcwggLiMIHLAgEBMA0GCSqGSIb3DQEBCwUAMB0xGzAZBgNVBAMTEll1YmljbyBVMkYgVGVzdCBDQTAeFw0xNDA1MTUxMjU4NTRaFw0xNDA2MTQxMjU4NTRaMB0xGzAZBgNVBAMTEll1YmljbyBVMkYgVGVzdCBFRTBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IABNsK2_Uhx1zOY9ym4eglBg2U5idUGU-dJK8mGr6tmUQflaNxkQo6IOc-kV4T6L44BXrVeqN-dpCPr-KKlLYw650wDQYJKoZIhvcNAQELBQADggIBAJVAa1Bhfa2Eo7TriA_jMA8togoA2SUE7nL6Z99YUQ8LRwKcPkEpSpOsKYWJLaR6gTIoV3EB76hCiBaWN5HV3-CPyTyNsM2JcILsedPGeHMpMuWrbL1Wn9VFkc7B3Y1k3OmcH1480q9RpYIYr-A35zKedgV3AnvmJKAxVhv9GcVx0_CewHMFTryFuFOe78W8nFajutknarupekDXR4tVcmvj_ihJcST0j_Qggeo4_3wKT98CgjmBgjvKCd3Kqg8n9aSDVWyaOZsVOhZj3Fv5rFu895--D4qiPDETozJIyliH-HugoQpqYJaTX10mnmMdCa6aQeW9CEf-5QmbIP0S4uZAf7pKYTNmDQ5z27DVopqaFw00MIVqQkae_zSPX4dsNeeoTTXrwUGqitLaGap5ol81LKD9JdP3nSUYLfq0vLsHNDyNgb306TfbOenRRVsgQS8tJyLcknSKktWD_Qn7E5vjOXprXPrmdp7g5OPvrbz9QkWa1JTRfo2n2AXV02LPFc-UfR9bWCBEIJBxvmbpmqt0MnBTHWnth2b0CU_KJTDCY3kAPLGbOT8A4KiI73pRW-e9SWTaQXskw3Ei_dHRILM_l9OXsqoYHJ4Dd3tbfvmjoNYggSw4j50l3unI9d1qR5xlBFpW5sLr8gKX4bnY4SR2nyNiOQNLyPc0B0nW502aMEUCIQDTGOX-i_QrffJDY8XvKbPwMuBVrOSO-ayvTnWs_WSuDQIgZ7fMAvD_Ezyy5jg6fQeuOkoJi8V2naCtzV-HTly8Nww=", "clientData": "eyAiY2hhbGxlbmdlIjogInlLQTB4MDc1dGpKLUdFN2ZLVGZuelRPU2FOVU9XUXhSZDlUV3o1YUZPZzgiLCAib3JpZ2luIjogImh0dHA6XC9cL2RlbW8uZXhhbXBsZS5jb20iLCAidHlwIjogIm5hdmlnYXRvci5pZC5maW5pc2hFbnJvbGxtZW50IiB9"}`)
	err = s.RegisterFin(u, r)
	if err != nil {
		t.Fatal(err)
	}
	u = getUser(&s, t)

	_, err = s.Sign(u)
	if err != nil {
		t.Fatal(err)
	}

	u = getUser(&s, t)
	if u.Devices[0].Challenge == "" {
		t.Fatal("failed to sign")
	}

	u.Devices[0].Challenge = "fEnc9oV79EaBgK5BoNERU5gPKM2XGYWrz4fUjgc0Q7g"
	s.Users.PutUser(u)

	r = strings.NewReader(`{ "signatureData": "AQAAAAQwRQIhAI6FSrMD3KUUtkpiP0jpIEakql-HNhwWFngyw553pS1CAiAKLjACPOhxzZXuZsVO8im-HStEcYGC50PKhsGp_SUAng==", "clientData": "eyAiY2hhbGxlbmdlIjogImZFbmM5b1Y3OUVhQmdLNUJvTkVSVTVnUEtNMlhHWVdyejRmVWpnYzBRN2ciLCAib3JpZ2luIjogImh0dHA6XC9cL2RlbW8uZXhhbXBsZS5jb20iLCAidHlwIjogIm5hdmlnYXRvci5pZC5nZXRBc3NlcnRpb24iIH0=", "keyHandle": "CTUayZo8hCBeC-sGQJChC0wW-bBg99bmOlGCgw8XGq4dLsxO3yWh9mRYArZxocP5hBB1pEGB3bbJYiM-5acc5w" }`)
	_, err = s.SignFin(u, r)
	if err != nil {
		t.Fatal(err)
	}
	u = getUser(&s, t)

	r = strings.NewReader(`{ "signatureData": "AQAAAAQwRQIhAI6FSrMD3KUUtkpiP0jpIEakql-HNhwWFngyw553pS1CAiAKLjACPOhxzZXuZsVO8im-HStEcYGC50PKhsGp_SUAng==", "clientData": "eyAiY2hhbGxlbmdlIjogImZFbmM5b1Y3OUVhQmdLNUJvTkVSVTVnUEtNMlhHWVdyejRmVWpnYzBRN2ciLCAib3JpZ2luIjogImh0dHA6XC9cL2RlbW8uZXhhbXBsZS5jb20iLCAidHlwIjogIm5hdmlnYXRvci5pZC5nZXRBc3NlcnRpb24iIH0=", "keyHandle": "CTUayZo8hCBeC-sGQJChC0wW-bBg99bmOlGCgw8XGq4dLsxO3yWh9mRYArZxocP5hBB1pEGB3bbJYiM-5acc5w" }`)
	_, err = s.SignFin(u, r)
	if err == nil {
		t.Fatal("The counters should have mismatched!")
	}
}

func getUser(s *Server, t *testing.T) User {
	u, err := s.Users.GetUser("test")
	if err != nil {
		t.Fatal(err)
	}
	return u

}

func TestClientData(t *testing.T) {
	var err error
	s := StdServer(nil, "")

	_, err = s.validateClientData("", "", nil)
	nonNil(t, err)

	// not base64
	_, err = s.validateClientData("", "{}", nil)
	nonNil(t, err)

	// not json
	_, err = s.validateClientData("", b64u([]byte("asdf")), nil)
	nonNil(t, err)

	// typ mismatch
	_, err = s.validateClientData("", b64u([]byte(`{"typ": "test"}`)), nil)
	nonNil(t, err)

	// origin mismatch
	_, err = s.validateClientData("test", b64u([]byte(`{"typ": "test"}`)), nil)
	nonNil(t, err)

	// origin mismatch
	_, err = s.validateClientData("test", b64u([]byte(`{"typ": "test", "origin": "asdf"}`)), nil)
	nonNil(t, err)

	// no devices
	_, err = s.validateClientData("test", b64u([]byte(`{"typ": "test", "origin": ""}`)), nil)
	nonNil(t, err)

	// no matching challenge
	_, err = s.validateClientData("test", b64u([]byte(`{typ: "test", "origin": "", "challenge": "a"}`)), []Device{Device{Challenge: "b"}})
	nonNil(t, err)

	// no devices
	_, err = s.validateClientData("test", b64u([]byte(`{typ: "test", "origin": "", "challenge": "a"}`)), []Device{Device{Challenge: "a"}})
	nonNil(t, err)
}

func TestRegister(t *testing.T) {
	var err error
	udb := &userDB{}
	s := StdServer(nil, "test")
	s.Users = udb

	// bad challenge
	_, err = s.Register(User{})
	nonNil(t, err)

	// no username, PutUser fails
	_, err = s.Register(User{})
	nonNil(t, err)

	r, err := s.Register(User{User: "test"})
	if err != nil {
		t.Fatal(err)
	}

	if r.Version != s.Version {
		t.Fatal("version mismatch")
	}
	if r.Challenge == "" {
		t.Fatal("no challenge")
	}
	if r.AppID != "test" {
		t.Fatal("AppId mismatch")
	}

	u, err := s.Users.GetUser("test")
	if err != nil {
		t.Fatal(err)
	}

	if u.User != "test" {
		t.Fatal("failed to put user")
	}
	if len(udb.Users["test"].Devices) != 1 {
		t.Fatal("strange number of devices")
	}
	if udb.Users["test"].Devices[0].Challenge != r.Challenge {
		t.Fatal("challenge mismatch")
	}
}

func nonNil(t *testing.T, n interface{}) {
	if n == nil {
		panic("should be nil")
		t.Fatalf("Should be nil: %+v", n)
	}
}
