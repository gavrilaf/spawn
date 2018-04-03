package config

import "github.com/gavrilaf/spawn/pkg/api/defs"

var ApiDefaultAccess = defs.Access{Locked: false, Device: true, Email: false, Scope: 0}

var ApiAccessConfig = []defs.EndpointAccess{
	defs.EndpointAccess{Group: gUser, Endpoint: eUserLogout, Access: defs.Access{Locked: true, Device: false, Email: false, Scope: 0}},
	defs.EndpointAccess{Group: gUser, Endpoint: eUserState, Access: defs.Access{Locked: false, Device: false, Email: false, Scope: 0}},
}
