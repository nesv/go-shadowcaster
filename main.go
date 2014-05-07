package main

import (
	"flag"
	"net/http"

	"github.com/boltdb/bolt"
	"github.com/golang/glog"
)

func main() {
	httpAddr := flag.String("http", "127.0.0.1:5000", "address and port to listen on")
	httpDocroot := flag.String("root", "www", "HTTP document root for static web files")
	dbPath := flag.String("db", "shadowcaster.db", "path to database file")
	flag.Parse()

	// Open up the database.
	glog.V(1).Infof("Opening database %v", *dbPath)
	db, err := bolt.Open(*dbPath, 0664)
	if err != nil {
		glog.Fatalln(err)
	}
	defer db.Close()
	defer glog.Infoln("Closing database")
	glog.V(1).Infoln("Running consistency check on database")
	if err := db.Check(); err != nil {
		glog.Fatalln("Database consistency check failed:", err)
	}
	glog.V(1).Infoln("Database consistency check passed")

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
