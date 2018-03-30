package auth

import (
	"encoding/hex"
	"time"

	"github.com/gavrilaf/spawn/pkg/api/types"
	"github.com/gavrilaf/spawn/pkg/cache/mdl"
	"github.com/gavrilaf/spawn/pkg/cryptx"
	db "github.com/gavrilaf/spawn/pkg/dbx/mdl"
	"github.com/gavrilaf/spawn/pkg/errx"

	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
	"gopkg.in/dgrijalva/jwt-go.v3"
)

// HandleLogin handles sign in proccess
// return AuthToken or error
func (mw *Middleware) HandleLogin(p LoginDTO, ctx LoginContext) (AuthTokenDTO, error) {
	log.Infof("auth.HandleLogin, %v, %v", p, ctx)

	p.FixLocale() // set 'en' locale & language if empty

	// Check client
	client, err := mw.bridge.GetClient(p.ClientID)
	if err != nil {
		log.Errorf("auth.HandleLogin, can't find client with ID = %v: (%v)", p.ClientID, err)
		return AuthTokenDTO{}, err
	}

	// Check signature
	if err = p.CheckSignature(client.Secret); err != nil {
		log.Errorf("auth.HandleLogin, invalid signature for %v, must be %v", p, p.GetSignature(client.Secret))
		return AuthTokenDTO{}, types.ErrInvalidSignature
	}

	// Check user
	user, err := mw.bridge.FindUser(p.Username)
	if err != nil {
		log.Errorf("auth.HandleLogin, find user error: (%v)", err)
		return AuthTokenDTO{}, types.ErrUserUnknown
	}

	log.Infof("Found user: %v", user)

	if !p.CheckPassword(user.PasswordHash) {
		log.Errorf("auth.HandleLogin, invalid password for %v", p.Username)
		return AuthTokenDTO{}, types.ErrUserUnknown
	}

	// Check device
	device, err := mw.bridge.GetDevice(user.ID, p.DeviceID)
	if device == nil {
		_, reason := errx.GetErrorReason(err)
		if reason != errx.ReasonNotFound {
			log.Errorf("auth.HandleLogin, check device error: (%v)", err)
			return AuthTokenDTO{}, err
		}

		log.Infof("Login with new device. User %v, device (%v, %v)", user.ID, p.DeviceID, p.DeviceName)
		newDevice := p.CreateDevice()

		newDevice.ID = p.DeviceID
		newDevice.Name = p.DeviceName
		newDevice.UserID = user.ID
		newDevice.Lang = p.Lang
		newDevice.Locale = p.Locale
		newDevice.IsConfirmed = false

		device, err = mw.bridge.AddDevice(user.ID, newDevice)
		if err != nil {
			log.Errorf("auth.HandleLogin, add device error: (%v, %v), (%v)", user.ID, p.DeviceID, err)
			return AuthTokenDTO{}, err
		}

		log.Infof("Added new device %v for user %v", device, user.ID)
	} else {
		log.Infof("Login with registered device %v for user %v", device, user.ID)
	}

	ctx.DeviceName = p.DeviceName

	return mw.makeLogin(client, user, device, &ctx)
}

// HandleRegister handles register proccess. Function creates new user and makes signing in.
// return AuthToken or error
func (mw *Middleware) HandleRegister(p RegisterDTO, ctx LoginContext) (AuthTokenDTO, error) {
	log.Infof("auth.HandleRegister, %v, %v", p, ctx)

	p.FixLocale() // set 'en' locale & language if empty

	// Check client
	client, err := mw.bridge.GetClient(p.ClientID)
	if err != nil {
		log.Errorf("auth.HandleRegister, can't find client with ID = %v: (%v)", p.ClientID, err)
		return AuthTokenDTO{}, err
	}

	// Check signature
	if p.CheckSignature(client.Secret) != nil {
		log.Errorf("auth.HandleRegister, invalid signature for %v, must be %v", p, p.GetSignature(client.Secret))
		return AuthTokenDTO{}, types.ErrInvalidSignature
	}

	// Check user already registered
	alredyExist, _ := mw.bridge.FindUser(p.Username)
	if alredyExist != nil {
		log.Errorf("auth.HandleRegister, user %v already exists", p.Username)
		return AuthTokenDTO{}, types.ErrUserAlreadyExist
	}

	// Create password hash
	pswHash, err := cryptx.GenerateHashedPassword(p.Password)
	if err != nil {
		log.Errorf("auth.HandleRegister, password generate error %v: (%v)", p.Username, err)
		return AuthTokenDTO{}, err
	}

	err = mw.bridge.RegisterUser(p.Username, pswHash, p.CreateDevice())
	if err != nil {
		log.Errorf("auth.HandleRegister, add to storage error %v: (%v)", p.Username, err)
		return AuthTokenDTO{}, err
	}

	user, err := mw.bridge.FindUser(p.Username)
	if err != nil {
		log.Errorf("auth.HandleRegister, find registered user error %v: (%v)", p.Username, err)
		return AuthTokenDTO{}, err
	}

	log.Infof("Created new user: %v", user)

	device, err := mw.bridge.GetDevice(user.ID, p.DeviceID)
	if err != nil {
		log.Errorf("auth.HandleRegister, find device error %v: (%v)", p.Username, err)
		return AuthTokenDTO{}, err
	}

	ctx.DeviceName = p.DeviceName

	return mw.makeLogin(client, user, device, &ctx)
}

