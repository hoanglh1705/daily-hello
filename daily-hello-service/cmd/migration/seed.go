package migration

import (
	config "daily-hello-service/config"

	"github.com/go-gormigrate/gormigrate/v2"
)

func migrateData(_ *config.Config) []*gormigrate.Migration {
	migrations := []*gormigrate.Migration{}
	// migrations = append(migrations,
	// &gormigrate.Migration{
	// 	ID: "202404031338",
	// 	Migrate: func(tx *gorm.DB) error {
	// 		changes := []string{
	// 			`INSERT INTO "roles" ("name", "code",description,created_at,updated_at,deleted_at) VALUES
	// 			('Admin','ADM','Admin','2024-04-03 23:49:59+07',NULL,NULL),
	// 			('Manager','MNG','Trưởng BP','2024-04-03 23:50:27+07',NULL,NULL),
	// 			('Supervisor','SPV','PP/TP','2024-04-03 23:50:28+07',NULL,NULL),
	// 			('Specialist','SPL','CV/CVCC','2024-04-03 23:50:28+07',NULL,NULL);`,
	// 			`COMMIT;`,
	// 		}

	// 		return migratehelper.ExecMultiple(tx, strings.Join(changes, " "))
	// 	},
	// 	Rollback: func(tx *gorm.DB) error {
	// 		changes := []string{
	// 			`ROLLBACK;`,
	// 		}

	// 		return migratehelper.ExecMultiple(tx, strings.Join(changes, " "))
	// 	},
	// },

	// )

	return migrations
}
