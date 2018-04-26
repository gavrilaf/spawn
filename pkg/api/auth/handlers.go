package auth

import (
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/dgrijalva/jwt-go.v3"

	"github.com/gavrilaf/spawn/pkg/api/defs"
	"github.com/gavrilaf/spawn/pkg/api/ginx"
	"github.com/gavrilaf/spawn/pkg/cache/mdl"
	"github.com/gavrilaf/spawn/pkg/cryptx"
	db "github.com/gavrilaf/spawn/pkg/dbx/mdl"
	"github.com/gavrilaf/spawn/pkg/errx"
)

const (
	TokenLifeTime   = time.Hour
	TokenMaxRefresh = time.Hour * 24
)

// handleSignIn -
// return AuthToken or error
func (self ApiImpl) handleSignIn(p LoginDTO, ctx LoginContext) (AuthTokenDTO, error) {
	log.Infof("auth.handleSignIn: %s, %s", p.String(), ctx.String())

	p.FixLocale() // set 'en' locale & language if empty

	// Check client
	client, err := self.getClient(p.ClientID)
	if err != nil {
		log.Errorf("auth.handleSignIn, can't find client with ID = %s: (%v)", p.ClientID, err)
		return AuthTokenDTO{}, err
	}

	// Check signature
	if err = p.CheckSignature(client.Secret); err != nil {
		log.Errorf("auth.handleSignIn, invalid signature for %s", p.Username)
		return AuthTokenDTO{}, defs.ErrInvalidSignature
	}

	// Check user
	user, err := self.findUser(p.Username)
	if err != nil {
		log.Errorf("auth.handleSignIn, find user error: (%v)", err)
		return AuthTokenDTO{}, defs.ErrUserUnknown
	}

	log.Infof("Found user: %s, %s", user.ID, user.Username)

	if !p.CheckPassword(user.PasswordHash) {
		log.Errorf("auth.handleSignIn, invalid password for %s", p.Username)
		return AuthTokenDTO{}, defs.ErrUserUnknown
	}

	// Check device
	device, err := self.getDevice(user.ID, p.DeviceID)
	if device == nil {
		_, reason := errx.GetErrorReason(err)
		if reason != errx.ReasonNotFound {
			log.Errorf("auth.handleSignIn, check device error: (%v)", err)
			return AuthTokenDTO{}, err
		}

		log.Infof("Login with new device. User %s, device (%s, %s)", user.ID, p.DeviceID, p.DeviceName)
		newDevice := p.CreateDevice()

		newDevice.DeviceID = p.DeviceID
		newDevice.Name = p.DeviceName
		newDevice.UserID = user.ID
		newDevice.Lang = p.Lang
		newDevice.Locale = p.Locale
		newDevice.IsConfirmed = false

		device, err = self.addDevice(user.ID, newDevice)
		if err != nil {
			log.Errorf("auth.handleSignIn, add device error: (%s, %s), (%v)", user.ID, p.DeviceID, err)
			return AuthTokenDTO{}, err
		}

		log.Infof("Added new device %s for user %s", device.DeviceID, user.ID)
	} else {
		log.Infof("Login with registered device %s for user %s", device.DeviceID, user.ID)

		// If user already signed in ?
		session, _ := self.ReadModel.FindSession(user.ID, device.DeviceID)
		if session != nil {
			log.Infof("User already signed in, invalidate old session, %s, %s, %s", session.ID, user.ID, p.DeviceID)

			// Invalidate old session
			err = self.ReadModel.DeleteSession(session.ID)
			if err != nil {
				log.Errorf("auth.handleSignIn, invalidate old session error: (%s, %s), (%v)", user.ID, p.DeviceID, err)
				return AuthTokenDTO{}, err
			}
		}
	}

	ctx.DeviceName = p.DeviceName

	return self.makeLogin(client, user, device, &ctx)
}

