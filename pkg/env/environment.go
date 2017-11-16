package env

type Environment struct{}

type CacheOptions struct{}

func GetEnvironment(path string) *Environment {
	return &Environment{}
}

func (e *Environment) Cache() *CacheOptions {
	return nil
}

func (e *Environment) GetName() string {
	return "Test"
}
