package migration

import (
	"daily-hello-service/config"
	"daily-hello-service/internal/diregistry"
	"go-libs/loghelper"
	"go-libs/migratehelper"
	"go-libs/sqlormhelper"

	"github.com/go-gormigrate/gormigrate/v2"
)

func main() {
	diregistry.BuildDIContainer()
	cfg := diregistry.GetDependency(diregistry.ConfigDIName).(*config.Config)
	_, err := loghelper.InitZapLogger(&loghelper.LoggerOptions{
		AppName:       "migrate",
		MaskingFields: cfg.SensitiveFields,
	})
	if err != nil {
		loghelper.Logger.Panic("Can't init zap logger", loghelper.Error(err))
	}

	StartMigrate(cfg)
}

func StartMigrate(cfg *config.Config) error {
	if !cfg.AutoMigration {
		return nil
	}
	sqlDB := diregistry.GetDependency(diregistry.DataBaseDIName).(sqlormhelper.SqlGormDatabase)
	db, err := sqlDB.GetConn()
	if err != nil {
		loghelper.Logger.Panic("error when get db conn", loghelper.Error(err))
	}

	migrationDB := []*gormigrate.Migration{}
	migrationDB = append(migrationDB, migrateTable()...)
	migrationDB = append(migrationDB, migrateData(cfg)...)
	migrationDB = append(migrationDB, migrateView()...)

	migration := migratehelper.NewGormMigration(db, migratehelper.GormMigrateOpts{Migrations: migrationDB})
	if err = migration.Run(db); err != nil {
		loghelper.Logger.Error("error when migrate data", loghelper.Error(err))
	}

	return nil
}
