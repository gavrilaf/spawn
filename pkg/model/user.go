package model

import (
	"time"
)

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

type DeviceInfo struct {
	ID          string
	Name        string
	Fingerprint []byte
	LoginTime   time.Time
	LoginIP     string
	LoginRegion string
}

type AuthInfo struct {
	Username         string
	PasswordHash     string
	IsLocked         bool
	IsEmailConfirmed bool
	Is2FARequired    bool
}

type PhoneNumber struct {
	CountryCode int
	Number      string
	IsConfirmed bool
}

type PersonalInfo struct {
	FirstName string
	LastName  string
	BirthDate int64
	PhoneNumber
}

type UserProfile struct {
	ID string
	AuthInfo
	PersonalInfo
}
