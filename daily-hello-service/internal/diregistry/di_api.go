package diregistry

import (
	"daily-hello-service/internal/handlers"
	"daily-hello-service/internal/services"
	"go-libs/dihelper"

	"github.com/sarulabs/di"
)

func initApiBuilder() {
	dihelper.APIsBuilder = func() []di.Def {
		arr := []di.Def{}
		arr = append(arr,
			di.Def{
				Name:  AuthAPIDIName,
				Scope: di.App,
				Build: func(ctn di.Container) (any, error) {
					svc := ctn.Get(AuthServiceDIName).(*services.AuthService)
					return handlers.NewAuthHandler(svc), nil
				},
				Close: func(obj any) error {
					return nil
				},
			},
			di.Def{
				Name:  UserAPIDIName,
				Scope: di.App,
				Build: func(ctn di.Container) (any, error) {
					svc := ctn.Get(UserServiceDIName).(*services.UserService)
					rbac := ctn.Get(RBACServiceDIName).(*services.RBACService)
					return handlers.NewUserHandler(svc, rbac), nil
				},
				Close: func(obj any) error {
					return nil
				},
			},
			di.Def{
				Name:  BranchAPIDIName,
				Scope: di.App,
				Build: func(ctn di.Container) (any, error) {
					cuc := ctn.Get(BranchServiceDIName).(*services.BranchService)
					rbac := ctn.Get(RBACServiceDIName).(*services.RBACService)
					return handlers.NewBranchHandler(cuc, rbac), nil
				},
				Close: func(obj any) error {
					return nil
				},
			},
			di.Def{
				Name:  BranchWifiAPIDIName,
				Scope: di.App,
				Build: func(ctn di.Container) (any, error) {
					svc := ctn.Get(BranchWifiServiceDIName).(*services.BranchWifiService)
					return handlers.NewBranchWifiHandler(svc), nil
				},
				Close: func(obj any) error {
					return nil
				},
			},
			di.Def{
				Name:  AttendanceAPIDIName,
				Scope: di.App,
				Build: func(ctn di.Container) (any, error) {
					svc := ctn.Get(AttendanceServiceDIName).(*services.AttendanceService)
					rbac := ctn.Get(RBACServiceDIName).(*services.RBACService)
					return handlers.NewAttendanceHandler(svc, rbac), nil
				},
				Close: func(obj any) error {
					return nil
				},
			},
			di.Def{
				Name:  DeviceAPIDIName,
				Scope: di.App,
				Build: func(ctn di.Container) (any, error) {
					svc := ctn.Get(DeviceServiceDIName).(*services.DeviceService)
					rbac := ctn.Get(RBACServiceDIName).(*services.RBACService)
					return handlers.NewDeviceHandler(svc, rbac), nil
				},
				Close: func(obj any) error {
					return nil
				},
			},
			di.Def{
				Name:  DashboardAPIDIName,
				Scope: di.App,
				Build: func(ctn di.Container) (any, error) {
					svc := ctn.Get(DashboardServiceDIName).(*services.DashboardService)
					rbac := ctn.Get(RBACServiceDIName).(*services.RBACService)
					return handlers.NewDashboardHandler(svc, rbac), nil
				},
				Close: func(obj any) error {
					return nil
				},
			},
		)

		return arr
	}
}
