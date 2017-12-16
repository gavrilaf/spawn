package backend

import (
	"time"

	log "github.com/sirupsen/logrus"
)

func (srv *Server) loadAuthCache() {
	log.Infof("Initializing auth read model...")

	log.Info("Loading clients")

	clients, err := srv.db.GetClients()
	if err != nil {
		log.Fatalf("Could not loading clients, %v", err)
	}

	for _, cl := range clients {
		err = srv.cache.AddClient(cl)
		if err != nil {
			log.Errorf("Could not add client %v, error: %v", cl, err)
		} else {
			log.Infof("Client %v added to the read model", cl.ID)
		}
	}

	log.Info("Loading users")

	counter := 0
	profiles, errs := srv.db.ReadAllUserProfiles()
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

	log.Infof("Read model init ok, loaded %d users", counter)
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
		log.Errorf("Could not save profile in cache: %v", err)
		return err
	}

	return nil
}

func (srv *Server) updateCachedUserDevices(userID string) error {
	devices, err := srv.db.GetUserDevicesEx(userID)
	if err != nil {
		log.Errorf("Could not read user devices from db: %v", err)
		return err
	}

	err = srv.cache.SetUserDevicesInfo(userID, devices)
	if err != nil {
		log.Errorf("Could not save devices into the cache: %v", err)
		return err
	}

	return nil
}
