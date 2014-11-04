package u2f

type U2F struct {
	Users UserDB
	AppID string
	// base
	AppIDB64 string
	Version  string
}
