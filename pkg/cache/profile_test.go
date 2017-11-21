package cache

import (
	//"fmt"
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	db "github.com/gavrilaf/spawn/pkg/dbx/model"
	"github.com/gavrilaf/spawn/pkg/env"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getUserProfileCache(t *testing.T) Cache {
	cache, err := Connect(env.GetEnvironment("Test"))
	require.Nil(t, err)
	require.NotNil(t, cache)
	return cache
}

func TestUserProfile(t *testing.T) {
	cache := getUserProfileCache(t)
	defer cache.Close()

	profile := db.UserProfile{
		ID:      "user-1",
		Country: "ua",
		PhoneNumber: db.PhoneNumber{
			CountryCode: 38,
			Number:      "067876123",
			IsConfirmed: false,
		},
		AuthInfo: db.AuthInfo{
			Username:     "testuser@test.com",
			PasswordHash: "password",
			Permissions: db.Permissions{
				IsLocked:         true,
				IsEmailConfirmed: true,
				Is2FARequired:    true}},
		PersonalInfo: db.PersonalInfo{
			FirstName: "FirstName",
			LastName:  "LastName",
			BirthDate: db.BirthdayDate(1966, 4, 21)}}

	err := cache.SetUserProfile(profile)
	assert.Nil(t, err)

	pr1, err := cache.GetUserProfile("user-1")
	assert.Nil(t, err)
	assert.NotNil(t, pr1)

	fmt.Printf("DB profile: %v\n", spew.Sdump(profile))
	fmt.Printf("Cache profile: %v\n", spew.Sdump(pr1))

	assert.Equal(t, profile.BirthDate, pr1.GetBirthDate())
}
