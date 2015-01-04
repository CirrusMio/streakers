package main

import (
  "encoding/json"
  "fmt"
  "github.com/codegangsta/negroni"
  "github.com/gorilla/mux"
  "net/http"
)

type Hacker struct {
  Name    string
  Hobbies []string
}

func main() {
  // classic provides Recovery, Logging, Static default middleware
  n := negroni.Classic()

  router := mux.NewRouter()
  router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(w, "Hello World!")
  })

  router.HandleFunc("/json/{hacker}", hacker_handler)

  // router goes last
  n.UseHandler(router)
  n.Run(":3000")
}

// learned from: http://www.alexedwards.net/blog/golang-response-snippets#json
func hacker_handler(w http.ResponseWriter, r *http.Request) {
  vars := mux.Vars(r) // from the request
  hacker := vars["hacker"]
  my_little_json := Hacker{hacker, []string{"music", "programming"}}

  js, err := json.Marshal(my_little_json)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}
