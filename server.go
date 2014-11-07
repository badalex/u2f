package u2f

// U2FServer
type U2FServer struct {
	Users   UserDB
	AppID   string
	Version string
}

// StdU2FServer standard server
func StdU2FServer(udb UserDB, appID string) U2FServer {
	return U2FServer{
		Users:   udb,
		AppID:   appID,
		Version: "U2F_V2",
	}
}
