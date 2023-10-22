package database

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Brodobrey-inc/TestService/config"
	"github.com/Brodobrey-inc/TestService/logging"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var (
	DB *sqlx.DB
)

func Initialize() {
	var err error
	DB, err = sqlx.Connect("postgres", toDSN())

	if err != nil {
		logging.LogFatalError(err, "Can't connect to database")
	}
	logging.LogDebug("Successfuly connected to database")

	driver, err := postgres.WithInstance(DB.DB, &postgres.Config{})
	if err != nil {
		logging.LogFatalError(err, "Failed to get postgres driver")
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations",
		"postgres", driver)
	if err != nil {
		logging.LogFatalError(err, "Failed to start migration")
	}

	logging.LogDebug("Start migration for database")
	err = m.Up()

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logging.LogError(err, "Migration did not complete successfully")
	} else {
		logging.LogDebug("Migration completed with no errors")
	}
}

func toDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.ServiceConfig.PostgresUser,
		config.ServiceConfig.PostgresPassword,
		config.ServiceConfig.PostgresHost,
		config.ServiceConfig.PostgresPort,
		config.ServiceConfig.PostgresDBName,
	)
}
