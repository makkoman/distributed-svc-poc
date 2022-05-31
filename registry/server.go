package registry

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

// TODO: get from config
const ServerPort = ":3000"
const ServicesURL = "http://localhost" + ServerPort + "/services"

type registry struct {
	registrations []Registration
	mutex         *sync.Mutex
}

func (r *registry) add(reg Registration) error {
	r.mutex.Lock()
	r.registrations = append(r.registrations, reg)
	// not deferring because we'll need to have control of the mutex
	r.mutex.Unlock()
	return nil
}

func (r *registry) remove(url string) error {
	for i := range r.registrations {
		if r.registrations[i].ServiceURL == url {
			r.mutex.Lock()
			r.registrations = append(r.registrations[:i], r.registrations[i+1:]...)
			r.mutex.Unlock()
			return nil
		}
	}
	return fmt.Errorf("service URL %s not found", url)
}

var reg = registry{
	registrations: make([]Registration, 0),
	mutex:         new(sync.Mutex),
}

type RegistryService struct{}

func (s RegistryService) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	log.Println("request received")
	switch req.Method {
	case http.MethodPost:
		dec := json.NewDecoder(req.Body)
		var r Registration
		err := dec.Decode(&r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		log.Printf("adding service: %v with URL: %v\n", r.ServiceName, r.ServiceURL)
		err = reg.add(r)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	case http.MethodDelete:
		payload, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		url := string(payload)
		log.Printf("removing service with URL: %v\n", url)
		err = reg.remove(url)
		if err != nil {
			log.Printf("error removing service:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
