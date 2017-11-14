package model

import (
	"time"
)

type Client struct {
	ID     string
	Secret []byte
}

type Permissions struct {
	IsLocked         bool  `db:"is_locked"`
	IsEmailConfirmed bool  `db:"is_email_confirmed"`
	Is2FARequired    bool  `db:"is_2fa_required"`
	Scopes           int64 `db:"scopes"`
}

type DeviceInfo struct {
	ID          string    `db:"device_id"`
	Name        string    `db:"device_name"`
	UserID      string    `db:"user_id"`
	IsConfirmed bool      `db:"is_confirmed"`
	Fingerprint []byte    `db:"fingerprint"`
	LoginTime   time.Time `db:"login_time"`
	LoginIP     string    `db:"login_ip"`
	LoginRegion string    `db:"login_region"`
	Locale      string    `db:"locale"`
	Lang        string    `db:"lang"`
}

type AuthInfo struct {
	Username     string `db:"username"`
	PasswordHash string `db:"password"`
	Permissions
}

type PhoneNumber struct {
	CountryCode int    `db:"phone_country_code"`
	Number      string `db:"phone_number"`
	IsConfirmed bool   `db:"is_phone_confirmed"`
}

type PersonalInfo struct {
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	BirthDate time.Time `db:"birth_date"`
	Country   string    `db:"country"`
	PhoneNumber
}

type UserProfile struct {
	ID string `db:"id"`
	AuthInfo
	PersonalInfo
}
