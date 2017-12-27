package dbx

import (
	"time"

	"github.com/gavrilaf/spawn/pkg/dbx/mdl"
	"github.com/gavrilaf/spawn/pkg/errx"
	"github.com/gavrilaf/spawn/pkg/utils"
	"github.com/satori/go.uuid"
)

const (
	getAllClients = `select * from public."Clients"`

	addUser = `INSERT INTO
		public."Users"(id, username, password)
		VALUES($1, $2, $3)`

	getAllUsers = `select * from public."Users"`

	getUserByID = `select * from public."Users" where id = $1`

	getUserByName = `select * from public."Users" where username = $1`

	updatePermission = `UPDATE public."Users"
		SET
			is_locked = $2,
			is_email_confirmed = $3,
			is_2fa_required = $4,
			scope = $5
		WHERE id = $1`

	updatePersonal = `UPDATE public."Users"
			SET
				first_name = $2,
				last_name = $3,
				birth_date = $4
			WHERE id = $1`

	updateCountry = `UPDATE public."Users"
			SET country = $2
			WHERE id = $1`

	updatePhone = `UPDATE public."Users"
		SET
			phone_country_code = $2,
			phone_number = $3,
			is_phone_confirmed = $4
		WHERE id = $1`

	addDevice = `INSERT INTO
		public."Devices"(device_id, user_id, device_name, is_confirmed, locale, lang)
		VALUES($1, $2, $3, $4, $5, $6)`

	updateDevice = `UPDATE public."Devices"
		SET
			device_name = $3,
			locale = $4,
			lang = $5
		WHERE device_id = $1 AND user_id = $2`

	confirmDevice = `UPDATE public."Devices"
		SET is_confirmed = true
		WHERE device_id = $1 AND user_id = $2`

	getUserDevices = `SELECT * FROM public."Devices" WHERE user_id = $1`
	getUserDevice  = `SELECT * FROM public."Devices" WHERE user_id = $1 AND device_id = $2`

	deleteDevice = `DELETE FROM public."Devices"
		WHERE device_id = $1 AND user_id = $2`

	setFingerprint = `UPDATE public."Devices"
		SET fingerprint = $3
		WHERE device_id = $1 AND user_id = $2`

	registerLogin = `INSERT INTO
			public."LoginsLog"(user_id, device_id, device_name, timestamp, user_agent, ip, region)
			VALUES($1, $2, $3, $4, $5, $6, $7)`

	getUserDevicesEx = `select distinct on (device_id)
			D.device_id as device_id, D.user_id as user_id, D.is_confirmed as is_confirmed,
			D.device_name as device_name, D.locale as locale, D.lang as lang,
			LL.timestamp as login_time, LL.ip as login_ip, LL.region as login_region, LL.user_agent as user_agent
		from
			public."Devices" D left join public."LoginsLog" LL on
				D.device_id = LL.device_id AND D.user_id = LL.user_id
		where D.user_id = $1
		order by device_id, login_time`
)

// Clients

func (db *Bridge) GetClients() ([]mdl.Client, error) {
	var clients []mdl.Client
	if err := db.conn.Select(&clients, getAllClients); err != nil {
		return nil, err
	}
	return clients, nil
}

// User profile

// RegisterUser
func (db *Bridge) RegisterUser(username string, password string, device mdl.DeviceInfo) (*mdl.UserProfile, error) {
	userID := uuid.NewV4().String()

	tx, err := db.conn.Beginx()
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(addUser, userID, username, password)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	_, err = tx.Exec(addDevice,
		device.ID,
		userID,
		device.Name,
		device.IsConfirmed,
		device.Locale,
		device.Lang)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return db.GetUserProfile(userID)
}

// GetUserProfile
func (db *Bridge) GetUserProfile(id string) (*mdl.UserProfile, error) {
	var profile mdl.UserProfile
	if err := db.conn.Get(&profile, getUserByID, id); err != nil {
		return nil, errx.ErrNotFound(Scope, "User with id %v not found: %v", id, err)
	}
	profile.BirthDate = utils.FixDbTimezone(profile.BirthDate)
	return &profile, nil
}

// FindUserProfile
func (db *Bridge) FindUserProfile(username string) (*mdl.UserProfile, error) {
	var profile mdl.UserProfile
	if err := db.conn.Get(&profile, getUserByName, username); err != nil {
		return nil, errx.ErrNotFound(Scope, "User with name %v not found: %v", username, err)
	}
	profile.BirthDate = utils.FixDbTimezone(profile.BirthDate)
	return &profile, nil
}

