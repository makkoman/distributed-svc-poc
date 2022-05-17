package log

import (
	"io/ioutil"
	stLog "log"
	"net/http"
	"os"
)

var log *stLog.Logger

type fileLog string

func (l fileLog) Write(data []byte) (int, error) {
	f, err := os.OpenFile(string(l), os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.FileMode(0600))
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return f.Write(data)
}

func Run(destination string) {
	log = stLog.New(fileLog(destination), "", stLog.LstdFlags)
}

func RegisterHandlers() {
	http.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		msg, err := ioutil.ReadAll(r.Body)
		if err != nil || len(msg) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		write(string(msg))
	})
}

func write(message string) {
	log.Printf("%v\n", message)
}
