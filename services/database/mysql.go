package mysqlDB

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	mysql "github.com/go-sql-driver/mysql"
)

var (
	MysqlDB *sql.DB
	dbErr   error
)

const (
	userName string = "root"
	password string = "davy@N.1"
	dbName   string = "SimpNote"
)

type Migration struct {
	TableName  string
	Definition []Column
}

type Column struct {
	Name  string
	Type  string
	Extra string
	Key   string
}

func InitDBConnection(logger *log.Logger, migrations ...Migration) error {
	dataSourceName := fmt.Sprintf("%s:%s@/%s", userName, password, dbName)

	MysqlDB, dbErr = sql.Open("mysql", dataSourceName)

	if dbErr != nil {
		logger.Fatalf("failed to connect to mysql\n Error is: %s\n", dbErr)
	}

	if pingErr := MysqlDB.Ping(); pingErr != nil {
		logger.Fatalf("failed to ping the mysql database\n Error is: %s\n", pingErr)
	}

	if logErr := mysql.SetLogger(logger); logErr != nil {
		logger.Fatalf("failed to set the mysql logger\n Error is: %s\n", logErr)
	}

	if migrationErr := migrateTables(logger, migrations...); migrationErr != nil {
		logger.Fatalf("failed to migrate the tables\n Error is: %s\n", migrationErr)
	}

	return nil
}

func migrateTables(logger *log.Logger, migrations ...Migration) error {

	tx, err := MysqlDB.Begin()
	if err != nil {
		return err
	}
	for _, m := range migrations {

		query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s(", m.TableName)

		for k, c := range m.Definition {
			definition := fmt.Sprintf("%s %s %s %s", c.Name, c.Type, c.Extra, c.Key)
			query = query + strings.Trim(definition, " ")
			if k != len(m.Definition)-1 {
				query = query + ","
			}
		}

		query = query + ")"

		_, err = tx.Exec(query)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return rollbackErr
			}
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

// CloseDBConnection closes the database connection and releases the db so that other apps can use it
func CloseDBConnection() {
	MysqlDB.Close()
}
