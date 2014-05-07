package main

import (
	"flag"
	"net/http"

	"github.com/golang/glog"
)

type config struct {
	IndexPath        string
	HTTPAddr         string
	HTTPDocumentRoot string
}

var Config config

func main() {
	httpAddr := flag.String("http", "127.0.0.1:5000", "address and port to listen on")
	httpDocroot := flag.String("root", "www", "HTTP document root for static web files")
	dataPath := flag.String("data", "/usr/local/var/lib/shadowcaster", "data directory (for indexes and such)")
	flag.Parse()

	Config = config{
		IndexPath:        *dataPath,
		HTTPAddr:         *httpAddr,
		HTTPDocumentRoot: *httpDocroot}

	// Run consistency checks on the indexes.
	glog.Infoln("Running consistency checks on the indexes")
	if err := CheckIndexes(*dataPath); err != nil {
		glog.Fatalln(err)
	}
	glog.Infoln("Consistency checks passed")

	// Set up the HTTP handling.
	http.HandleFunc("/movies/", HandleMovies)
	http.HandleFunc("/movies/setdir", HandleSetMovieDir)
	http.HandleFunc("/movies/status", HandleMovieStatus)
	http.HandleFunc("/tv/", HandleTV)
	http.HandleFunc("/music/", HandleMusic)
	http.HandleFunc("/pictures/", HandlePictures)
	http.HandleFunc("/settings/", HandleSettings)
	http.Handle("/", http.FileServer(http.Dir(*httpDocroot)))
	glog.Infof("Listening on %v", *httpAddr)
	if err := http.ListenAndServe(*httpAddr, nil); err != nil {
		glog.Fatalln(err)
	}
	glog.Infof("ShadowCaster offline")
}
