package diregistry

import (
	"daily-hello-service/config"
	"daily-hello-service/internal/repositories"
	"daily-hello-service/internal/services"
	"time"

	"go-libs/dihelper"

	"github.com/sarulabs/di"
)

func initUseCasesBuilder() {
	dihelper.UsecasesBuilder = func() []di.Def {
		arr := []di.Def{}
		arr = append(arr,
			di.Def{
				Name:  AuthServiceDIName,
				Scope: di.App,
				Build: func(ctn di.Container) (any, error) {
					cfg := ctn.Get(ConfigDIName).(*config.Config)
					userRepo := ctn.Get(UserRepositoryDIName).(*repositories.UserRepository)
					tokenRepo := ctn.Get(TokenRepositoryDIName).(repositories.TokenRepository)
					return services.NewAuthService(
						userRepo,
						tokenRepo,
						cfg.JwtConfig.SecretKey,
						cfg.JwtConfig.Duration,
						cfg.JwtConfig.DurationRefresh,
					), nil
				},
				Close: func(obj any) error {
					return nil
				},
			},
			di.Def{
				Name:  UserServiceDIName,
				Scope: di.App,
				Build: func(ctn di.Container) (any, error) {
					userRepo := ctn.Get(UserRepositoryDIName).(*repositories.UserRepository)
					return services.NewUserService(userRepo), nil
				},
				Close: func(obj any) error {
					return nil
				},
			},
			di.Def{
				Name:  BranchServiceDIName,
				Scope: di.App,
				Build: func(ctn di.Container) (any, error) {
					branchRepository := ctn.Get(BranchRepositoryDIName).(repositories.BranchRepository)
					return services.NewBranchService(branchRepository), nil
				},
				Close: func(obj any) error {
					return nil
				},
			},
			di.Def{
				Name:  BranchWifiServiceDIName,
				Scope: di.App,
				Build: func(ctn di.Container) (any, error) {
					branchWifiRepo := ctn.Get(BranchWifiRepositoryDIName).(repositories.BranchWifiRepository)
					branchRepo := ctn.Get(BranchRepositoryDIName).(repositories.BranchRepository)
					return services.NewBranchWifiService(branchWifiRepo, branchRepo), nil
				},
				Close: func(obj any) error {
					return nil
				},
			},
			di.Def{
				Name:  AttendanceServiceDIName,
				Scope: di.App,
				Build: func(ctn di.Container) (any, error) {
					cfg := ctn.Get(ConfigDIName).(*config.Config)
					attendanceRepo := ctn.Get(AttendanceRepositoryDIName).(*repositories.AttendanceRepository)
					branchRepo := ctn.Get(BranchRepositoryDIName).(repositories.BranchRepository)
					branchWifiRepo := ctn.Get(BranchWifiRepositoryDIName).(repositories.BranchWifiRepository)
					locationService := services.NewLocationService(branchRepo, branchWifiRepo)
					location, err := loadAttendanceTimezone(cfg.Database.Timezone)
					if err != nil {
						return nil, err
					}
					return services.NewAttendanceService(attendanceRepo, branchRepo, locationService, location), nil
				},
				Close: func(obj any) error {
					return nil
				},
			},
			di.Def{
				Name:  DeviceServiceDIName,
				Scope: di.App,
				Build: func(ctn di.Container) (any, error) {
					deviceRepo := ctn.Get(DeviceRepositoryDIName).(*repositories.DeviceRepository)
					return services.NewDeviceService(deviceRepo), nil
				},
				Close: func(obj any) error {
					return nil
				},
			},
		)

		return arr
	}
}

func loadAttendanceTimezone(name string) (*time.Location, error) {
	if name == "" {
		name = "Asia/Ho_Chi_Minh"
	}

	location, err := time.LoadLocation(name)
	if err == nil {
		return location, nil
	}

	return time.FixedZone("GMT+7", 7*60*60), nil
}
