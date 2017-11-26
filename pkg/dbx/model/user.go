package model

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
)

func emptyBirthDate() time.Time {
	t, _ := time.Parse(time.RFC3339, "1900-01-01T00:00:00+00:00")
	return t.UTC()
}

func BirthdayDate(year int, month int, day int) time.Time {
	t, _ := time.Parse(time.RFC3339, fmt.Sprintf("%d-%02d-%02dT00:00:00+00:00", year, month, day))
	return t.UTC()
}

func BirthdayDateFromTime(t time.Time) time.Time {
	return t.Truncate(24 * time.Hour).UTC()
}

var EmptyBirthDate = emptyBirthDate()

///

type Client struct {
	ID     string
	Secret []byte
}

// Device

type DeviceInfo struct {
	ID          string `db:"device_id"`
	Name        string `db:"device_name"`
	UserID      string `db:"user_id"`
	IsConfirmed bool   `db:"is_confirmed"`
	Fingerprint []byte `db:"fingerprint"`
	Locale      string `db:"locale"`
	Lang        string `db:"lang"`
}

type DeviceInfoEx struct {
	LoginTime   pq.NullTime    `db:"login_time"`
	LoginIP     sql.NullString `db:"login_ip"`
	LoginRegion sql.NullString `db:"login_region"`
	DeviceInfo
}

func (d DeviceInfoEx) GetLoginTime() *time.Time {
	if d.LoginTime.Valid {
		return &d.LoginTime.Time
	}
	return nil
}

func (d DeviceInfoEx) GetLoginIP() string {
	if d.LoginIP.Valid {
		return d.LoginIP.String
	}
	return ""
}

func (d DeviceInfoEx) GetLoginRegion() string {
	if d.LoginIP.Valid {
		return d.LoginIP.String
	}
	return ""
}

// User profile

type Permissions struct {
	IsLocked         bool  `db:"is_locked" structs:"is_locked"`
	IsEmailConfirmed bool  `db:"is_email_confirmed" structs:"is_email_confirmed"`
	Is2FARequired    bool  `db:"is_2fa_required" structs:"is_2fa_required"`
	Scopes           int64 `db:"scopes" structs:"scopes"`
}

type AuthInfo struct {
	Username     string `db:"username" structs:"username"`
	PasswordHash string `db:"password" structs:"password"`
	Permissions
}

type PhoneNumber struct {
	CountryCode int    `db:"phone_country_code" structs:"phone_country_code"`
	Number      string `db:"phone_number" structs:"phone_number"`
	IsConfirmed bool   `db:"is_phone_confirmed" structs:"is_phone_confirmed"`
}

type PersonalInfo struct {
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	BirthDate time.Time `db:"birth_date"`
}

type UserProfile struct {
	ID      string `db:"id"`
	Country string `db:"country"`
	PhoneNumber
	AuthInfo
	PersonalInfo
}

func (p *UserProfile) GetBirthdayDate() time.Time {
	return BirthdayDateFromTime(p.BirthDate)
}

// User logs

type LoginInfo struct {
	UserID   string    `db:"user_id"`
	DeviceID string    `db:"device_id"`
	Time     time.Time `db:"time"`
	IP       string    `db:"ip"`
	Region   string    `db:"region"`
}
