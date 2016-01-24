package main

import (
  "os"
  "log"
  "github.com/jghiloni/willitserve/common"
  "github.com/jghiloni/willitserve/mysql"
  "github.com/cloudfoundry-community/go-cfenv"
  "net/http"
  "html/template"
)

var templates *template.Template = (*template.Template)(nil)

var testers map[string]common.ServiceTester = map[string]common.ServiceTester{
  "mysql": mysql.ServiceTester{},
}

type Result struct {
  ServiceName string
  TesterName string
  Accepted bool
  Success bool
  Notes string
}

func handle(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/html")
  appEnv, _ := cfenv.Current()

  results := []Result{}

  for _, servicesSlice := range appEnv.Services {
    for _, service := range servicesSlice {
      for key, tester := range testers {
        result := Result{
          ServiceName: service.Name,
          TesterName: key,
        }
        if tester.AcceptService(service) {
          result.Accepted = true
          success, err := tester.TestService(service)

          result.Success = success
      	  if err != nil {
      	    result.Notes = err.Error()
      	  }
        } else {
          result.Accepted = false
        }

        results = append(results, result)
      }
    }
  }

  err := templates.ExecuteTemplate(w, "tester.html", results)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func main() {

  templates = template.Must(template.ParseFiles("templates/tester.html"))

  http.HandleFunc("/", handle)

  err := http.ListenAndServe(":" + os.Getenv("PORT"), nil)
  if err != nil {
    log.Fatal(err)
  }
}
