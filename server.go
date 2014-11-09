package u2f

// Server
type Server struct {
	Users   UserDB
	AppID   string
	Version string
}

// StdServer standard server
func StdServer(udb UserDB, appID string) Server {
	return Server{
		Users:   udb,
		AppID:   appID,
		Version: "U2F_V2",
	}
}
