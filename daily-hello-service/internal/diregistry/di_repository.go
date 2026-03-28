package diregistry

import (
	"daily-hello-service/internal/repositories"
	"go-libs/dihelper"
	"go-libs/sqlormhelper"

	"github.com/sarulabs/di"
)

func initRepositoriesBuilder() {
	dihelper.RepositoriesBuilder = func() []di.Def {
		arr := []di.Def{}
		arr = append(arr,
			di.Def{
				Name:  TokenRepositoryDIName,
				Scope: di.App,
				Build: func(ctn di.Container) (any, error) {
					sql := ctn.Get(DataBaseDIName).(sqlormhelper.SqlGormDatabase)
					db, err := sql.GetConn()
					if err != nil {
						return nil, err
					}
					return repositories.NewTokenRepository(db), nil
				},
				Close: func(obj any) error {
					return nil
				},
			},
			di.Def{
				Name:  UserRepositoryDIName,
				Scope: di.App,
				Build: func(ctn di.Container) (any, error) {
					sql := ctn.Get(DataBaseDIName).(sqlormhelper.SqlGormDatabase)
					db, err := sql.GetConn()
					if err != nil {
						return nil, err
					}
					return repositories.NewUserRepository(db), nil
				},
				Close: func(obj any) error {
					return nil
				},
			},
			di.Def{
				Name:  BranchRepositoryDIName,
				Scope: di.App,
				Build: func(ctn di.Container) (any, error) {
					sql := ctn.Get(DataBaseDIName).(sqlormhelper.SqlGormDatabase)
					db, err := sql.GetConn()
					if err != nil {
						return nil, err
					}
					return repositories.NewBranchRepository(db), nil
				},
				Close: func(obj any) error {
					return nil
				},
			},
			di.Def{
				Name:  AttendanceRepositoryDIName,
				Scope: di.App,
				Build: func(ctn di.Container) (any, error) {
					sql := ctn.Get(DataBaseDIName).(sqlormhelper.SqlGormDatabase)
					db, err := sql.GetConn()
					if err != nil {
						return nil, err
					}
					return repositories.NewAttendanceRepository(db), nil
				},
				Close: func(obj any) error {
					return nil
				},
			},
			di.Def{
				Name:  BranchWifiRepositoryDIName,
				Scope: di.App,
				Build: func(ctn di.Container) (any, error) {
					sql := ctn.Get(DataBaseDIName).(sqlormhelper.SqlGormDatabase)
					db, err := sql.GetConn()
					if err != nil {
						return nil, err
					}
					return repositories.NewBranchWifiRepository(db), nil
				},
				Close: func(obj any) error {
					return nil
				},
			},
			di.Def{
				Name:  ShiftRepositoryDIName,
				Scope: di.App,
				Build: func(ctn di.Container) (any, error) {
					sql := ctn.Get(DataBaseDIName).(sqlormhelper.SqlGormDatabase)
					db, err := sql.GetConn()
					if err != nil {
						return nil, err
					}
					return repositories.NewShiftRepository(db), nil
				},
				Close: func(obj any) error {
					return nil
				},
			},
			di.Def{
				Name:  AttendanceSummaryRepositoryDIName,
				Scope: di.App,
				Build: func(ctn di.Container) (any, error) {
					sql := ctn.Get(DataBaseDIName).(sqlormhelper.SqlGormDatabase)
					db, err := sql.GetConn()
					if err != nil {
						return nil, err
					}
					return repositories.NewAttendanceSummaryRepository(db), nil
				},
				Close: func(obj any) error {
					return nil
				},
			},
			di.Def{
				Name:  DeviceRepositoryDIName,
				Scope: di.App,
				Build: func(ctn di.Container) (any, error) {
					sql := ctn.Get(DataBaseDIName).(sqlormhelper.SqlGormDatabase)
					db, err := sql.GetConn()
					if err != nil {
						return nil, err
					}
					return repositories.NewDeviceRepository(db), nil
				},
				Close: func(obj any) error {
					return nil
				},
			},
		)

		return arr
	}
}
