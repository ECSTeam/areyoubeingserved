package mysql

import (
  "fmt"
  "regexp"
  "github.com/cloudfoundry-community/go-cfenv"
  "database/sql"
  _ "github.com/go-sql-driver/mysql"
  "github.com/jghiloni/willitserve/common"
)

type ServiceTester struct {}


var mysqlRegexp	*regexp.Regexp =
  regexp.MustCompile(`mysql://(?P<user>\w+)\:(?P<password>\w+)\@(?P<host>[^:]+)(?:\:(?P<port>\d+))?/(?P<name>[^\?]+)`)

var mysqlJdbcRegexp *regexp.Regexp =
  regexp.MustCompile(`jdbc:mysql://(?P<host>[^:]+)(?:\:(?P<port>\d+))?/(?P<name>[^\?]+)\?user=(?P<user>\w+)\x26password=(?P<password>\w+)`)

func buildDSNString(credentials map[string]interface{}) string {
  dsn := common.BuildDatabaseDSN(credentials)

  if !common.IsDatabaseDSNComplete(dsn) {
    if jdbcUrl, jdbcOk := credentials["jdbcUrl"]; jdbcOk {
      dsn = common.ParseDatabaseDSN(mysqlJdbcRegexp, jdbcUrl.(string))
      if !common.IsDatabaseDSNComplete(dsn) {
        if uri, uriOk := credentials["uri"]; uriOk {
          dsn = common.ParseDatabaseDSN(mysqlRegexp, uri.(string))
        }
      }
    }
  }

  if common.IsDatabaseDSNComplete(dsn) {
    return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dsn.User, dsn.Password, dsn.Host, dsn.Port, dsn.Name)
  }

  return ""
}

func (m ServiceTester) AcceptService(service cfenv.Service) bool {
  return common.AcceptDatabaseType("mysql", service)
}

func (m ServiceTester) TestService(service cfenv.Service) (bool, error) {
  dsnString := buildDSNString(service.Credentials)

  if dsnString == "" {
    return false, fmt.Errorf("Unable to construct DSN string")
  }

  db, err := sql.Open("mysql", dsnString)
  if err != nil {
    return false, err
  }

  return common.TestDatabaseService(db, "SELECT 'success'")
}
