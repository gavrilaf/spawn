package cache

type Client struct {
	ID     string
	Secret []byte
}

type Session struct {
	ID           string
	RefreshToken string
	ClientID     string
	ClientSecret []byte
	UserID       string
	DeviceID     string
}

type PersonalInfo struct {
	FirstName string
	LastName  string
}

type AuthInfo struct {
	Username     string
	PasswordHash string
	IsLocked     bool
}

type UserProfile struct {
	ID string
	AuthInfo
	PersonalInfo
}

type UserCache interface {
	AddClient(client Client) error
	FindClient(id string) (*Client, error)

	AddSession(session Session) error
	FindSession(id string) (*Session, error)
	DeleteSession(id string) error

	AddUser(profile UserProfile, devices []string) error
	FindProfile(id string) (*UserProfile, error)

	AddDevice(userId string, deviceId string) error
	DeleteDevice(userId string, deviceId string) error
	IsDeviceExists(userId string, deviceId string) (bool, error)
}