// UpdateUserPermissions
func (db *Bridge) UpdateUserPermissions(id string, permissions mdl.Permissions) error {
	_, err := db.conn.Exec(updatePermission,
		id,
		permissions.IsLocked,
		permissions.IsEmailConfirmed,
		permissions.Is2FARequired,
		permissions.Scope)

	return err
}

// UpdatePersonalInfo
func (db *Bridge) UpdateUserPersonalInfo(id string, info mdl.PersonalInfo) error {
	_, err := db.conn.Exec(updatePersonal,
		id,
		info.FirstName,
		info.LastName,
		info.BirthDate)

	return err
}

// UpdateCountryInfo
func (db *Bridge) UpdateUserCountry(id string, country string) error {
	_, err := db.conn.Exec(updateCountry, id, country)
	return err
}

// UpdateUserPhoneNumber
func (db *Bridge) UpdateUserPhoneNumber(id string, phone mdl.PhoneNumber) error {
	_, err := db.conn.Exec(updatePhone,
		id,
		phone.CountryCode,
		phone.Number,
		phone.IsConfirmed)

	return err

}

// ReadAllUserProfiles
func (db *Bridge) ReadAllUserProfiles() (<-chan *mdl.UserProfile, <-chan error) {
	result := make(chan *mdl.UserProfile)
	errors := make(chan error)

	go func() {
		profile := mdl.UserProfile{}

		rows, err := db.conn.Queryx(getAllUsers)
		if err != nil {
			errors <- err
			return
		}

		for rows.Next() {
			err := rows.StructScan(&profile)
			if err != nil {
				errors <- err
				return
			}
			profile.BirthDate = utils.FixDbTimezone(profile.BirthDate)
			result <- &profile
		}

		result <- nil
	}()

	return result, errors
}

// Devices

// AddDevice
func (db *Bridge) AddDevice(device mdl.DeviceInfo) error {
	_, err := db.conn.Exec(addDevice,
		device.ID,
		device.UserID,
		device.Name,
		device.IsConfirmed,
		device.Locale,
		device.Lang)

	return err
}

// ConfirmDevice
func (db *Bridge) ConfirmDevice(userID string, deviceID string) error {
	_, err := db.conn.Exec(confirmDevice, deviceID, userID)
	return err
}

// UpdateDevice
func (db *Bridge) UpdateDevice(device mdl.DeviceInfo) error {
	_, err := db.conn.Exec(updateDevice,
		device.ID,
		device.UserID,
		device.Name,
		device.Locale,
		device.Lang)

	return err
}

// Remove device
func (db *Bridge) RemoveDevice(userID string, deviceID string) error {
	_, err := db.conn.Exec(deleteDevice, deviceID, userID)
	return err
}

// SetDeviceFingerprint
func (db *Bridge) SetDeviceFingerprint(userID string, deviceID string, fingerprint []byte) error {
	_, err := db.conn.Exec(setFingerprint, deviceID, userID, fingerprint)
	return err
}

// GetUserDevice
func (db *Bridge) GetUserDevice(userID string, deviceID string) (*mdl.DeviceInfo, error) {
	var device mdl.DeviceInfo
	if err := db.conn.Get(&device, getUserDevice, userID, deviceID); err != nil {
		return nil, errx.ErrNotFound(Scope, "Device (%v, %v) not found: %v", userID, deviceID, err)
	}
	return &device, nil
}

// GetUserDevices
func (db *Bridge) GetUserDevices(userID string) ([]mdl.DeviceInfo, error) {
	var devices []mdl.DeviceInfo
	if err := db.conn.Select(&devices, getUserDevices, userID); err != nil {
		return nil, err
	}
	return devices, nil
}

// GetUserDevicesEx
func (db *Bridge) GetUserDevicesEx(userID string) ([]mdl.DeviceInfoEx, error) {
	var devices []mdl.DeviceInfoEx
	if err := db.conn.Select(&devices, getUserDevicesEx, userID); err != nil {
		return nil, err
	}
	return devices, nil
}

// User logs

// LogUserLogin
func (db *Bridge) LogUserLogin(userID string, deviceID string, userAgent string, ip string, region string) error {
	device, err := db.GetUserDevice(userID, deviceID)
	if device == nil {
		return errx.ErrNotFound(Scope, "Device (%v, %v) not found. Error: %v", userID, deviceID, err)
	}
	_, err = db.conn.Exec(registerLogin, userID, deviceID, device.Name, time.Now(), userAgent, ip, region)
	return err
}
