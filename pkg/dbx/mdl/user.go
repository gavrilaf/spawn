package mdl

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
)

///

type Client struct {
	ID          string `db:"id" structs:"id"`
	Secret      []byte `db:"secret" structs:"secret"`
	IsActive    bool   `db:"is_active" structs:"is_active"`
	Description string `db:"description" structs:"description"`
	DefScope    int64  `db:"def_scope" structs:"def_scope"`
}

// Device

type DeviceInfo struct {
	ID          string `db:"device_id" structs:"device_id"`
	Name        string `db:"device_name" structs:"device_name"`
	UserID      string `db:"user_id" structs:"user_id"`
	IsConfirmed bool   `db:"is_confirmed" structs:"is_confirmed"`
	Fingerprint []byte `db:"fingerprint" structs:"fingerprint"`
	Locale      string `db:"locale" structs:"locale"`
	Lang        string `db:"lang" structs:"lang"`
}

type DeviceInfoEx struct {
	LoginTime   pq.NullTime    `db:"login_time"  structs:"login_time"`
	LoginIP     sql.NullString `db:"login_ip"  structs:"login_ip"`
	UserAgent   sql.NullString `db:"user_agent"  structs:"user_agent"`
	LoginRegion sql.NullString `db:"login_region"  structs:"login_region"`
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

func (d DeviceInfoEx) GetUserAgent() string {
	if d.UserAgent.Valid {
		return d.UserAgent.String
	}
	return ""
}

func (d DeviceInfoEx) GetLoginRegion() string {
	if d.LoginIP.Valid {
		return d.LoginRegion.String
	}
	return ""
}

// User profile

type Permissions struct {
	IsLocked         bool  `db:"is_locked" structs:"is_locked"`
	IsEmailConfirmed bool  `db:"is_email_confirmed" structs:"is_email_confirmed"`
	Is2FARequired    bool  `db:"is_2fa_required" structs:"is_2fa_required"`
	Scope            int64 `db:"scope" structs:"scope"`
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
	FirstName string    `db:"first_name" structs:"first_name"`
	LastName  string    `db:"last_name" structs:"last_name"`
	BirthDate time.Time `db:"birth_date" structs:"birth_date"`
}

type UserProfile struct {
	ID      string `db:"id" structs:"id"`
	Country string `db:"country" structs:"country"`
	PhoneNumber
	AuthInfo
	PersonalInfo
}

// User logs

type LoginInfo struct {
	UserID    string    `db:"user_id" structs:"user_id"`
	DeviceID  string    `db:"device_id" structs:"device_id"`
	Time      time.Time `db:"timestamp" structs:"timestamp"`
	IP        string    `db:"ip" structs:"ip"`
	UserAgent string    `db:"user_agent" structs:"user_agent"`
	Region    string    `db:"region" structs:"region"`
}
