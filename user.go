package u2f

// A User
type User struct {
	// Priv is yours to do what you please with.
	// For example, if you have an sql backend you could store the tuple or
	// primary key here to make updating easier.
	Priv interface{}

	// User containts the username
	User string

	// Enrolled holds if they have signed up
	Enrolled bool

	// A list of associated U2F devices with for this User
	Devices []Device
}

// Holds a U2F device
type Device struct {
	//  Priv is yours to do what you please with
	Priv      interface{}
	KeyHandle string
	PubKey    string
	Cert      string
	Challenge string
	Counter   uint32
}

// UserDB interface
type UserDB interface {
	// GetUser from a username. It is assumed you have done any needed
	// password authentication before this point
	GetUser(user string) (User, error)

	// PutUser Update the user
	PutUser(u User) error
}
