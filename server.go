package main

import (
  "fmt"
  "github.com/codegangsta/negroni"
  "github.com/gorilla/mux"
  "net/http"
)

func main() {
  // classic provides Recovery, Logging, Static default middleware
  n := negroni.Classic()

  router := mux.NewRouter()
  router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(w, "Hello World!")
  })

  // router goes last
  n.UseHandler(router)
  n.Run(":3000")
}
