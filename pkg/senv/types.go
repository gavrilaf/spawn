package senv

import "time"

type RPCOptions struct {
	URL     string
	Timeout time.Duration
}

type RedisOptions struct {
	URL         string
	MaxIdle     int
	IdleTimeout time.Duration
}

type DBOptions struct {
	Driver     string
	DataSource string
}

type Environment struct {
	name  string
	rpc   RPCOptions
	redis RedisOptions
	db    DBOptions
}

////////////////////////////////////////////////////////////////////////////////////

func (env Environment) GetName() string {
	return env.name
}

func (env Environment) GetRPCOpts() RPCOptions {
	return env.rpc
}

func (env Environment) GetRedisOpts() RedisOptions {
	return env.redis
}

func (env Environment) GetDBOpts() DBOptions {
	return env.db
}
