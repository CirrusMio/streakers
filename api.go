// Package streakers is an application that tracks a user's commit streak and
// provides encouragement to contribute to open source projects.
package streakers

import (
	"encoding/json"
	"fmt"
	"net/http"
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
	db     *bolt.DB
	names  chan string
	ticker *time.Ticker
	closec chan bool
}

// event is the JSON structure that is returned from Github's events API.
type event struct {
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"createdAt"`
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
		w.WriteHeader(400)
		return
	}

	switch r.Method {
	case "POST":
		createHacker(w, r)
	case "GET":
		showHacker(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		w.WriteHeader(405)
	}
}

// Close shuts down the poll goroutines and stops the ticker.
func (api *API) Close() error {
	close(api.closec)
	api.ticker.Stop()

	return nil
}

// createHacker schedules a hacker update based on their Github handle.
func (api *API) createHacker(w http.ResponseWriter, r *http.Request) {
	api.names <- name

	w.WriteHeader(202)
}

// showHacker retrieves a hacker based on their Github handle and returns
// their last commit time.
func (api *API) showHacker(w http.ResponseWriter, r *http.Request) {
	w.Header.Set("Content-Type", "application/json")

	api.db.View(func(tx *bolt.Tx) error {
		body := tx.Bucket("hackers").Get([]byte(n))
		if body == nil {
			http.NotFound(w, r)

			return fmt.Errorf("no information available for %s", n)
		}

		// Schedule an update.
		api.names <- n

		return w.Write(body)
	})
}

// poll checks for updates for users that it should retrieve from Github.
// Additionally, it will reschedule all known users every PollDuration.
// Returns when closec is closed.
func (api *API) poll() {
	for {
		select {
		case n := <-api.names:
			hacker, err := updatedHacker(n)
			if err != nil {
				continue
			}

			b, err := json.Marshal(hacker)
			if err != nil {
				continue
			}

			api.db.Update(func(tx *bolt.Tx) error {
				return tx.Bucket("hackers").Put([]byte(name), b)
			})
		case <-api.ticker.C:
			h.db.View(func(tx *bolt.Tx) error {
				c := tx.Bucket("hackers").Cursor()
				for k, _ := c.First(); k != nil; k, _ = c.Next() {
					api.names <- k
				}
			})
		case <-api.closec:
			return
		}
	}
}

func updatedHacker(name string) (*hacker, error) {
	res, err := http.Get(fmt.Sprintf(eventsURL, name))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("received status code %d", res.StatusCode)
	}

	var events []eventsJSON
	if err := json.NewDecoder(res.Body).Decode(&events); err != nil {
		return nil, err
	}

	for _, e := range events {
		if e.Type != "PushEvent" {
			continue
		}

		// Not entirely true since this is the timestamp for when the push event
		// occurred rather than when the commit was actually made, but it's close
		// enough for now.
		return &Hacker{
			Name:       name,
			LastCommit: e.CreatedAt,
		}
	}

	return nil, fmt.Errorf("%s has never committed code", name)
}
