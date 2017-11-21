package auth

import (
	mdl "github.com/gavrilaf/spawn/pkg/cache/model"
	"github.com/gavrilaf/spawn/pkg/cryptx"
	db "github.com/gavrilaf/spawn/pkg/dbx/model"

	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"gopkg.in/dgrijalva/jwt-go.v3"

	"encoding/hex"
	"time"
)

func (mw *Middleware) HandleLogin(p LoginDTO) (AuthTokenDTO, error) {
	mw.Log.Infof("auth.HandleLogin, %v", p)

	p.FixLocale()

	// Check client
	client, err := mw.Stg.FindClient(p.ClientID)
	if err != nil {
		mw.Log.Errorf("auth.HandleLogin, can't find client with ID = %v: (%v)", p.ClientID, err)
		return AuthTokenDTO{}, err
	}

	// Check signature
	if err = p.CheckSignature(client.Secret); err != nil {
		mw.Log.Errorf("auth.HandleLogin, invalid signature for %v, must be %v", p, p.GetSignature(client.Secret))
		return AuthTokenDTO{}, errInvalidSignature
	}

	// Check user
	user, err := mw.Stg.FindUser(p.Username)
	if err != nil {
		mw.Log.Errorf("auth.HandleLogin, find user error: (%v)", err)
		return AuthTokenDTO{}, err
	}

	if user == nil {
		mw.Log.Errorf("auth.HandleLogin, unknown user: (%v)", p.Username)
		return AuthTokenDTO{}, errUserUnknown
	}

	if !p.CheckPassword(user.PasswordHash) {
		mw.Log.Errorf("auth.HandleLogin, invalid password for %v", p.Username)
		return AuthTokenDTO{}, errUserUnknown
	}

	// Check device
	device, err := mw.Stg.FindDevice(user.ID, p.DeviceID)
	if err != nil {
		mw.Log.Errorf("auth.HandleLogin, check device error: (%v)", err)
		return AuthTokenDTO{}, err
	}

	if device == nil {
		newDevice := p.CreateDevice()

		newDevice.UserID = user.ID
		newDevice.IsConfirmed = false

		device, err = mw.Stg.AddDevice(user.ID, newDevice)
		if err != nil {
			mw.Log.Errorf("auth.HandleLogin, add device: (%v, %v), (%v)", user.ID, p.DeviceID, err)
			return AuthTokenDTO{}, err
		}
	} else {
		// Update device locale if needed
	}

	return mw.makeLogin(client, *user, *device)
}

func (mw *Middleware) HandleRegister(p RegisterDTO) (AuthTokenDTO, error) {
	mw.Log.Infof("auth.HandleRegister, %v", p)

	p.FixLocale()

	// Check client
	client, err := mw.Stg.FindClient(p.ClientID)
	if err != nil {
		mw.Log.Errorf("auth.HandleRegister, can't find client with ID = %v: (%v)", p.ClientID, err)
		return AuthTokenDTO{}, err
	}

	// Check signature
	if p.CheckSignature(client.Secret) != nil {
		mw.Log.Errorf("auth.HandleRegister, invalid signature for %v, must be %v", p, p.GetSignature(client.Secret))
		return AuthTokenDTO{}, errInvalidSignature
	}

	// Check user already registered
	alredyExist, _ := mw.Stg.FindUser(p.Username)
	if alredyExist != nil {
		mw.Log.Errorf("auth.HandleRegister, user %v already exists", p.Username)
		return AuthTokenDTO{}, errUserAlreadyExist
	}

	// Create password hash
	pswHash, err := cryptx.GenerateHashedPassword(p.Password)
	if err != nil {
		mw.Log.Errorf("auth.HandleRegister, password generate error %v: (%v)", p.Username, err)
		return AuthTokenDTO{}, err
	}

	err = mw.Stg.RegisterUser(p.Username, pswHash, p.CreateDevice())
	if err != nil {
		mw.Log.Errorf("auth.HandleRegister, add to storage error %v: (%v)", p.Username, err)
		return AuthTokenDTO{}, err
	}

	user, err := mw.Stg.FindUser(p.Username)
	if err != nil {
		mw.Log.Errorf("auth.HandleRegister, find registered user error %v: (%v)", p.Username, err)
		return AuthTokenDTO{}, err
	}

	device, err := mw.Stg.FindDevice(user.ID, p.DeviceID)
	if err != nil {
		mw.Log.Errorf("auth.HandleRegister, find device error %v: (%v)", p.Username, err)
		return AuthTokenDTO{}, err
	}

	return mw.makeLogin(client, *user, *device)
}

