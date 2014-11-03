package u2f

type User struct {
	User      string
	Enrolled  bool
	KeyHandle string
	PubKey    string
	Cert      string
	Challenge string
	Counter   uint32
}

type Users interface {
	GetUser(user string) (User, error)
	PutUser(u User)
}
