package config

import "github.com/gavrilaf/spawn/pkg/api/defs"

var ApiDefaultAccess = defs.Access{NeedDevice: true, NeedEmail: false, MinScope: 0}

var ApiAccessConfig = []defs.EndpointAccess{
	defs.EndpointAccess{Group: gUser, Endpoint: eUserLogout, Access: defs.Access{NeedDevice: false, NeedEmail: false, MinScope: 0}},
}