func (mw *Middleware) HandleRefresh(p RefreshDTO) (AuthTokenDTO, error) {
	token, _ := mw.parseToken(p.AuthToken)
	claims := token.Claims.(jwt.MapClaims)

	sessionId := claims["session_id"].(string)
	origIat := int64(claims["orig_iat"].(float64))

	mw.Log.Infof("auth.HandleRefresh, sesson = %v, iat = %v", sessionId, origIat)

	if origIat < time.Now().Add(-mw.MaxRefresh).Unix() {
		return AuthTokenDTO{}, errTokenExpired
	}

	session, err := mw.Stg.FindSession(sessionId)
	if err != nil {
		mw.Log.Errorf("auth.HandleRefresh, can't find session: (%v)", err)
		return AuthTokenDTO{}, err
	}

	if p.RefreshToken != session.RefreshToken {
		mw.Log.Errorf("auth.HandleRefresh, invalid refresh token!!!")
		return AuthTokenDTO{}, errTokenInvalid
	}

	// Create the token
	newToken := jwt.New(jwt.GetSigningMethod(SigningAlgorithm))
	newClaims := newToken.Claims.(jwt.MapClaims)

	// Copy claims from original token
	for key := range claims {
		newClaims[key] = claims[key]
	}

	now := time.Now()
	expire := now.Add(mw.Timeout)
	claims["exp"] = expire.Unix()

	tokenString, err := token.SignedString(session.ClientSecret)
	if err != nil {
		mw.Log.Errorf("auth.HandleRefresh, can't create signature %v: (%v)", sessionId, err)
		return AuthTokenDTO{}, err
	}

	return AuthTokenDTO{
		AuthToken:    tokenString,
		RefreshToken: "",
		Expire:       expire,
		Permissions: PermissionsDTO{
			IsDeviceConfirmed: session.IsDeviceConfirmed,
			IsEmailConfirmed:  session.IsEmailConfirmed,
			Is2FARequired:     session.Is2FARequired,
			IsLocked:          session.IsLocked,
			Scopes:            session.Scopes,
		}}, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (mw *Middleware) CheckAccess(userId string, clientId string, c *gin.Context) bool {
	return true
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (mw *Middleware) makeLogin(client db.Client, user mdl.AuthUser, device mdl.AuthDevice) (AuthTokenDTO, error) {
	sessionId := generateSessionID()
	refreshToken := generateRefreshToken(sessionId)

	session := mdl.Session{
		ID:                sessionId,
		RefreshToken:      refreshToken,
		ClientID:          client.ID,
		ClientSecret:      client.Secret,
		UserID:            user.ID,
		DeviceID:          device.DeviceID,
		IsDeviceConfirmed: device.IsConfirmed,
		Locale:            device.Locale,
		Lang:              device.Lang,
		Permissions:       user.Permissions,
	}

	err := mw.Stg.StoreSession(session)
	if err != nil {
		mw.Log.Errorf("auth.makeLogin, add to storage error (%v), (%v)", session, err)
		return AuthTokenDTO{}, err
	}

	// Create the token
	token := jwt.New(jwt.GetSigningMethod(SigningAlgorithm))
	claims := token.Claims.(jwt.MapClaims)

	now := time.Now()
	expire := now.Add(mw.Timeout)
	claims["session_id"] = session.ID
	claims["aud"] = session.ClientID
	claims["exp"] = expire.Unix()
	claims["orig_iat"] = now.Unix()
	claims["iss"] = Realm

	// Add custom claims here

	tokenString, err := token.SignedString(session.ClientSecret)
	if err != nil {
		mw.Log.Errorf("auth.makeLogin, can't create signature: (%v)", err)
		return AuthTokenDTO{}, err
	}

	return AuthTokenDTO{
		AuthToken:    tokenString,
		RefreshToken: refreshToken,
		Expire:       expire,
		Permissions: PermissionsDTO{
			IsDeviceConfirmed: device.IsConfirmed,
			IsEmailConfirmed:  user.IsEmailConfirmed,
			Is2FARequired:     user.Is2FARequired,
			IsLocked:          user.IsLocked,
			Scopes:            user.Scopes,
		}}, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func generateSessionID() string {
	return uuid.NewV4().String()
}

func generateRefreshToken(sessionId string) string {
	k, _ := cryptx.GenerateSaltedKey(sessionId)
	return hex.EncodeToString(k)
}
