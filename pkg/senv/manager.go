package senv

func GetEnvironment(path string) *Environment {
	switch path {
	case "Test":
		return getTestEnv()
	default:
		return nil
	}
}
