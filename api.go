// Package streakers is an application that tracks a user's commit streak and
// provides encouragement to contribute to open source projects.
package streakers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/boltdb/bolt"
)

var (
	// NumPollers configures the number of goroutines to spawn to poll Github
	// for new events.
	NumPollers = 3

	// PollInterval configures how often hackers are requeued to check for updates.
	PollInterval = 2 * time.Hour
)

// eventsURL is the URL to pull updates for the user from.
const eventsURL = "https://api.github.com/users/%s/events"

// Hacker is a person who writes the open source code.
type Hacker struct {
	Name       string    `json:"name"`
	LastCommit time.Time `json:"lastCommit",omitempty`
}

// API provides an http.Handlers for creating and retrieving information
// about hackers.
type API struct {
	db          *bolt.DB
	rateLimited int64

	names  chan string
	ticker *time.Ticker
	closec chan bool
}

// event is the JSON structure that is returned from Github's events API.
type event struct {
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"created_at"`
}

// NewAPI sets up a new Streakers API.
func NewAPI(db *bolt.DB) *API {
	api := &API{
		db:     db,
		names:  make(chan string, 20),
		ticker: time.NewTicker(PollInterval),
		closec: make(chan bool),
	}

	for i := 0; i < NumPollers; i++ {
		go api.poll()
	}

	return api
}

// ServeHTTP handles the /hackers API endpoint.
func (api *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if n := r.FormValue("name"); n == "" {
		log.Println("received request missing `name` parameter")

		w.WriteHeader(400)
		return
	}

	switch r.Method {
	case "POST":
		api.createHacker(w, r)
	case "GET":
		api.showHacker(w, r)
	default:
		log.Printf("received request with unsupported method %s\n", r.Method)

		w.Header().Set("Allow", "GET, POST")
		w.WriteHeader(405)
	}
}

// Close shuts down the poll goroutines and stops the ticker.
func (api *API) Close() error {
	log.Println("closing API")

	close(api.closec)
	api.ticker.Stop()

	return nil
}

// createHacker schedules a hacker update based on their Github handle.
func (api *API) createHacker(w http.ResponseWriter, r *http.Request) {
	n := r.FormValue("name")
	api.names <- n

	log.Printf("scheduled update for %s\n", n)

	w.WriteHeader(202)
}

// showHacker retrieves a hacker based on their Github handle and returns
// their last commit time.
func (api *API) showHacker(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	n := r.FormValue("name")

	if err := api.db.View(func(tx *bolt.Tx) error {
		body := tx.Bucket([]byte("hackers")).Get([]byte(n))
		if body == nil {
			http.NotFound(w, r)

			return fmt.Errorf("no information available for %s", n)
		}

		select {
		case api.names <- n:
			// Schedule an early update if we can.
			log.Printf("scheduled update for %s\n", n)
		default:
			// If the queue is full drop it on the floor and let the pollers handle it.
		}

		_, err := w.Write(body)
		return err
	}); err != nil {
		log.Printf("failed lookup for %s: %v\n", n, err)
	}
}

// poll checks for updates for users that it should retrieve from Github.
// Additionally, it will reschedule all known users every PollDuration.
// Returns when closec is closed.
func (api *API) poll() {
	log.Println("starting poller")

	for {
		// If we are rate limited sleep the goroutine until we are past the reset
		// point.
		if s := api.rateLimited - time.Now().Unix(); s > 0 {
			time.Sleep(time.Duration(s) * time.Second)
		}

		select {
		case n := <-api.names:
			log.Printf("received update request for %s\n", n)

			hacker, err := api.updatedHacker(n)
			if err != nil {
				log.Printf("failed to update %s: %v\n", n, err)
				continue
			}

			b, err := json.Marshal(hacker)
			if err != nil {
				log.Printf("failed to marshal JSON for %s: %v\n", n, err)
				continue
			}

			if err := api.db.Update(func(tx *bolt.Tx) error {
				return tx.Bucket([]byte("hackers")).Put([]byte(n), b)
			}); err != nil {
				log.Printf("failed to store %s: %v\n", n, err)
			}
		case <-api.ticker.C:
			log.Printf("beginning full update")

			api.db.View(func(tx *bolt.Tx) error {
				c := tx.Bucket([]byte("hackers")).Cursor()
				for k, _ := c.First(); k != nil; k, _ = c.Next() {
					log.Printf("sending update request for %s\n", k)
					api.names <- string(k)
				}
				return nil
			})
		case <-api.closec:
			log.Println("closing poller")
			return
		}
	}
}

func (api *API) updatedHacker(name string) (*Hacker, error) {
	url := fmt.Sprintf(eventsURL, name)

	log.Printf("making request to %s\n", url)

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	switch res.StatusCode {
	case 200:
		if res.Header.Get("X-RateLimit-Remaining") == "0" {
			ts, err := strconv.ParseInt(res.Header.Get("X-RateLimit-Reset"), 10, 64)
			if err != nil {
				return nil, err
			}

			atomic.StoreInt64(&api.rateLimited, ts)
		}
	case 403:
		return nil, fmt.Errorf("exceeded rate limit; dropping update for %s\n", name)
	default:
		return nil, fmt.Errorf("received status code %d", res.StatusCode)
	}

	var events []event
	if err := json.NewDecoder(res.Body).Decode(&events); err != nil {
		return nil, err
	}

	for _, e := range events {
		switch e.Type {
		// TODO: Decide what events should count as activity.
		case "PushEvent", "CreateEvent", "ForkEvent":
		default:
			continue
		}

		return &Hacker{
			Name:       name,
			LastCommit: e.CreatedAt,
		}, nil
	}

	return nil, fmt.Errorf("%s has never committed code", name)
}
