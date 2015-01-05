package main

import (
  "encoding/json"
  "fmt"
  "github.com/codegangsta/negroni"
  "github.com/gorilla/mux"
  "github.com/joho/godotenv"
  "log"
  "net/http"
  "os"
)

type Hacker struct {
  Name    string
  Hobbies []string
}

func main() {
  env := godotenv.Load()
  if env != nil {
    log.Fatal("Error loading .env file")
  }

  db_user := os.Getenv("DATABASE_USER")
  logger := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
  logger.Println("got db_user: ", db_user)

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
  params := mux.Vars(r) // from the request
  my_little_json := Hacker{params["hacker"], []string{"music", "programming"}}

  js, err := json.Marshal(my_little_json)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}
