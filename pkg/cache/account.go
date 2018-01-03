package cache

import (
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/gavrilaf/spawn/pkg/cache/mdl"
	db "github.com/gavrilaf/spawn/pkg/dbx/mdl"
	"github.com/gavrilaf/spawn/pkg/errx"
)

func (br *Bridge) SetUserAccounts(userID string, accounts []db.Account) error {
	conn := br.get()
	defer conn.Close()

	prefix := "account:" + userID + ":"

	for _, dbAcc := range accounts {
		key := prefix + dbAcc.ID
		account := mdl.CreateAcountFromDbModel(dbAcc)
		_, err := conn.Do("HMSET", redis.Args{}.Add(key).AddFlat(&account)...)
		if err != nil {
			removeAllUserAccount(conn, userID)
			return err
		}
	}

	return nil
}

func (br *Bridge) GetUserAccounts(userID string) ([]mdl.Account, error) {
	conn := br.get()
	defer conn.Close()

	pattern := "account:" + userID + ":" + "*"

	keys, err := getKeys(conn, pattern)
	if err != nil {
		return nil, err
	}

	accounts := make([]mdl.Account, len(keys))

	for indx, key := range keys {
		v, err := redis.Values(conn.Do("HGETALL", key))

		if err != nil {
			return nil, err
		}
		if len(v) == 0 {
			return nil, errx.ErrKeyNotFound(Scope, key)
		}

		var account mdl.Account
		if err := redis.ScanStruct(v, &account); err != nil {
			return nil, err
		}
		accounts[indx] = account
	}

	return accounts, nil
}

func (br *Bridge) UpdateUserAccountStatus(userID string, accountID string, status mdl.AccountStatus) error {
	conn := br.get()
	defer conn.Close()

	return updateAccount(conn, userID, accountID, func(account *mdl.Account) *mdl.Account {
		account.Status = status
		return account
	})
}

func (br *Bridge) UpdateUserAccountBalance(userID string, accountID string, balance string) error {
	conn := br.get()
	defer conn.Close()

	return updateAccount(conn, userID, accountID, func(account *mdl.Account) *mdl.Account {
		account.Balance = balance
		return account
	})
}

func (br *Bridge) AddUserAccount(userID string, account db.Account) error {
	conn := br.get()
	defer conn.Close()

	key := "account:" + userID + ":" + account.ID
	p := mdl.CreateAcountFromDbModel(account)
	_, err := conn.Do("HMSET", redis.Args{}.Add(key).AddFlat(&p)...)

	return err
}

func (br *Bridge) GetUserAccount(userID string, accountID string) (*mdl.Account, error) {
	conn := br.get()
	defer conn.Close()

	return getAccount(conn, userID, accountID)
}

func (br *Bridge) ClearUserAccounts(userID string) error {
	conn := br.get()
	defer conn.Close()

	return removeAllUserAccount(conn, userID)
}

/////////////////////////////////////////////////////////////////////////////////
// Private

func removeAllUserAccount(conn redis.Conn, userID string) error {
	_, err := deleteKeys(conn, "account:"+userID+":"+"*")
	return err
}

func getAccount(conn redis.Conn, userID string, accountID string) (*mdl.Account, error) {
	key := "account:" + userID + ":" + accountID

	v, err := redis.Values(conn.Do("HGETALL", key))

	if err != nil {
		return nil, err
	}
	if len(v) == 0 {
		return nil, errx.ErrKeyNotFound(Scope, key)
	}

	var account mdl.Account
	if err := redis.ScanStruct(v, &account); err != nil {
		return nil, err
	}

	return &account, nil
}

func updateAccount(conn redis.Conn, userID string, accountID string, updater func(account *mdl.Account) *mdl.Account) error {
	account, err := getAccount(conn, userID, accountID)
	if err != nil {
		return err
	}

	account = updater(account)
	account.Updated = time.Now().UTC().Unix()

	key := "account:" + userID + ":" + accountID
	_, err = conn.Do("HMSET", redis.Args{}.Add(key).AddFlat(account)...)

	return err
}
