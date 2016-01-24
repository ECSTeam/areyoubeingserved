package common

import (
  "github.com/cloudfoundry-community/go-cfenv"
)

type ServiceTester interface {
  AcceptService(cfenv.Service) bool
  TestService(cfenv.Service) (bool, error)
}

func GetString(credentials map[string]interface{}, keys ...string) string {
  for _, key := range keys {
    if val, ok := credentials[key]; ok {
      return val.(string)
    }
  }

  return ""
}
