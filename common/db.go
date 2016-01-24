package common

import (
  "regexp"
  "strings"
  "database/sql"
  "github.com/cloudfoundry-community/go-cfenv"
)

type DatabaseDSN struct {
  Host string
  Port string
  User string
  Password string
  Name string
}

func IsDatabaseDSNComplete(dsn DatabaseDSN) bool {
  return dsn.Host != "" &&
    dsn.Port != "" &&
    dsn.Name != "" &&
    dsn.User != "" &&
    dsn.Password != ""
}

func BuildDatabaseDSN(credentials map[string]interface{}) DatabaseDSN {
  return DatabaseDSN {
    Host: GetString(credentials, "hostname", "host"),
    Port: GetString(credentials, "port"),
    User: GetString(credentials, "username", "user"),
    Password: GetString(credentials, "password"),
    Name: GetString(credentials, "name"),
  }
}

func ParseDatabaseDSN(urlRegexp *regexp.Regexp, url string) DatabaseDSN {
  defer func() {
    _ = recover()
  }()

  dsn := DatabaseDSN{}

  match := urlRegexp.FindStringSubmatch(url)

  for i, name := range urlRegexp.SubexpNames() {
    if i != 0 {
      if name == "host" {
        dsn.Host = match[i]
      }

      if name == "port" {
        dsn.Port = match[i]
      }

      if name == "user" {
        dsn.User = match[i]
      }

      if name == "password" {
        dsn.Password = match[i]
      }

      if name == "name" {
        dsn.Name = match[i]
      }
    }
  }

  return dsn
}

func isDbTypeUrl(dbType string, url interface{}) bool {
  if url == nil || dbType == "" {
    return false
  }

  lcstr := strings.ToLower(url.(string))
  return strings.Contains(lcstr, dbType + "://")
}

func AcceptDatabaseType(dbType string, service cfenv.Service) bool {
  for _, tag := range service.Tags {
    lctag := strings.ToLower(tag)
    if lctag == dbType {
      return true
    }
  }

  if isDbTypeUrl(dbType, service.Credentials["jdbcUrl"]) {
    return true
  }

  if isDbTypeUrl(dbType, service.Credentials["url"]) {
    return true
  }

  if isDbTypeUrl(dbType, service.Credentials["uri"]) {
    return true
  }

  return false
}

func TestDatabaseService(db *sql.DB, testSql string) (bool, error) {
  _, err := db.Query(testSql)
  return err == nil, err
}
