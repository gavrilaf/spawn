package defs

const (
	Realm = "Spawn"

	TokenHeadName    = "Bearer"
	SigningAlgorithm = "HS256"
	TokenLookup      = "Authorization"

	AuthTypeSimple = "simple"

	EndpointKey = "EnpointKey"
	SessionKey  = "Session"
)

var EmptySuccessResponse = map[string]interface{}{"success": true}

type Endpoint struct {
	Path   string
	Method string
}

type Access struct {
	Locked bool
	Device bool
	Email  bool
	Scope  int
}

type EndpointAccess struct {
	Group string
	Endpoint
	Access
}

func GetEndpointKey(group string, endpoint Endpoint) string {
	return group + endpoint.Path + ":" + endpoint.Method
}
