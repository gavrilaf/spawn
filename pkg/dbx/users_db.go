package dbx

import (
	"fmt"
	"time"

	mdl "github.com/gavrilaf/spawn/pkg/dbx/model"
	"github.com/satori/go.uuid"
)

func EmptyBirthDate() time.Time {
	t1, _ := time.Parse(time.RFC3339, "1900-01-01T00:00:00+00:00")
	return t1
}

//time.Date(1900, time.January, 0, 0, 0, 0, 0, time.UTC)

const (
	addUser = `INSERT INTO
		public."Users"(id, username, password)
		VALUES($1, $2, $3)`

	getAllUsers = `select * from public."Users"`

	getUserByID = `select * from public."Users" where id = $1`

	getUserByName = `select * from public."Users" where username = $1`

	updatePermission = `UPDATE public."Users"
		SET
			is_locked = $1,
			is_email_confirmed = $2,
			is_2fa_required = $3,
			scopes = $4
		WHERE id = $5`

	updatePersonal = `UPDATE public."Users"
			SET
				first_name = $1,
				last_name = $2,
				birth_date = $3
			WHERE id = $4`

	updateCountry = `UPDATE public."Users"
			SET country = $1
			WHERE id = $2`

	updatePhone = `UPDATE public."Users"
		SET
			phone_country_code = $1,
			phone_number = $2,
			is_phone_confirmed = $3
		WHERE id = $4`

	addDevice = `INSERT INTO
		public."Devices"(device_id, user_id, device_name, is_confirmed, locale, lang)
		VALUES($1, $2, $3, $4, $5, $6)`

	confirmDevice = `UPDATE public."Devices"
		SET is_confirmed = true
		WHERE device_id = $1 AND user_id = $2`

	getUserDevices = `SELECT * FROM public."Devices" WHERE user_id = $1`
	getUserDevice  = `SELECT * FROM public."Devices" WHERE user_id = $1 AND device_id = $2`

	deleteDevice = `DELETE FROM public."Devices"
		WHERE device_id = $1 AND user_id = $2`

	setFingerprint = `UPDATE public."Devices"
		SET fingerprint = $1,
		WHERE device_id = $2 AND user_id = $3`

	registerLogin = `INSERT INTO
			public."LoginsLog"(user_id, device_id, device_name, time, ip, region)
			VALUES($1, $2, $3, $4, $5, $6)`

	getUserDevicesEx = `select distinct on (device_id)
			D.device_id as device_id, D.user_id as user_id,
			D.device_name as device_name,
			LL.time as login_time, LL.ip as login_ip, LL.region as login_region
		from
			public."Devices" D left join public."LoginsLog" LL on
				D.device_id = LL.device_id AND D.user_id = LL.user_id
		where D.user_id = $1
		order by device_id, login_time`
)

// User profile

func (db *Bridge) RegisterUser(username string, password string, device mdl.DeviceInfo) (*mdl.UserProfile, error) {
	userID := uuid.NewV4().String()

	tx := db.Db.MustBegin()
	tx.MustExec(addUser, userID, username, password)
	tx.MustExec(addDevice,
		device.ID,
		userID,
		device.Name,
		device.IsConfirmed,
		device.Locale,
		device.Lang)

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return db.GetUserProfile(userID)
}

func (db *Bridge) GetUserProfile(id string) (*mdl.UserProfile, error) {
	var profile mdl.UserProfile
	if err := db.Db.Get(&profile, getUserByID, id); err != nil {
		return nil, err
	}

	return &profile, nil
}

func (db *Bridge) FindUserProfile(username string) (*mdl.UserProfile, error) {
	var profile mdl.UserProfile
	if err := db.Db.Get(&profile, getUserByName, username); err != nil {
		return nil, err
	}

	return &profile, nil
}

func (db *Bridge) UpdateUserPermissions(id string, permissions mdl.Permissions) error {
	_, err := db.Db.Exec(updatePermission,
		permissions.IsLocked,
		permissions.IsEmailConfirmed,
		permissions.Is2FARequired,
		permissions.Scopes,
		id)

	return err
}

func (db *Bridge) UpdateUserPersonalInfo(id string, info mdl.PersonalInfo) error {
	_, err := db.Db.Exec(updatePersonal,
		info.FirstName,
		info.LastName,
		info.BirthDate,
		id)

	return err
}

func (db *Bridge) UpdateUserCountry(id string, country string) error {
	_, err := db.Db.Exec(updateCountry, country, id)
	return err
}

func (db *Bridge) UpdateUserPhoneNumber(id string, phone mdl.PhoneNumber) error {
	_, err := db.Db.Exec(updatePhone,
		phone.CountryCode,
		phone.Number,
		phone.IsConfirmed,
		id)
	return err

}

func (db *Bridge) ReadAllUserProfiles() (<-chan *mdl.UserProfile, <-chan error) {
	result := make(chan *mdl.UserProfile)
	errors := make(chan error)

	go func() {
		profile := mdl.UserProfile{}

		rows, err := db.Db.Queryx(getAllUsers)
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
			result <- &profile
		}

		result <- nil
	}()

	return result, errors
}

// Devices

func (db *Bridge) AddDevice(userId string, device mdl.DeviceInfo) error {
	_, err := db.Db.Exec(addDevice,
		device.ID,
		userId,
		device.Name,
		device.IsConfirmed,
		device.Locale,
		device.Lang)

	return err
}

func (db *Bridge) ConfirmDevice(userId string, deviceId string) error {
	_, err := db.Db.Exec(confirmDevice, deviceId, userId)
	return err
}

func (db *Bridge) RemoveDevice(userId string, deviceId string) error {
	_, err := db.Db.Exec(deleteDevice, deviceId, userId)
	return err
}

func (db *Bridge) SetDeviceFingerprint(userID string, deviceID string, fingerprint []byte) error {
	_, err := db.Db.Exec(setFingerprint, deviceID, userID, fingerprint)
	return err
}

func (db *Bridge) GetUserDevice(userID string, deviceID string) (*mdl.DeviceInfo, error) {
	var device mdl.DeviceInfo
	if err := db.Db.Get(&device, getUserDevice, userID, deviceID); err != nil {
		return nil, err
	}
	return &device, nil
}

func (db *Bridge) GetUserDevices(userID string) ([]mdl.DeviceInfo, error) {
	var devices []mdl.DeviceInfo
	if err := db.Db.Select(&devices, getUserDevices, userID); err != nil {
		return nil, err
	}
	return devices, nil
}

func (db *Bridge) GetUserDevicesEx(userID string) ([]mdl.DeviceInfoEx, error) {
	var devices []mdl.DeviceInfoEx
	if err := db.Db.Select(&devices, getUserDevicesEx, userID); err != nil {
		return nil, err
	}
	return devices, nil
}

// User logs

func (db *Bridge) LogUserLogin(userID string, deviceID string, ip string, region string) error {
	device, err := db.GetUserDevice(userID, deviceID)
	if device == nil {
		return fmt.Errorf("Device (%v, %v) not found", userID, deviceID)
	}
	_, err = db.Db.Exec(registerLogin, userID, deviceID, device.Name, time.Now(), ip, region)
	return err
}
