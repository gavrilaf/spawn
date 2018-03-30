package config

import "github.com/gavrilaf/spawn/pkg/api/types"

var ApiDefaultAccess = types.Access{NeedDevice: true, NeedEmail: false, MinScope: 0}

var ApiAccessConfig = []types.EndpointAccess{
	types.EndpointAccess{Group: gUser, Endpoint: eUserLogout, Access: types.Access{NeedDevice: false, NeedEmail: false, MinScope: 0}},
}