// handleSignUp -
// return AuthToken or error
func (self ApiImpl) handleSignUp(p RegisterDTO, ctx LoginContext) (AuthTokenDTO, error) {
	log.Infof("auth.handleSignUp, %s, %s", p.String(), ctx.String())

	p.FixLocale() // set 'en' locale & language if empty

	// Check client
	client, err := self.getClient(p.ClientID)
	if err != nil {
		log.Errorf("auth.handleSignUp, can't find client with ID = %s: (%v)", p.ClientID, err)
		return AuthTokenDTO{}, err
	}

	// Check signature
	if p.CheckSignature(client.Secret) != nil {
		log.Errorf("auth.handleSignUp, invalid signature for %s", p.Username)
		return AuthTokenDTO{}, defs.ErrInvalidSignature
	}

	// Check user already registered
	alredyExist, _ := self.findUser(p.Username)
	if alredyExist != nil {
		log.Errorf("auth.handleSignUp, user %s already exists", p.Username)
		return AuthTokenDTO{}, defs.ErrUserAlreadyExist
	}

	// Create password hash
	pswHash, err := cryptx.GenerateHashedPassword(p.Password)
	if err != nil {
		log.Errorf("auth.handleSignUp, password generate error %s: (%v)", p.Username, err)
		return AuthTokenDTO{}, err
	}

	err = self.registerUser(p.Username, pswHash, p.CreateDevice())
	if err != nil {
		log.Errorf("auth.handleSignUp, add to storage error %s: (%v)", p.Username, err)
		return AuthTokenDTO{}, err
	}

	user, err := self.findUser(p.Username)
	if err != nil {
		log.Errorf("auth.handleSignUp, find registered user error %s: (%v)", p.Username, err)
		return AuthTokenDTO{}, err
	}

	log.Infof("Created new user: %s, %s", user.ID, user.Username)

	device, err := self.getDevice(user.ID, p.DeviceID)
	if err != nil {
		log.Errorf("auth.handleSignUp, find device error %s: (%v)", p.Username, err)
		return AuthTokenDTO{}, err
	}

	ctx.DeviceName = p.DeviceName

	return self.makeLogin(client, user, device, &ctx)
}

// handleRefresh -
// return AuthToken (with empty RefreshToken) or error
func (self ApiImpl) handleRefresh(p RefreshDTO) (AuthTokenDTO, error) {
	token, _ := ginx.ParseToken(p.AuthToken, func(id string) (interface{}, error) {
		cl, err := self.getClient(id)
		if err != nil {
			return nil, err
		}
		return cl.Secret, nil
	})

	claims := token.Claims.(jwt.MapClaims)
	sessionID := claims["session_id"].(string)
	origIat := int64(claims["orig_iat"].(float64))
	nonce := int64(claims["nonce"].(float64))

	if origIat < time.Now().Add(-TokenMaxRefresh).Unix() {
		return AuthTokenDTO{}, defs.ErrTokenExpired
	}

	session, err := self.getSession(sessionID)
	if err != nil {
		log.Errorf("auth.handleRefresh, can't find session: (%v)", err)
		return AuthTokenDTO{}, err
	}

	if p.RefreshToken != session.RefreshToken {
		log.Errorf("auth.handleRefresh, invalid refresh token")
		return AuthTokenDTO{}, defs.ErrTokenInvalid
	}

	if nonce != session.Nonce {
		log.Errorf("auth.handleRefresh, invalid nonce %d, required %d", nonce, session.Nonce)
		return AuthTokenDTO{}, defs.ErrTokenExpired
	}

	// Create the token
	newToken := jwt.New(jwt.GetSigningMethod(defs.SigningAlgorithm))
	newClaims := newToken.Claims.(jwt.MapClaims)

	// Copy claims from original token
	for key := range claims {
		newClaims[key] = claims[key]
	}

	now := time.Now()
	expire := now.Add(TokenLifeTime)
	claims["exp"] = expire.Unix()

	// Update nonce
	nonce += 1
	claims["nonce"] = nonce
	session.Nonce = nonce

	err = self.updateSession(*session)
	if err != nil {
		log.Errorf("auth.handleRefresh, can't update session %s: (%v)", sessionID, err)
		return AuthTokenDTO{}, err
	}

	tokenString, err := token.SignedString(session.ClientSecret)
	if err != nil {
		log.Errorf("auth.handleRefresh, can't create signature %s: (%v)", sessionID, err)
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

func (self ApiImpl) makeLogin(client *db.Client, user *mdl.AuthUser, device *mdl.AuthDevice, ctx *LoginContext) (AuthTokenDTO, error) {
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

	sessionID, err := self.addSession(session)
	if err != nil {
		log.Errorf("auth.makeLogin, add to storage error: (%v)", err)
		return AuthTokenDTO{}, err
	}

	// Create the token
	token := jwt.New(jwt.GetSigningMethod(defs.SigningAlgorithm))
	claims := token.Claims.(jwt.MapClaims)

	now := time.Now()
	expire := now.Add(TokenLifeTime)
	claims["session_id"] = sessionID
	claims["aud"] = session.ClientID
	claims["exp"] = expire.Unix()
	claims["orig_iat"] = now.Unix()
	claims["iss"] = defs.Realm
	claims["nonce"] = int64(1)

	// Add custom claims here

	tokenString, err := token.SignedString(session.ClientSecret)
	if err != nil {
		log.Errorf("auth.makeLogin, can't create signature: (%v)", err)
		return AuthTokenDTO{}, err
	}

	session.ID = sessionID
	err = self.handleLogin(session, *ctx)
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
