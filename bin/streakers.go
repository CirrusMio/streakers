// streakers is an application that tracks a user's commit streak and
// provides encouragement to contribute to open source projects.
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/CirrusMio/streakers"
	"github.com/boltdb/bolt"
)

var (
	dbPath       = flag.String("db", "streakers.db", "path to DB file")
	staticPath   = flag.String("static", "static", "path to static files")
	addr         = flag.String("addr", ":8080", "host:port to run on")
	clientID     = flag.String("id", "", "Github client ID")
	clientSecret = flag.String("secret", "", "Github client secret")
)

func main() {
	flag.Parse()

	db, err := bolt.Open(*dbPath, 0666, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("hackers"))
		return err
	}); err != nil {
		log.Fatal(err)
	}

	api := streakers.NewAPI(db)
	defer api.Close()

	api.SetCredentials(*clientID, *clientSecret)

	http.Handle("/hackers", api)
	http.Handle("/", http.FileServer(http.Dir(*staticPath)))

	log.Fatal(http.ListenAndServe(*addr, nil))
}
