package dbx

import (
	//"database/sql"
	"fmt"
	mdl "github.com/gavrilaf/spawn/pkg/model"
	"github.com/satori/go.uuid"
)

const (
	addUserQuery = `INSERT INTO 
		public."Users"(id, username, password, is_locked, is_email_confirmed, is_2fa_required) 
		VALUES($1, $2, $3, $4, $5, $6)`

	addDeviceQuery = `INSERT INTO 
		public."Devices"(device_id, user_id, device_name, is_confirmed, fingerprint, login_time, login_ip, login_region, locale, lang) 
		VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	confirmDeviceQuery = `UPDATE public."Devices"
		SET is_confirmed = true 
		WHERE device_id = $1 AND user_id = $2`

	deleteDeviceQuery = `DELETE FROM public."Devices"
		WHERE device_id = $1 AND user_id = $2`
)

func (db *Bridge) RegisterUser(username string, password string, device mdl.DeviceInfo) (*mdl.UserProfile, error) {

	// generate id
	userId := uuid.NewV4().String()

	tx := db.Db.MustBegin()

	tx.MustExec(addUserQuery, userId, username, password, false, false, false)
	tx.MustExec(addDeviceQuery,
		device.ID,
		userId,
		device.Name,
		device.IsConfirmed,
		device.Fingerprint,
		device.LoginTime,
		device.LoginIP,
		device.LoginRegion,
		device.Locale,
		device.Lang)

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return db.GetUserProfile(userId)
}

func (db *Bridge) UpdatePermissions(id string, permissons *mdl.Permissions) error {
	return fmt.Errorf("not implemented")
}

func (db *Bridge) UpdatePersonalInfo(id string, info *mdl.PersonalInfo) error {
	return fmt.Errorf("not implemented")
}

func (db *Bridge) GetUserProfile(id string) (*mdl.UserProfile, error) {
	var profile mdl.UserProfile
	if err := db.Db.Get(&profile, `select * from public."Users" where id = $1`, id); err != nil {
		return nil, err
	}

	return &profile, nil
}

func (db *Bridge) FindUserProfile(username string) (*mdl.UserProfile, error) {
	var profile mdl.UserProfile
	if err := db.Db.Get(&profile, `select * from public."Users" where username = $1`, username); err != nil {
		return nil, err
	}

	return &profile, nil
}

func (db *Bridge) AddDevice(userId string, device mdl.DeviceInfo) error {
	_, err := db.Db.Exec(addDeviceQuery,
		device.ID,
		userId,
		device.Name,
		device.IsConfirmed,
		device.Fingerprint,
		device.LoginTime,
		device.LoginIP,
		device.LoginRegion,
		device.Locale,
		device.Lang)

	return err
}

func (db *Bridge) ConfirmDevice(userId string, deviceId string) error {
	_, err := db.Db.Exec(confirmDeviceQuery, deviceId, userId)
	return err
}

func (db *Bridge) RemoveDevice(userId string, deviceId string) error {
	_, err := db.Db.Exec(deleteDeviceQuery, deviceId, userId)
	return err
}

func (db *Bridge) GetDevices(userId string) ([]mdl.DeviceInfo, error) {
	var devices []mdl.DeviceInfo

	if err := db.Db.Select(&devices, `select * from public."Devices" where user_id = $1`, userId); err != nil {
		return nil, err
	}

	return devices, nil
}
