package common

import (
  "github.com/cloudfoundry-community/go-cfenv"
  "strconv"
)

type ServiceTester interface {
  AcceptService(cfenv.Service) bool
  TestService(cfenv.Service) (bool, error)
}

func GetString(credentials map[string]interface{}, keys ...string) string {
  for _, key := range keys {
    if val, ok := credentials[key]; ok {
      s, ok := val.(string)
      if !ok {
        switch val.(type) {
        case int:
          s = strconv.Itoa(val.(int))
        default:
          s = ""
        }
      }
      return s
    }
  }

  return ""
}
