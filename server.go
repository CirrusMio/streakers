package main

import (
  "encoding/json"
  "fmt"
  "net/http"
  "log"
  "github.com/codegangsta/negroni"
  "github.com/gorilla/mux"
  "github.com/PuerkitoBio/goquery"
)

type Hacker struct {
  Name    string
  Hobbies []string
  Streak  int
  Today   string
}


func main() {
  // classic provides Recovery, Logging, Static default middleware
  n := negroni.Classic()

  router := mux.NewRouter()
  router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(w, "Hello World!")
  })

  // GET /hackers/chase
  router.HandleFunc("/hackers/{hacker}", hacker_handler)

  // router goes last
  n.UseHandler(router)
  n.Run(":3000")
}

// learned from: http://www.alexedwards.net/blog/golang-response-snippets#json
func hacker_handler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r) // from the request
  hacker := vars["hacker"]
  my_little_json := Hacker{hacker, []string{"music", "programming"}, 1, contrib_calendar_scraper(hacker)}

  js, err := json.Marshal(my_little_json)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}

func contrib_calendar_scraper(hacker_name string) string {
  url := "https://api.github.com/users/" + hacker_name + "/calendar"
  doc, err := goquery.NewDocument(url)
  if err != nil {
    log.Fatal(err)
  }
}
