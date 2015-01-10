package main

import (
  _ "database/sql"
  "encoding/json"
  "fmt"
  "github.com/codegangsta/negroni"
  "github.com/gorilla/mux"
  "github.com/jinzhu/gorm"
  "github.com/joho/godotenv"
  _ "github.com/lib/pq"
  "github.com/unrolled/render"
  "log"
  "net/http"
  "os"
  "io/ioutil"
  "time"
)

type Api struct {
  DB gorm.DB
}

type Hacker struct {
  Id    int64
  Name  string `sql:"not null;unique"`
  Today bool
}

type GitHubEvent struct {
  Created_At string `json: "created_at"`
}

func main() {
  InitEnv()

  api := Api{}
  api.InitDB()
  api.InitSchema()

  // classic provides Recovery, Logging, Static default middleware
  n := negroni.Classic()
  // for easy template rendering
  r := render.New(render.Options{
    Directory:  "app/templates",
    Layout:     "layout",
    Extensions: []string{".tmpl", ".html"},
  })

  router := mux.NewRouter()

  router.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
    // Assumes you have a template in ./app/templates called "index.html"
    r.HTML(w, http.StatusOK, "index", nil)
  })

  // GET /hackers/chase
  router.HandleFunc("/hackers/{github_username}", api.HackerHandler).Methods("GET")

  // POST /hackers
  router.HandleFunc("/hackers", api.CreateHackerHandler).Methods("POST")

  // router goes last
  n.UseHandler(router)
  n.Run(":3000")
}

// learned from: http://www.alexedwards.net/blog/golang-response-snippets#json
func (api *Api) HackerHandler(w http.ResponseWriter, r *http.Request) {
  params := mux.Vars(r) // from the request
  // need to figure out how to recall a db record from params/Vars
  // return/display JSON dump of saved object
  // my_little_json := Hacker{1, params["github_username"], today("github_username")}
  hacker := Hacker{}
  if err := api.DB.Where("name = ?", params["github_username"]).First(&hacker).Error; err != nil {
    // This really should redirect to a Not Found route. Or return a not found JSON
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  js, err := json.Marshal(&hacker)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}

func (api *Api) CreateHackerHandler(w http.ResponseWriter, r *http.Request) {
  err := r.ParseForm()
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  // save data
  hacker := Hacker{Name: r.Form.Get("name")}
  if err := api.DB.Where(Hacker{Name: r.Form.Get("name")}).FirstOrCreate(&hacker).Error; err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  // make JSON
  js, err := json.Marshal(&hacker)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }

  // return JSON
  w.Header().Set("Content-Type", "application/json")
  w.Write(js)
}

func today(h string) bool {

  var gh []GitHubEvent
  personal_url := fmt.Sprintf("https://api.github.com/users/%s/events", h)
  body := read_github_events(personal_url)
  err := json.Unmarshal(body, &gh)
  if err != nil {
    log.Fatal(err.Error())
  }

  today_string := time.Now().String()[0:10]
  for _, item := range gh {
    answer, err := regexp.MatchString(today_string, item.Created_At)

    if err != nil {
      log.Fatal(err.Error())
    }

    if answer {
      return true
    }
  }

  return false
}


func read_github_events(url string) []byte {
  response, err := http.Get(url)

  if err != nil {
    log.Fatal(err.Error())
  }

  body, err := ioutil.ReadAll(response.Body)

  if err != nil {
    log.Fatal(err.Error())
  }

  return body
}

func InitEnv() {
  env := godotenv.Load()
  if env != nil {
    log.Fatal("Error loading .env file")
  }
}

func (api *Api) InitDB() {
  var err error
  db_user := os.Getenv("DATABASE_USER")
  db_pass := os.Getenv("DATABASE_PASS")
  db_host := os.Getenv("DATABASE_HOST")
  db_name := os.Getenv("DATABASE_NAME")
  logger := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
  logger.Println("got db_user: ", db_user)

  api.DB, err = gorm.Open("postgres", fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable", db_user, db_name, db_pass, db_host))
  if err != nil {
    log.Fatal("database connection error: ", err)
  }
  api.DB.LogMode(true)
}

func (api *Api) InitSchema() {
  api.DB.AutoMigrate(&Hacker{})
}
