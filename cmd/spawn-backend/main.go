package main

import (
	log "github.com/sirupsen/logrus"
)

func main() {

	log.Info("Spawn backend started")

	log.Info("Preparing cache")
	cache, err := BuildCache()
	if err != nil {
		log.Errorf("Build cache error: %v", err)
	}
	log.Info("Cache is ready")

	cache.PrintState()

}
