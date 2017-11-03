package auth

import (
	"github.com/gavrilaf/spawn/pkg/cryptx"

	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"gopkg.in/dgrijalva/jwt-go.v3"

	"encoding/hex"
	"time"
)

func (mw *Middleware) HandleLogin(p *LoginDTO) (*AuthTokenDTO, error) {
	mw.Log.Infof("auth.HandleLogin, %v", p)

	// Check client
	client, err := mw.Storage.FindClientByID(p.ClientID)
	if err != nil {
		mw.Log.Errorf("auth.HandleLogin, can't find client with ID = %v: (%v)", p.ClientID, err)
		return nil, err
	}

	// Check signature
	if err = p.CheckSignature(client.Secret()); err != nil {
		mw.Log.Errorf("auth.HandleLogin, invalid signature for %v", p)
		return nil, errInvalidSignature
	}

	// Check user
	user, err := mw.Storage.FindUserByUsername(p.Username)
	if err != nil {
		mw.Log.Errorf("auth.HandleLogin, can't find user = %v: (%v)", p.Username, err)
		return nil, err
	}

	if !p.CheckPassword(user.PasswordHash) {
		mw.Log.Errorf("auth.HandleLogin, invalid password for %v", p.Username)
		return nil, errUserUnknown
	}

	// Check device
	if !p.CheckDevice(user.Devices) {
		// TODO: Send email about new device
		mw.Log.Errorf("auth.HandleLogin, unknown device for %v: %v", p.Username, p.DeviceID)
		return nil, errDeviceUnknown
	}

	// Generate token

	sessionId := mw.GenerateSessionID()
	refreshToken := mw.GenerateRefreshToken(sessionId)

	session := Session{ID: sessionId, RefreshToken: refreshToken, ClientID: client.ID(), ClientSecret: client.Secret(), UserID: user.ID}

	err = mw.Storage.StoreSession(session)
	if err != nil {
		mw.Log.Errorf("auth.HandleLogin, add to storage error %v: (%v) - (%v)", p.Username, session, err)
		return nil, err
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

	// Add custom claims

	tokenString, err := token.SignedString(session.ClientSecret)
	if err != nil {
		mw.Log.Errorf("auth.HandleLogin, can't create signature %v: (%v)", p.Username, err)
		return nil, err
	}

	return &AuthTokenDTO{AuthToken: tokenString, RefreshToken: refreshToken, Expire: expire}, nil
}

func (mw *Middleware) HandleRefresh(p *RefreshDTO) (*AuthTokenDTO, error) {
	token, _ := mw.parseToken(p.AuthToken)
	claims := token.Claims.(jwt.MapClaims)

	sessionId := claims["session_id"].(string)
	origIat := int64(claims["orig_iat"].(float64))

	mw.Log.Infof("auth.HandleRefresh, sesson = %v, iat = %v", sessionId, origIat)

	if origIat < time.Now().Add(-mw.MaxRefresh).Unix() {
		return nil, errTokenExpired
	}

	session, err := mw.Storage.FindSessionByID(sessionId)
	if err != nil {
		mw.Log.Errorf("auth.HandleRefresh, can't find session: (%v)", err)
		return nil, err
	}

	if p.RefreshToken != session.RefreshToken {
		mw.Log.Errorf("auth.HandleRefresh, invalid refresh token!!!")
		return nil, errTokenInvalid
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
		return nil, err
	}

	return &AuthTokenDTO{AuthToken: tokenString, RefreshToken: "", Expire: expire}, nil
}

func (mw *Middleware) HandleRegister(p *RegisterDTO) (*UserRegisteredDTO, error) {
	mw.Log.Infof("auth.HandleRegister, %v", p)

	// Check client
	client, err := mw.Storage.FindClientByID(p.ClientID)
	if err != nil {
		mw.Log.Errorf("auth.HandleRegister, can't find client with ID = %v: (%v)", p.ClientID, err)
		return nil, err
	}

	// Check signature
	if p.CheckSignature(client.Secret()) != nil {
		mw.Log.Errorf("auth.HandleRegister, invalid signature for %v", p)
		return nil, errInvalidSignature
	}

	// Create password hash

	pswHash, err := cryptx.GenerateHashedPassword(p.Password)
	if err != nil {
		mw.Log.Errorf("auth.HandleRegister, password generate error %v: (%v)", p, err)
		return nil, err
	}

	err = mw.Storage.AddUser(p.ClientID, p.DeviceID, p.Username, pswHash)
	if err != nil {
		mw.Log.Errorf("auth.HandleRegister, add to storage error %v: (%v)", p, err)
		return nil, err
	}

	return &UserRegisteredDTO{}, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func (mw *Middleware) CheckAccess(userId string, clientId string, c *gin.Context) bool {
	return true
}

func (mw *Middleware) GenerateSessionID() string {
	return uuid.NewV4().String()
}

func (mw *Middleware) GenerateRefreshToken(sessionId string) string {
	k, _ := cryptx.GenerateSaltedKey(sessionId)
	return hex.EncodeToString(k)
}