// HandleRefresh handles refresh auth token proccess
// return AuthToken (with empty RefreshToken) or error
func (mw *Middleware) HandleRefresh(p RefreshDTO) (AuthTokenDTO, error) {
	token, _ := mw.parseToken(p.AuthToken)
	claims := token.Claims.(jwt.MapClaims)

	sessionID := claims["session_id"].(string)
	origIat := int64(claims["orig_iat"].(float64))
	nonce := int64(claims["nonce"].(float64))

	log.Infof("auth.HandleRefresh, sesson = %v, iat = %d, nonce = %d", sessionID, origIat, nonce)

	if origIat < time.Now().Add(-mw.maxRefresh).Unix() {
		return AuthTokenDTO{}, types.ErrTokenExpired
	}

	session, err := mw.bridge.GetSession(sessionID)
	if err != nil {
		log.Errorf("auth.HandleRefresh, can't find session: (%v)", err)
		return AuthTokenDTO{}, err
	}

	if p.RefreshToken != session.RefreshToken {
		log.Errorf("auth.HandleRefresh, invalid refresh token")
		return AuthTokenDTO{}, types.ErrTokenInvalid
	}

	if nonce != session.Nonce {
		log.Errorf("auth.HandleRefresh, invalid nonce %d, required %d", nonce, session.Nonce)
		return AuthTokenDTO{}, types.ErrTokenExpired
	}

	// Create the token
	newToken := jwt.New(jwt.GetSigningMethod(SigningAlgorithm))
	newClaims := newToken.Claims.(jwt.MapClaims)

	// Copy claims from original token
	for key := range claims {
		newClaims[key] = claims[key]
	}

	now := time.Now()
	expire := now.Add(mw.timeout)
	claims["exp"] = expire.Unix()

	// Update nonce
	nonce += 1
	claims["nonce"] = nonce
	session.Nonce = nonce

	err = mw.bridge.UpdateSession(*session)
	if err != nil {
		log.Errorf("auth.HandleRefresh, can't update session %v: (%v)", sessionID, err)
		return AuthTokenDTO{}, err
	}

	tokenString, err := token.SignedString(session.ClientSecret)
	if err != nil {
		log.Errorf("auth.HandleRefresh, can't create signature %v: (%v)", sessionID, err)
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
			Scopes:            session.Scope,
		}}, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (mw *Middleware) makeLogin(client *db.Client, user *mdl.AuthUser, device *mdl.AuthDevice, ctx *LoginContext) (AuthTokenDTO, error) {
	refreshToken := generateRefreshToken()

	session := mdl.Session{
		RefreshToken:      refreshToken,
		Nonce:             1,
		ClientID:          client.ID,
		ClientSecret:      client.Secret,
		UserID:            user.ID,
		DeviceID:          device.DeviceID,
		IsDeviceConfirmed: device.IsConfirmed,
		Locale:            device.Locale,
		Lang:              device.Lang,
		Permissions:       user.Permissions,
	}

	sessionID, err := mw.bridge.AddSession(session)
	if err != nil {
		log.Errorf("auth.makeLogin, add to storage error (%v), (%v)", session, err)
		return AuthTokenDTO{}, err
	}

	// Create the token
	token := jwt.New(jwt.GetSigningMethod(SigningAlgorithm))
	claims := token.Claims.(jwt.MapClaims)

	now := time.Now()
	expire := now.Add(mw.timeout)
	claims["session_id"] = sessionID
	claims["aud"] = session.ClientID
	claims["exp"] = expire.Unix()
	claims["orig_iat"] = now.Unix()
	claims["iss"] = Realm
	claims["nonce"] = int64(1)

	// Add custom claims here

	tokenString, err := token.SignedString(session.ClientSecret)
	if err != nil {
		log.Errorf("auth.makeLogin, can't create signature: (%v)", err)
		return AuthTokenDTO{}, err
	}

	session.ID = sessionID
	err = mw.bridge.HandlerLogin(session, *ctx)
	if err != nil {
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
			Scopes:            user.Scope,
		}}, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func generateRefreshToken() string {
	k, _ := cryptx.GenerateSaltedKey(uuid.NewV4().String())
	return hex.EncodeToString(k)
}
