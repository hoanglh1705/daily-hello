package diregistry

import (
	"daily-hello-service/config"
	"go-libs/dihelper"

	"github.com/sarulabs/di"
)

// DI Path
const (
	// Redis
	CacheHelperDIName       string = "RedisCacheHelper"
	RedisClientHelperDIName string = "RedisClientHelper"

	// * Config
	ConfigDIName string = "Config"

	// * Helper
	ModelConverterDIName   string = "ModelConverter"
	AdapterConverterDIName string = "AdapterConverter"

	DataBaseDIName string = "Database"
	GormDBDIName   string = "GormDB"

	// * Repository
	BaseRepositoryDIName              string = "BaseRepository"
	AuditLogRepositoryDIName          string = "AuditLogRepository"
	UserRepositoryDIName              string = "UserRepository"
	BranchRepositoryDIName            string = "BranchRepository"
	AttendanceRepositoryDIName        string = "AttendanceRepository"
	BranchWifiRepositoryDIName        string = "BranchWifiRepository"
	ShiftRepositoryDIName             string = "ShiftRepository"
	AttendanceSummaryRepositoryDIName string = "AttendanceSummaryRepository"
	DeviceRepositoryDIName            string = "DeviceRepository"

	// * Adapter
	JWTMiddlewareDIName string = "JWTMiddleware"

	// * Usecase
	HelloWorldUsecaseDIName string = "HelloWorldUsecase"

	// * Api

	// Admin
	AdminCompanyAPIDIName string = "AdminCompanyAPI"

	// Agent
	AgentCompanyAPIDIName string = "AgentCompanyAPI"

	// integration
	IntegrationStaffAPIDIName string = "IntegrationStaffAPI"
)

func BuildDIContainer() {
	initBuilder()
	dihelper.BuildLibDIContainer()
}

func GetDependency(name string) any {
	return dihelper.GetLibDependency(name)
}

func initBuilder() {
	initConfigBuilder()
	initHelpersBuilder()
	initRepositoriesBuilder()
	initAdaptersBuilder()
	initUseCasesBuilder()
	initApiBuilder()
	dihelper.ConfigsBuilder = func() []di.Def {
		arr := []di.Def{}
		arr = append(arr, di.Def{
			Name:  ConfigDIName,
			Scope: di.App,
			Build: func(ctn di.Container) (any, error) {
				cfg, err := config.Load()
				return cfg, err
			},
			Close: func(obj any) error {
				return nil
			},
		})

		return arr
	}
}
