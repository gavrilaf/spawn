package senv

import (
	"fmt"
	"time"
)

type BackendOptions struct {
	URL       string
	QueueName string
	Timeout   time.Duration
}

type CacheOptions struct {
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
	back  BackendOptions
	cache CacheOptions
	db    DBOptions
}

////////////////////////////////////////////////////////////////////////////////////

func (env Environment) GetName() string {
	return env.name
}

func (env Environment) GetBackOpts() BackendOptions {
	return env.back
}

func (env Environment) GetCacheOpts() CacheOptions {
	return env.cache
}

func (env Environment) GetDBOpts() DBOptions {
	return env.db
}

func (env Environment) String() string {
	return fmt.Sprintf("Environment {Name=%s\n\tDB{Driver = %s, Datasource = %s}\n\tCache{URL=%s, Idle=%d, IdleTimeout=%d}\n\tBackend{URL=%s, Queue=%s, Timeout=%d}",
		env.name, env.db.Driver, env.db.DataSource, env.cache.URL, env.cache.MaxIdle, env.cache.IdleTimeout, env.back.URL, env.back.QueueName, env.back.Timeout)
}
