package proxy

import (
	"path/filepath"
	"log"
	"os"
	"net"
)

func ListenToNewSock(newsock string, sigc chan os.Signal) (l net.Listener, err error) {
	// extract directory for newsock
	dir, _ := filepath.Split(newsock)
	// attempt to create dir and ignore if it's already existing
	_ = os.Mkdir(dir, 0777)
	l, err = net.Listen("unix", newsock)
	if err != nil {
		panic(err)
	}
	os.Chmod(newsock, 0666)
	log.Println("[gk-soxy] Listening on " + newsock)
	go func(c chan os.Signal) {
		sig := <-c
		log.Printf("[gk-soxy] Caught signal %s: shutting down.\n", sig)
		if err := l.Close(); err != nil {
			panic(err)
		}
		os.Exit(0)
	}(sigc)
	return
}
