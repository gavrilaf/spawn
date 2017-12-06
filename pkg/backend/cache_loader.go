package backend

import (
	"time"

	log "github.com/sirupsen/logrus"
)

func (srv *Server) loadAuthCache() {
	log.Infof("Loading users auth cache")
	profiles, errs := srv.db.ReadAllUserProfiles()
	counter := 0

readLoop:
	for {
		select {
		case p := <-profiles:
			if p == nil {
				break readLoop
			}
			devices, err := srv.db.GetUserDevices(p.ID)
			if err != nil {
				log.Errorf("Load devices for user %v error: %v", p.ID, err)
				continue readLoop
			}

			srv.cache.SetUserAuthInfo(*p, devices)
			counter++
			continue readLoop

		case err := <-errs:
			log.Errorf("Read profiles error: %v", err)
			srv.updateServerState(StateError)
			break readLoop

		case <-time.After(1 * time.Second):
			log.Errorf("Read profiles timeout")
			srv.updateServerState(StateError)
			break readLoop
		}
	}

	log.Infof("Auth cache init ok, loaded %d users", counter)
	srv.wg.Done()
}

func (srv *Server) updateCachedUserProfile(id string) error {
	profile, err := srv.db.GetUserProfile(id)
	if err != nil {
		log.Errorf("Could not read profile from db: %v", err)
		return err
	}

	err = srv.cache.SetUserProfile(*profile)
	if err != nil {
		log.Errorf("Could not read profile in cache: %v", err)
		return err
	}

	// Notify about update

	return nil
}
